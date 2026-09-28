package configs

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Session  SessionConfig
	Internal InternalConfig
}

type AppConfig struct {
	Port string
	Name string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type SessionConfig struct {
	Secret          string
	ExpirationHours int `mapstructure:"expiration_hours"`
}

// InternalConfig contains credentials used only for service-to-service calls.
type InternalConfig struct {
	ServiceSecret string `mapstructure:"service_secret"`
}

var Cfg *Config

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("❌ Error reading config file: %s", err)
	}

	Cfg = &Config{}
	if err := viper.Unmarshal(Cfg); err != nil {
		log.Fatalf("❌ Unable to decode into struct: %s", err)
	}
	log.Println("✅ Config loaded successfully")
}
