package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// IDGeneratorService generates unique IDs for tracking tasks
type IDGeneratorService struct{}

// NewIDGeneratorService creates a new ID generator service
func NewIDGeneratorService() *IDGeneratorService {
	return &IDGeneratorService{}
}

// GenerateTaskID generates a unique task ID
func (s *IDGeneratorService) GenerateTaskID() string {
	return s.generateID("task")
}

// GenerateJobID generates a unique job ID
func (s *IDGeneratorService) GenerateJobID() string {
	return s.generateID("job")
}

// generateID generates a unique ID with prefix
func (s *IDGeneratorService) generateID(prefix string) string {
	// Generate 8 random bytes
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based ID if random generation fails
		return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
	}
	
	// Convert to hex and add prefix
	return fmt.Sprintf("%s_%s", prefix, hex.EncodeToString(bytes))
}

