package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Name     string
	Username string
	Password string
	Port     uint
	Host     string
	LogFile  string
}

type LoggerConfig struct {
	FileName string
	MaxSize  int
	Level    string
}

type AuthConfig struct {
	OAuth2Key         string
	OAuth2Secret      string
	RedirectURL       string
	TokenHourLifeSpan string
	JWTSecret         string
}

type Config struct {
	// PROD, DEV, DOCKER
	AppEnv string

	// Server Config
	Host string
	Port uint

	// Frontend URL
	FrontendURL string

	// Logger Config
	Log LoggerConfig

	// Database Config
	Db DatabaseConfig

	// Auth Config
	Auth AuthConfig
}

// All configurations
var allConfigurations = struct {

	// Configuration for app environment : DEV
	Dev Config

	// Configuration for app environment : PROD
	Prod Config

	// Configuration for app environment : DOCKER
	Docker Config
}{}

var currentConfig Config

func GetConfig() Config {
	return currentConfig
}

func InitConfig() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	appEnv := os.Getenv("APP_ENV")

	// Load JSON config
	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error loading config.json file")
	}
	defer configFile.Close()

	// Parse JSON config
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&allConfigurations)

	if err != nil {
		fmt.Println("Error decoding config.json file")
	}

	// Set current config
	switch appEnv {
	case "DEV":
		currentConfig = allConfigurations.Dev
	case "PROD":
		currentConfig = allConfigurations.Prod
	case "DOCKER":
		currentConfig = allConfigurations.Docker
	default:
		panic("Error setting current config")
	}

	fmt.Println("Config loaded successfully")
}
