package utils

import (
	"fmt"
	"log"
	"os"
	"simple-crud-notes/utils/enum"
	"time"
)

type LogEntry struct {
	Timestamp       time.Time `json:"timestamp"`
	StatusCode      int       `json:"status_code"`
	Method          string    `json:"method"`
	Endpoint        string    `json:"endpoint"`
	RequestHeader   string    `json:"request_header"`
	RequestBody     string    `json:"request_body"`
	RequestParams   string    `json:"request_params"`
	ResponseMessage string    `json:"response_message"`
	ProcessTime     string    `json:"process_time"`
}

var logChan = make(chan LogEntry, 1000) // Buffer channel

func init() {
	go processLogs()
}

func processLogs() {
	logFileName := fmt.Sprintf("logs/apps-%s.log", time.Now().Format("2006-01-02"))
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		return
	}
	defer file.Close()

	for entry := range logChan {
		log.Printf("logChan %+v", entry)

		var logType string
		var statusCode = entry.StatusCode

		switch {
		case statusCode >= 200 && statusCode < 300:
			logType = enum.LOG_TYPE_SUCCESS
		case statusCode >= 400 && statusCode < 500:
			logType = enum.LOG_TYPE_WARN
		default:
			logType = enum.LOG_TYPE_ERROR
		}

		logLine := fmt.Sprintf(
			"%s|%s|%d|%s|%s|%s|%s|%s|%s|%s\n",
			entry.Timestamp.Format("2006-01-02 15:04:05.000"),
			logType,
			entry.StatusCode,
			entry.Method,
			entry.Endpoint,
			toJSONString(entry.RequestHeader),
			toJSONString(entry.RequestBody),
			toJSONString(entry.RequestParams),
			toJSONString(entry.ResponseMessage),
			entry.ProcessTime,
		)

		if _, err := file.WriteString(logLine); err != nil {
			fmt.Println("Log write error:", err)
		}
	}
}

func LogAsync(entry LogEntry) {
	select {
	case logChan <- entry:
	default:
		// Fallback jika channel penuh
		fmt.Println("Log channel overflow. Entry dropped.")
	}
}
