package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

var (
    Port      string
    DBUrl     string
    JwtSecret string
)

func LoadEnv() {
    if err := godotenv.Load(); err != nil {
        log.Println(" .env not found, using default env vars")
    }

    Port = os.Getenv("PORT")
    DBUrl = os.Getenv("DB_URL")
    JwtSecret = os.Getenv("JWT_SECRET")

    if Port == "" || DBUrl == "" || JwtSecret == "" {
        log.Fatal(" Missing required environment variables")
    }
}