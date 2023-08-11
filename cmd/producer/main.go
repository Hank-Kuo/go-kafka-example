package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Hank-Kuo/go-kafka-example/config"
	"github.com/Hank-Kuo/go-kafka-example/internal/server"
	"github.com/Hank-Kuo/go-kafka-example/pkg/database"
	"github.com/Hank-Kuo/go-kafka-example/pkg/kafka"
	"github.com/Hank-Kuo/go-kafka-example/pkg/logger"
	"github.com/Hank-Kuo/go-kafka-example/pkg/tracer"
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
	srv.Run()
	// if err = ; err != nil {
	// 	apiLogger.Fatal(err)
	// }

}
