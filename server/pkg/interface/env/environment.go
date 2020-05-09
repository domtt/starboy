package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	User, Pass, Name, Host string
}

type EnvConfig struct {
	Port, GithubClientID, GithubClientSecret, WebAppURL, ServerURL string
	Production                                                     bool
	DB                                                             DBConfig
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
	port := getEnv("PORT", "8080")
	dbHost := "localhost"
	webApp := "http://localhost:3000"
	serverURL := "http://localhost:" + port

	if prod {
		dbHost = "db"
		webApp = "" // same as server
		serverURL = "http://starboy.dev"
	}

	config = EnvConfig{
		Production:         prod,
		Port:               port,
		GithubClientID:     mustGet("GITHUB_CLIENT_ID"),
		GithubClientSecret: mustGet("GITHUB_CLIENT_SECRET"),
		DB: DBConfig{
			Name: mustGet("POSTGRES_DB"),
			Pass: mustGet("POSTGRES_PASSWORD"),
			User: mustGet("POSTGRES_USER"),
			Host: dbHost,
		},
		WebAppURL: webApp,
		ServerURL: serverURL,
	}
	return &config
}

func Config() *EnvConfig {
	if len(config.Port) == 0 {
		log.Fatal("Attempted to fetch config before initialization")
	}
	return &config
}
