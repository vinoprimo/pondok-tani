package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type SMTPConfig struct {
	Host        string
	Port        string
	EnableSSL   bool
	Email       string
	DisplayName string
	Password    string
	Timeout     time.Duration
}

var SMTP SMTPConfig

func LoadSMTPConfig() {
	SMTP = SMTPConfig{
		Host:        strings.TrimSpace(os.Getenv("SMTP_HOST")),
		Port:        strings.TrimSpace(os.Getenv("SMTP_PORT")),
		EnableSSL:   strings.ToLower(strings.TrimSpace(os.Getenv("SMTP_ENABLE_SSL"))) == "true",
		Email:       strings.TrimSpace(os.Getenv("SMTP_EMAIL")),
		DisplayName: strings.TrimSpace(os.Getenv("SMTP_DISPLAY_NAME")),
		Password:    strings.TrimSpace(os.Getenv("SMTP_PASSWORD")),
		Timeout:     30 * time.Second,
	}

	if timeout := strings.TrimSpace(os.Getenv("SMTP_TIMEOUT")); timeout != "" {
		if millis, err := strconv.Atoi(timeout); err == nil && millis > 0 {
			SMTP.Timeout = time.Duration(millis) * time.Millisecond
		}
	}
}
