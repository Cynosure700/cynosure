package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func buildDefaultMySQLDSN() string {
	host := getenv("MYSQL_HOST", "127.0.0.1")
	port := getenv("MYSQL_PORT", "3306")
	user := getenv("MYSQL_USER", getenv("DB_USER", "root"))
	password := getenv("MYSQL_PASSWORD", getenv("DB_PASSWORD", ""))
	database := getenv("MYSQL_DATABASE", getenv("DB_NAME", "vibe_coding"))
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&multiStatements=true&loc=Local", user, password, host, port, database)
}

func getenv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func getenvInt(key string, fallback int) (int, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}

func getenvIntOrDefault(key string, fallback int) int {
	v, err := getenvInt(key, fallback)
	if err != nil {
		return fallback
	}
	return v
}

func getenvBool(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func parseCSVList(value string) []string {
	parts := strings.Split(value, ",")
	items := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item == "" {
			continue
		}
		if _, exists := seen[item]; exists {
			continue
		}
		seen[item] = struct{}{}
		items = append(items, item)
	}
	return items
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
