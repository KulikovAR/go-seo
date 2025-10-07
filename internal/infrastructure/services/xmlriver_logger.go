package services

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type XMLRiverLogger struct {
	logger *log.Logger
	file   *os.File
}

func NewXMLRiverLogger() (*XMLRiverLogger, error) {
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	logFile := filepath.Join(logDir, fmt.Sprintf("xmlriver_%s.log", time.Now().Format("2006-01-02")))

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	logger := log.New(file, "[XMLRiver] ", log.LstdFlags|log.Lshortfile)

	return &XMLRiverLogger{
		logger: logger,
		file:   file,
	}, nil
}

func (l *XMLRiverLogger) LogRequest(url string, params map[string]string) {
	l.logger.Printf("REQUEST: %s", url)
	l.logger.Printf("PARAMS: %+v", params)
	l.logger.Println("---")
}

func (l *XMLRiverLogger) LogResponse(statusCode int, body []byte) {
	l.logger.Printf("RESPONSE STATUS: %d", statusCode)
	l.logger.Printf("RESPONSE BODY:\n%s", string(body))
	l.logger.Println("==========================================")
}

func (l *XMLRiverLogger) LogError(err error, context string) {
	l.logger.Printf("ERROR [%s]: %v", context, err)
}

func (l *XMLRiverLogger) Close() error {
	return l.file.Close()
}

func parseXMLResponseWithLogging(body io.Reader, logger *XMLRiverLogger) (*SearchResponse, error) {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		logger.LogError(err, "failed to read response body")
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	logger.LogResponse(200, bodyBytes)

	var searchResp SearchResponse
	if err := xml.Unmarshal(bodyBytes, &searchResp); err != nil {
		logger.LogError(err, "failed to parse XML response")
		return nil, fmt.Errorf("failed to parse XML response: %w", err)
	}

	return &searchResp, nil
}
