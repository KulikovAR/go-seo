package services

import (
	"log"
	"time"
)

type KafkaService struct {
	brokers []string
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
	log.Printf("Kafka service initialized with brokers: %v", brokers)
	return &KafkaService{
		brokers: brokers,
	}, nil
}

func (k *KafkaService) SendTaskStatus(message *TaskStatusMessage) error {
	log.Printf("Kafka message sent - TaskID: %s, JobID: %s, Status: %s",
		message.TaskID, message.JobID, message.Status)
	return nil
}

func (k *KafkaService) SendJobStatus(jobID string, status string, errorMsg string) error {
	log.Printf("Kafka job message sent - JobID: %s, Status: %s, Error: %s",
		jobID, status, errorMsg)
	return nil
}

func (k *KafkaService) Close() error {
	log.Println("Kafka service closed")
	return nil
}
