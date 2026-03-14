package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	HealthURL     string
	CheckInterval time.Duration
	Timeout       time.Duration

	SMTPHost string
	SMTPPort string
	SMTPUser string
	SMTPPass string
	MailFrom string
	MailTo   string

	TorchiToken    string
	HeartbeatHour  int
}

func Load(envFile string) (Config, error) {
	if err := LoadDotEnv(envFile); err != nil {
		return Config{}, fmt.Errorf("failed to load %s: %w", envFile, err)
	}

	interval, _ := strconv.Atoi(getEnv("CHECK_INTERVAL_SEC", "60"))
	timeout, _ := strconv.Atoi(getEnv("HTTP_TIMEOUT_SEC", "10"))

	heartbeatHour, _ := strconv.Atoi(getEnv("HEARTBEAT_HOUR", "9"))

	cfg := Config{
		HealthURL:     getEnv("HEALTH_URL", "https://torchi.app/api/health"),
		CheckInterval: time.Duration(interval) * time.Second,
		Timeout:       time.Duration(timeout) * time.Second,
		SMTPHost:      getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:      getEnv("SMTP_PORT", "587"),
		SMTPUser:      os.Getenv("SMTP_USER"),
		SMTPPass:      os.Getenv("SMTP_PASS"),
		MailFrom:      os.Getenv("SMTP_USER"),
		MailTo:        os.Getenv("MAIL_TO"),
		TorchiToken:   os.Getenv("TORCHI_TOKEN"),
		HeartbeatHour: heartbeatHour,
	}

	if cfg.SMTPUser == "" || cfg.SMTPPass == "" || cfg.MailTo == "" {
		return Config{}, fmt.Errorf("SMTP_USER, SMTP_PASS, MAIL_TO are required")
	}
	if cfg.TorchiToken == "" {
		return Config{}, fmt.Errorf("TORCHI_TOKEN is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
