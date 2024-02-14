package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/IBM/sarama"
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
		createT(admin, *createTopic, *partitions, *replication, *sessionTimeout)
	case *deleteTopic != "":
		deleteT(admin, *deleteTopic)
	case *listTopic:
		listT(admin)
	default:
		fmt.Println("Usage: kafcmd [options]")
		flag.PrintDefaults()
		os.Exit(1)
	}

}

func listT(admin sarama.ClusterAdmin) {
	topics, err := admin.ListTopics()
	if err != nil {
		log.Fatalf("Error listing kafka topics: %v\n", err)
	}
	for topic := range topics {
		fmt.Println(topic)
	}

}

func createT(admin sarama.ClusterAdmin, topic string, partitions, replication int, sessionTimeout time.Duration) {
	topics, err := admin.ListTopics()
	if err != nil {
		log.Fatalf("Error listing Kafka topics: %v", err)
	}

	if topicExists(topics, topic) {
		fmt.Printf("Kafka topic '%s' already exists\n", topic)
		return
	}

	topicConfig := &sarama.TopicDetail{
		NumPartitions:     int32(partitions),
		ReplicationFactor: int16(replication),
		ConfigEntries:     nil,
	}

	err = admin.CreateTopic(topic, topicConfig, false)
	if err != nil {
		log.Fatalf("Error creating Kafka topic: %v", err)
	}

	fmt.Printf("Kafka topic '%s' created successfully with partitions %d, replication factor %d\n", topic, partitions, replication)
}

func deleteT(admin sarama.ClusterAdmin, topic string) {
	topics, err := admin.ListTopics()
	if err != nil {
		log.Fatalf("Error listing Kafka topics: %v", err)
	}

	if !topicExists(topics, topic) {
		fmt.Printf("Kafka topic '%s' does not exist\n", topic)
		return
	}

	err = admin.DeleteTopic(topic)
	if err != nil {
		log.Fatalf("Error deleting Kafka topic: %v", err)
	}

	fmt.Printf("Kafka topic '%s' deleted successfully\n", topic)
}

func topicExists(topics map[string]sarama.TopicDetail, topic string) bool {
	_, exists := topics[topic]
	return exists
}
