package logger

import (
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/exp/slog"
)

var LogLevel = map[string]slog.Level{
	"DEBUG": slog.LevelDebug,
	"INFO":  slog.LevelInfo,
	"WARN":  slog.LevelWarn,
	"ERROR": slog.LevelError,
}

func GetLoglevel() string {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file: %v", err)
	}

	logLevel := os.Getenv("LogLevel")
	if logLevel == "" {
		slog.Error("LogLevel is not set in environment variables")
		os.Exit(1)
	}

	return logLevel
}

func GetLogger() *slog.Logger {
	logLevel := GetLoglevel()
	opts := &slog.HandlerOptions{
		Level: slog.Level(LogLevel[logLevel]),
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)

	return logger

}
