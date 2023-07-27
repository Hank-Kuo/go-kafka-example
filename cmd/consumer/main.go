package main

import (
	"context"
	"fmt"
	"go-kafka-example/pkg/kafka"
	"os"
	"os/signal"
)

func Consumer(ctx context.Context) {
	r := kafka.NewReader([]string{"localhost:9092"}, "topic1")
	fmt.Println("listen message")
	for {
		m, err := r.ReadMessage(ctx)
		fmt.Println(err)
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}
}
func main() {

	ctx := context.Background()
	// defer r.Close()
	Consumer(ctx)
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

}
