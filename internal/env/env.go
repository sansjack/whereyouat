package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	TCPAddress          string
	HTTPAddress         string
	MMDB_GITHUB_API_URL string
	DB_DIR              string
	DB_FILENAME         string
	TAG_FILE            string
}

var cfg *Config

func Load() error {
	file, err := os.Open(".env")
	if err != nil {
		return fmt.Errorf("failed to open .env file: %w", err)
	}
	defer file.Close()

	envMap := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, `"'`)

		envMap[key] = value
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading .env file: %w", err)
	}

	cfg = &Config{
		TCPAddress:          getEnv(envMap, "TCP_ADDRESS", ":1234"),
		HTTPAddress:         getEnv(envMap, "HTTP_ADDRESS", ":8080"),
		MMDB_GITHUB_API_URL: getEnv(envMap, "MMDB_GITHUB_API_URL", "https://api.github.com/repos/yourusername/yourrepo/releases/latest"),
		DB_DIR:              getEnv(envMap, "DB_DIR", "data"),
		DB_FILENAME:         getEnv(envMap, "DB_FILENAME", "GeoLite2-Country.mmdb"),
		TAG_FILE:            getEnv(envMap, "TAG_FILE", "db_version.txt"),
	}

	return nil
}

// unsure if this is bad practice
func Get() *Config {
	if cfg == nil {
		panic("env.Load() must be called before env.Get()")
	}
	return cfg
}

func getEnv(envMap map[string]string, key, defaultValue string) string {
	if value, ok := envMap[key]; ok {
		return value
	}
	return defaultValue
}
