package main

import (
	"context"
	"fmt"
	"go-kafka-example/config"
	"go-kafka-example/internal/server"
	"go-kafka-example/pkg/database"
	"go-kafka-example/pkg/kafka"
	"go-kafka-example/pkg/logger"
	"go-kafka-example/pkg/tracer"

	"log"
)

func main() {
	log.Println("Starting go-example-kafka-api server")
	cfg, err := config.GetConf()

	if err != nil {
		panic(fmt.Errorf("load config: %v", err))
	}

	apiLogger := logger.NewApiLogger(cfg)
	apiLogger.InitLogger()

	// init database
	db, err := database.ConnectDB(&cfg.Database)
	if err != nil {
		panic(fmt.Errorf("load database: %v", err))
	}
	defer db.Close()

	traceProvider, err := tracer.NewJaeger(cfg)
	if err != nil {
		apiLogger.Fatal("Cannot create tracer", err)
	} else {
		apiLogger.Info("Jaeger connected")
	}
	defer func() {
		if err := traceProvider.Shutdown(context.Background()); err != nil {
			apiLogger.Error("Cannot shutdown tracer", err)
		}
	}()

	kakfaWriter, err := kafka.NewWriter(cfg.Kafka, "user_email")
	if err != nil {
		apiLogger.Fatal("Can't connect with kafka", err)
	}
	defer kakfaWriter.Close()

	// init server
	srv := server.NewServer(cfg, db, kakfaWriter, apiLogger)
	if err = srv.Run(); err != nil {
		apiLogger.Fatal(err)
	}

}
