package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	AppEnv   string `mapstructure:"APP_ENV"`
	HTTPPort string `mapstructure:"HTTP_PORT"`
	MongoURI string `mapstructure:"MONGO_URI"`
	MongoDB  string `mapstructure:"MONGO_DB"`
}

func Load() (*Config, error) {
	viper.SetDefault("APP_ENV", "local")
	viper.SetDefault("HTTP_PORT", "8080")
	viper.SetDefault("MONGO_URI", "mongodb://localhost:27017")
	viper.SetDefault("MONGO_DB", "stream_rec_db")

	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
