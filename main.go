package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/IBM/sarama"
)

func main() {
	// Define flags
	brokerList := flag.String("brokers", "localhost:9092", "Comma-separated list of Kafka brokers")
	createTopic := flag.String("create", "", "Name of the Kafka topic to create")
	deleteTopic := flag.String("delete", "", "Name of the Kafka topic to delete")
	listTopics := flag.Bool("list", false, "List all the Kafka topics")
	listCG := flag.Bool("listcg", false, "List all the Kafka Consumer Groups")
	replicationFactor := flag.Int("replication", 1, "Replication factor for the topic")
	partitions := flag.Int("partitions", 1, "Number of partitions for the topic")
	// Parse flags
	flag.Parse()

	// Create Kafka configuration
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0 // Use the desired Kafka version

	// Create Kafka admin client
	adminClient, err := sarama.NewClusterAdmin(strings.Split(*brokerList, ","), config)
	if err != nil {
		log.Fatalf("Error creating Kafka admin client: %v", err)
	}
	defer func() {
		if err := adminClient.Close(); err != nil {
			log.Fatalf("Error closing Kafka admin client: %v", err)
		}
	}()

	// Create or delete topic based on input
	if *createTopic != "" {
		createKafkaTopic(adminClient, *createTopic, *partitions, *replicationFactor)
	} else if *deleteTopic != "" {
		deleteKafkaTopic(adminClient, *deleteTopic)
	} else if *listTopics {
		listKafkaTopics(adminClient)
	} else if *listCG {
		listConsumerGroups(adminClient)
	} else {
		fmt.Println("Unknown Command")
	}
}

func listKafkaTopics(adminClient sarama.ClusterAdmin) {
	topics, err := adminClient.ListTopics()
	if err != nil {
		log.Fatalf("Error listing kafka topics: %v\n", err)
	}
	for topic := range topics {
		fmt.Println(topic)
	}

}

// createKafkaTopic creates a Kafka topic
func createKafkaTopic(adminClient sarama.ClusterAdmin, topic string, partitions int, replicationFactor int) {
	// Check if the topic already exists
	topics, err := adminClient.ListTopics()
	if err != nil {
		log.Fatalf("Error listing Kafka topics: %v", err)
	}

	if topicExists(topics, topic) {
		fmt.Printf("Kafka topic '%s' already exists\n", topic)
		return
	}

	// Create topic configuration
	topicConfig := &sarama.TopicDetail{
		NumPartitions:     int32(partitions),
		ReplicationFactor: int16(replicationFactor),
		ConfigEntries:     nil,
	}

	// Create the topic
	err = adminClient.CreateTopic(topic, topicConfig, false)
	if err != nil {
		log.Fatalf("Error creating Kafka topic: %v", err)
	}

	fmt.Printf("Kafka topic '%s' created successfully with replication factor %d\n", topic, replicationFactor)
}

// deleteKafkaTopic deletes a Kafka topic
func deleteKafkaTopic(adminClient sarama.ClusterAdmin, topic string) {
	// Check if the topic exists
	topics, err := adminClient.ListTopics()
	if err != nil {
		log.Fatalf("Error listing Kafka topics: %v", err)
	}

	if !topicExists(topics, topic) {
		fmt.Printf("Kafka topic '%s' does not exist\n", topic)
		return
	}

	// Delete the topic
	err = adminClient.DeleteTopic(topic)
	if err != nil {
		log.Fatalf("Error deleting Kafka topic: %v", err)
	}

	fmt.Printf("Kafka topic '%s' deleted successfully\n", topic)
}

// topicExists checks if a topic exists in the list of topics
func topicExists(topics map[string]sarama.TopicDetail, topic string) bool {
	_, exists := topics[topic]
	return exists
}

func listConsumerGroups(adminClient sarama.ClusterAdmin) {
	groups, err := adminClient.ListConsumerGroups()
	if err != nil {
		log.Fatal("Error listing Kafka consumer groups: %v", err)
	}

	for group := range groups {
		fmt.Println(group)
	}
}
