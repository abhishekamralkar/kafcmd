package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/abhishekamralkar/kafcmd/pkg/kafcmd"
)

var (
	brokers        = flag.String("brokers", "localhost:9092", "Kafka broker addresses (comma-separated)")
	createTopic    = flag.String("create", "", "Create a Kafka topic")
	deleteTopic    = flag.String("delete", "", "Delete a Kafka topic")
	listTopic      = flag.Bool("list", false, "List all Kafka topics")
	partitions     = flag.Int("partitions", 1, "Number of partitions for the topic")
	replication    = flag.Int("replication", 1, "Replication factor for the topic")
	sessionTimeout = flag.Duration("session-timeout", 10*time.Second, "Session timeout for topic creation (e.g., 10s)")
)

func main() {
	flag.Parse()
	config := sarama.NewConfig()
	admin, err := sarama.NewClusterAdmin(strings.Split(*brokers, ","), config)
	if err != nil {
		log.Fatalf("Failed to create kafka admin client: %v", err)
	}

	defer func() {
		err = admin.Close()
		if err != nil {
			log.Fatalf("Error closing kafka admin client: %v", err)
		}
	}()

	switch {
	case *createTopic != "":
		kafcmd.Create(admin, *createTopic, *partitions, *replication, *sessionTimeout)
	case *deleteTopic != "":
		kafcmd.Delete(admin, *deleteTopic)
	case *listTopic:
		kafcmd.List(admin)
	default:
		fmt.Println("Usage: kafcmd [options]")
		flag.PrintDefaults()
		os.Exit(1)
	}

}
