package config

import "os"

type Config struct {
    HTTPPort string
    DBConn   string
    RedisAddr string
}

func Load() *Config {
    return &Config{
        HTTPPort:  getEnv("HTTP_PORT", "8080"),
        DBConn:    getEnv("DB_CONN", "postgres://postgres:postgres@localhost:5432/gorest?sslmode=disable"),
        RedisAddr: getEnv("REDIS_ADDR", "localhost:6379"),
    }
}

func getEnv(key, defaultVal string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return defaultVal
}