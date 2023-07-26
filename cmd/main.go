package main

import (
	"context"
	"fmt"
	"go-kafka-example/config"

	"go-kafka-example/internal/server"
	"go-kafka-example/pkg/database"
	"go-kafka-example/pkg/logger"
	"go-kafka-example/pkg/tracer"

	"log"

	"github.com/segmentio/kafka-go"
)

func newKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

/*
kafkaURL := "localhost:9092"
topic := "topic1"

writer := newKafkaWriter(kafkaURL, topic)
defer writer.Close()
fmt.Println("start producing ... !!")

	for i := 0; ; i++ {
		key := fmt.Sprintf("Key-%d", i)
		msg := kafka.Message{
			Value: []byte(fmt.Sprintf("name-%d", i)),
		}
		err := writer.WriteMessages(context.Background(), msg)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("produced", key)
		}
		time.Sleep(1 * time.Second)
	}
*/
func main() {

	log.Println("Starting iShare-api server")
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

	// // init server
	srv := server.NewServer(cfg, db, apiLogger)
	if err = srv.Run(); err != nil {
		apiLogger.Fatal(err)
	}

}
