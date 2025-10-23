package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type KafkaService struct {
	brokers    []string
	httpClient *http.Client
	enabled    bool
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
		brokers: brokers,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		enabled: enabled,
	}

	if enabled {
		if err := service.checkConnection(); err != nil {
			log.Printf("Warning: Kafka connection failed: %v. Service will work in log-only mode.", err)
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

func (k *KafkaService) checkConnection() error {
	if len(k.brokers) == 0 {
		return fmt.Errorf("no brokers configured")
	}

	broker := k.brokers[0]
	if broker == "" {
		return fmt.Errorf("empty broker address")
	}

	resp, err := k.httpClient.Get(broker + "/topics")
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Kafka returned status: %d", resp.StatusCode)
	}

	return nil
}

func (k *KafkaService) createTopicsIfNotExist() {
	topics := []string{"tracking-status", "tracking-jobs"}

	for _, topic := range topics {
		if err := k.createTopic(topic); err != nil {
			log.Printf("Warning: Failed to create topic %s: %v", topic, err)
		} else {
			log.Printf("Topic %s is ready", topic)
		}
	}
}

func (k *KafkaService) createTopic(topicName string) error {
	exists, err := k.topicExists(topicName)
	if err != nil {
		return fmt.Errorf("failed to check if topic exists: %w", err)
	}

	if exists {
		log.Printf("Topic %s already exists", topicName)
		return nil
	}

	topicConfig := map[string]interface{}{
		"name": topicName,
		"configs": map[string]string{
			"cleanup.policy": "delete",
			"retention.ms":   "604800000",
		},
	}

	jsonData, err := json.Marshal(topicConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal topic config: %w", err)
	}

	req, err := http.NewRequest("POST", k.brokers[0]+"/topics", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := k.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create topic: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create topic, status: %d, body: %s", resp.StatusCode, string(body))
	}

	log.Printf("Successfully created topic: %s", topicName)
	return nil
}

func (k *KafkaService) topicExists(topicName string) (bool, error) {
	resp, err := k.httpClient.Get(k.brokers[0] + "/topics/" + topicName)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

func (k *KafkaService) SendTaskStatus(message *TaskStatusMessage) error {
	if !k.enabled {
		log.Printf("Kafka message sent (log-only) - TaskID: %s, JobID: %s, Status: %s",
			message.TaskID, message.JobID, message.Status)
		return nil
	}

	return k.sendMessage("tracking-status", message)
}

func (k *KafkaService) SendJobStatus(jobID string, status string, errorMsg string) error {
	if !k.enabled {
		log.Printf("Kafka job message sent (log-only) - JobID: %s, Status: %s, Error: %s",
			jobID, status, errorMsg)
		return nil
	}

	message := map[string]interface{}{
		"job_id":    jobID,
		"status":    status,
		"error":     errorMsg,
		"timestamp": time.Now(),
	}

	return k.sendMessage("tracking-jobs", message)
}

func (k *KafkaService) sendMessage(topic string, message interface{}) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	kafkaMsg := map[string]interface{}{
		"records": []map[string]interface{}{
			{
				"value": string(jsonData),
			},
		},
	}

	kafkaData, err := json.Marshal(kafkaMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal kafka message: %w", err)
	}

	url := fmt.Sprintf("%s/topics/%s", k.brokers[0], topic)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(kafkaData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/vnd.kafka.json.v2+json")

	resp, err := k.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to send message, status: %d, body: %s", resp.StatusCode, string(body))
	}

	log.Printf("Successfully sent message to topic %s", topic)
	return nil
}

func (k *KafkaService) Close() error {
	log.Println("Kafka service closed")
	return nil
}
