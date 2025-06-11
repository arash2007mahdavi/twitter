package configs

import (
	"os"

	"github.com/spf13/viper"
)

func GetConfig() *Config {
	path := getPath(os.Getenv("APP_ENV"))
	v := parseConfig(path, "yml")
	config := loadConfig(v)
	return config
}

func loadConfig(v *viper.Viper) *Config {
	var cfg Config
	err := v.Unmarshal(&cfg)
	if err != nil {
		panic("error in parsing config")
	}
	return &cfg
}

func parseConfig(filename string, filetype string) *viper.Viper {
	v := viper.New()
	v.SetConfigName(filename)
	v.SetConfigType(filetype)
	v.AddConfigPath(".")
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		panic("error in loading config")
	}
	return v
}

func getPath(env string) string {
	if env == "docker" {
		return "../configs/docker-config"
	} else if env == "production" {
		return "../configs/production-config"
	} else {
		return "../configs/development-config"
	}
}