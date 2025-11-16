package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type KafkaService struct {
	producer sarama.SyncProducer
	admin    sarama.ClusterAdmin
	enabled  bool
}

type TaskStatusMessage struct {
	TaskID    string      `json:"task_id"`
	JobID     string      `json:"job_id"`
	Status    string      `json:"status"`
	Timestamp time.Time   `json:"timestamp"`
	Error     string      `json:"error,omitempty"`
	Result    *TaskResult `json:"result,omitempty"`
}

type TaskResult struct {
	KeywordID int    `json:"keyword_id"`
	SiteID    int    `json:"site_id"`
	Source    string `json:"source"`
	Rank      int    `json:"rank"`
	URL       string `json:"url"`
	Title     string `json:"title"`
	Success   bool   `json:"success"`
}

func NewKafkaService(brokers []string) (*KafkaService, error) {
	enabled := len(brokers) > 0 && brokers[0] != ""

	service := &KafkaService{
		enabled: enabled,
	}

	if enabled {
		if err := service.initKafka(brokers); err != nil {
			log.Printf("Warning: Kafka initialization failed: %v. Service will work in log-only mode.", err)
			service.enabled = false
		} else {
			log.Printf("Kafka service initialized with brokers: %v", brokers)
			service.createTopicsIfNotExist()
		}
	} else {
		log.Println("Kafka service initialized in log-only mode (no brokers configured)")
	}

	return service, nil
}

func (k *KafkaService) initKafka(brokers []string) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Timeout = 10 * time.Second

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return fmt.Errorf("failed to create producer: %w", err)
	}

	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		producer.Close()
		return fmt.Errorf("failed to create admin client: %w", err)
	}

	k.producer = producer
	k.admin = admin

	return nil
}

func (k *KafkaService) createTopicsIfNotExist() {
	topics := []string{"tracking-status", "tracking-jobs"}

	for _, topic := range topics {
		if err := k.createTopic(topic); err != nil {
			log.Printf("Warning: Failed to create topic %s: %v (topics will be auto-created on first message if enabled)", topic, err)
		} else {
			log.Printf("Topic %s is ready", topic)
		}
	}
}

func (k *KafkaService) createTopic(topicName string) error {
	topics, err := k.admin.ListTopics()
	if err != nil {
		return fmt.Errorf("failed to list topics: %w", err)
	}

	if _, exists := topics[topicName]; exists {
		log.Printf("Topic %s already exists", topicName)
		return nil
	}

	topicDetail := &sarama.TopicDetail{
		NumPartitions:     3,
		ReplicationFactor: 1,
		ConfigEntries: map[string]*string{
			"cleanup.policy":   stringPtr("delete"),
			"retention.ms":     stringPtr("604800000"),
			"compression.type": stringPtr("snappy"),
		},
	}

	err = k.admin.CreateTopic(topicName, topicDetail, false)
	if err != nil {
		return fmt.Errorf("failed to create topic: %w", err)
	}

	log.Printf("Successfully created topic: %s", topicName)
	return nil
}

func stringPtr(s string) *string {
	return &s
}

func (k *KafkaService) SendTaskStatus(message *TaskStatusMessage) error {
	if !k.enabled {
		log.Printf("Kafka message sent (log-only) - TaskID: %s, JobID: %s, Status: %s",
			message.TaskID, message.JobID, message.Status)
		return nil
	}

	return k.sendMessage("tracking-status", message)
}

func (k *KafkaService) SendJobStatus(jobID string, status string, errorMsg string, percent ...int) error {
	if !k.enabled {
		if len(percent) > 0 {
			log.Printf("Kafka job message sent (log-only) - JobID: %s, Status: %s, Error: %s, Percent: %d",
				jobID, status, errorMsg, percent[0])
		} else {
			log.Printf("Kafka job message sent (log-only) - JobID: %s, Status: %s, Error: %s",
				jobID, status, errorMsg)
		}
		return nil
	}

	message := map[string]interface{}{
		"job_id":    jobID,
		"status":    status,
		"error":     errorMsg,
		"timestamp": time.Now(),
	}

	if len(percent) > 0 {
		message["percent"] = percent[0]
	}

	return k.sendMessage("tracking-jobs", message)
}

func (k *KafkaService) sendMessage(topic string, message interface{}) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Printf("ERROR: Failed to marshal Kafka message for topic %s: %v", topic, err)
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	kafkaMessage := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(jsonData),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("content-type"),
				Value: []byte("application/json"),
			},
		},
	}

	partition, offset, err := k.producer.SendMessage(kafkaMessage)
	if err != nil {
		log.Printf("ERROR: Failed to send message to Kafka topic %s: %v", topic, err)
		return fmt.Errorf("failed to send message: %w", err)
	}

	log.Printf("Successfully sent message to topic %s, partition %d, offset %d", topic, partition, offset)
	return nil
}

func (k *KafkaService) Close() error {
	if k.producer != nil {
		k.producer.Close()
	}
	if k.admin != nil {
		k.admin.Close()
	}
	log.Println("Kafka service closed")
	return nil
}
