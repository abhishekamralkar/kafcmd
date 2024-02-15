package kafcmd

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

func List(admin sarama.ClusterAdmin) {
	topics, err := admin.ListTopics()
	if err != nil {
		log.Fatalf("Error listing kafka topics: %v\n", err)
	}
	for topic := range topics {
		fmt.Println(topic)
	}

}

func Create(admin sarama.ClusterAdmin, topic string, partitions, replication int, sessionTimeout time.Duration) {
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

func Delete(admin sarama.ClusterAdmin, topic string) {
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
