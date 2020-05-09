package env

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	User, Pass, Name, Host string
}

type EnvConfig struct {
	Port, GithubClientID, GithubClientSecret string
	Production                               bool
	DB                                       DBConfig
}

var config EnvConfig

func mustGet(key string) string {
	v := os.Getenv(key)
	if len(v) == 0 {
		log.Fatal("Failed to get env variable: " + key)
	}
	return v
}

func getEnv(key, defaultValue string) string {
	v := os.Getenv(key)
	if len(v) == 0 {
		return defaultValue
	}
	return v
}

func hasEnv(key string) bool {
	_, ok := os.LookupEnv(key)
	return ok
}

func Load() *EnvConfig {
	godotenv.Load()
	prod := hasEnv("PRODUCTION")
	fmt.Println(os.Getenv("PRODUCTION"))
	dbHost := "locahost"
	if prod {
		dbHost = "db"
	}
	config = EnvConfig{
		Port:               getEnv("PORT", "8080"),
		GithubClientID:     mustGet("GITHUB_CLIENT_ID"),
		GithubClientSecret: mustGet("GITHUB_CLIENT_SECRET"),
		DB: DBConfig{
			Name: mustGet("POSTGRES_DB"),
			Pass: mustGet("POSTGRES_PASSWORD"),
			User: mustGet("POSTGRES_USER"),
			Host: dbHost,
		},
		Production: prod,
	}
	return &config
}

func Config() *EnvConfig {
	if len(config.Port) == 0 {
		log.Fatal("Attempted to fetch config before initialization")
	}
	return &config
}
