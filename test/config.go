package test

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	StripeSecret        string `mapstructure:"STRIPE_SECRET"`
	WebHookSecret       string `mapstructure:"WEBHOOK_SECRET"`
	WebHookUrl          string `mapstructure:"WEBHOOK_URL"`
	TerminalId          string `mapstructure:"TERMINAL_ID"`
	FirebaseDatabaseUrl string `mapstructure:"FIREBASE_DATABASE_URL"`
}

func LoadConfig() (*Config, error) {
	config := Config{}

	viper.SetDefault("STRIPE_SECRET", "not_configured")
	viper.SetDefault("WEBHOOK_SECRET", "not_configured")
	viper.SetDefault("WEBHOOK_URL", "not_configured")
	viper.SetDefault("TERMINAL_ID", "not_configured")
	viper.SetDefault("FIREBASE_DATABASE_URL", "not_configured")

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("config file not loaded, using environment variables: %s \n", err)
	}

	log.Println(viper.AllKeys())
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	log.Println(config)
	return &config, nil
}
