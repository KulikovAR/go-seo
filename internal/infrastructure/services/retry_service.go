package services

import (
	"fmt"
	"time"
)

// RetryService handles retry logic with exponential backoff
type RetryService struct {
	maxRetries int
	baseDelay  time.Duration
}

// NewRetryService creates a new retry service
func NewRetryService(maxRetries int, baseDelay time.Duration) *RetryService {
	return &RetryService{
		maxRetries: maxRetries,
		baseDelay:  baseDelay,
	}
}

// RetryConfig represents retry configuration
type RetryConfig struct {
	MaxRetries int
	BaseDelay  time.Duration
}

// DefaultRetryConfig returns default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries: 5,
		BaseDelay:  10 * time.Second,
	}
}

// ExecuteWithRetry executes a function with retry logic
func (r *RetryService) ExecuteWithRetry(fn func() error) error {
	var lastErr error
	
	for attempt := 0; attempt <= r.maxRetries; attempt++ {
		if attempt > 0 {
			// Calculate delay with exponential backoff
			delay := r.calculateDelay(attempt)
			time.Sleep(delay)
		}
		
		err := fn()
		if err == nil {
			return nil
		}
		
		lastErr = err
	}
	
	return fmt.Errorf("operation failed after %d attempts, last error: %w", r.maxRetries+1, lastErr)
}

// calculateDelay calculates the delay for the given attempt
func (r *RetryService) calculateDelay(attempt int) time.Duration {
	// Exponential backoff: baseDelay * 2^(attempt-1)
	// Attempt 1: baseDelay
	// Attempt 2: baseDelay * 2
	// Attempt 3: baseDelay * 4
	// Attempt 4: baseDelay * 8
	// Attempt 5: baseDelay * 16
	
	multiplier := 1 << (attempt - 1) // 2^(attempt-1)
	return r.baseDelay * time.Duration(multiplier)
}

// GetRetryDelay returns the delay for a specific retry attempt
func (r *RetryService) GetRetryDelay(attempt int) time.Duration {
	if attempt <= 0 {
		return 0
	}
	return r.calculateDelay(attempt)
}

