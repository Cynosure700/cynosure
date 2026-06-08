package config

import (
	"fmt"
	"os"
	"strings"
)

// getenv 仅用于读取不应写入 config.json 的敏感信息（如密钥、密码）。
func getenv(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

func buildMySQLDSN(fileCfg fileConfig) string {
	host := firstNonEmpty(fileCfg.MySQLHost, "1.12.217.28")
	port := firstNonEmpty(fileCfg.MySQLPort, "3306")
	user := firstNonEmpty(fileCfg.MySQLUser, "root")
	password := firstNonEmpty(getenv("MYSQL_PASSWORD"), "213140")
	database := firstNonEmpty(fileCfg.MySQLDatabase, "vibe_coding")
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&multiStatements=true&loc=Local", user, password, host, port, database)
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

func intOrDefault(value, fallback int) int {
	if value <= 0 {
		return fallback
	}
	return value
}
