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
}

func Load() (Config, error) {
	if err := LoadDotEnv(".env"); err != nil {
		return Config{}, fmt.Errorf("failed to load .env: %w", err)
	}

	interval, _ := strconv.Atoi(getEnv("CHECK_INTERVAL_SEC", "60"))
	timeout, _ := strconv.Atoi(getEnv("HTTP_TIMEOUT_SEC", "10"))

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
	}

	if cfg.SMTPUser == "" || cfg.SMTPPass == "" || cfg.MailFrom == "" || cfg.MailTo == "" {
		return Config{}, fmt.Errorf("SMTP_USER, SMTP_PASS, MAIL_FROM, MAIL_TO are required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
