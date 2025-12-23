package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGray   = "\033[90m"
)

type Logger struct {
	level     LogLevel
	debugMode bool
}

func New(level LogLevel) *Logger {
	return &Logger{level: level}
}

func NewFromEnv() *Logger {
	_ = godotenv.Load()

	rawLevel := strings.ToUpper(os.Getenv("LOG_LEVEL_THRESHOLD"))
	rawDebug := strings.ToLower(os.Getenv("DEBUG"))

	level := parseLevel(rawLevel)
	debug := rawDebug == "true"

	logger := &Logger{level: level, debugMode: debug}

	logger.Info("💡 Logger Mode: %v", rawDebug)
	logger.Info("⚙️ Log Level: %s", level)
	return logger
}

func parseLevel(s string) LogLevel {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return DEBUG
	case "WARN":
		return WARN
	case "ERROR":
		return ERROR
	default:
		return INFO
	}
}

// Logging methods
func (l *Logger) shouldLog(level LogLevel) bool {
	order := map[LogLevel]int{
		DEBUG: 0,
		INFO:  1,
		WARN:  2,
		ERROR: 3,
	}
	return order[level] >= order[l.level]
}

func (l *Logger) Info(msg string, args ...interface{}) {
	if l.shouldLog(INFO) {
		log.Printf(colorBlue+"[INFO]  "+colorReset+msg, args...)
	}
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	if l.shouldLog(WARN) {
		log.Printf(colorYellow+"[WARN]  "+colorReset+msg, args...)
	}
}

func (l *Logger) Error(msg string, args ...interface{}) {
	if l.shouldLog(ERROR) {
		log.Printf(colorRed+"[ERROR] "+colorReset+msg, args...)
	}
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	if l.debugMode && l.shouldLog(DEBUG) {
		log.Printf(colorGray+"[DEBUG] "+colorReset+msg, args...)
	}
}

func (l *Logger) LogStructured(level string, msg string, args ...interface{}) {
	if !l.shouldLog(LogLevel(level)) {
		return
	}
	entry := map[string]interface{}{
		"level":   level,
		"message": fmt.Sprintf(msg, args...),
		"time":    time.Now().Format(time.RFC3339),
	}
	jsonBytes, _ := json.Marshal(entry)
	fmt.Println(string(jsonBytes))
}

func (l *Logger) LogPretty(title string, v interface{}) {
	if !l.shouldLog(INFO) {
		return
	}	
	fmt.Println("------------------------------------------------------------------------")
	fmt.Println("* ", title)
	fmt.Println("------------------------------------------------------------------------")

	// Handle error types explicitly
	if err, ok := v.(error); ok {
		fmt.Println(err.Error())
		// return
	}

	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("Failed to marshal object:", err)
		return
	}
	fmt.Println(string(bytes))
	fmt.Println("_____")
}
