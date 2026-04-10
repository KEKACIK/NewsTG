package config

import "fmt"

type Config struct {
	Debug            bool
	TelegramBotToken string
	TelegramChatID   int64

	PostgresHost     string
	PostgresPort     int
	PostgresName     string
	PostgresUser     string
	PostgresPassword string
}

func NewConfig() *Config {
	initConfig()

	return &Config{
		Debug: getBoolEnv("DEBUG", false),

		TelegramBotToken: getStrEnv("TELEGRAM_BOT_TOKEN", ""),
		TelegramChatID:   getInt64Env("TELEGRAM_CHAT_ID", 0),

		PostgresHost:     getStrEnv("POSTGRES_HOST", "127.0.0.1"),
		PostgresPort:     getIntEnv("POSTGRES_PORT", 5432),
		PostgresName:     getStrEnv("POSTGRES_NAME", "postgres"),
		PostgresUser:     getStrEnv("POSTGRES_USER", "postgres"),
		PostgresPassword: getStrEnv("POSTGRES_PASSWORD", "postgres"),
	}
}

func (c *Config) GetPostgresDsn() string {
	url := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", c.PostgresUser, c.PostgresPassword, c.PostgresHost, c.PostgresPort, c.PostgresName)

	return url
}
