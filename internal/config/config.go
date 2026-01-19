package config

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		Postgres PostgresConfig
		HTTP     HTTPConfig
	}
	PostgresConfig struct {
		Username string
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Name     string
		SSLMode  string `mapstructure:"sslmode"`
		Password string
	}
	HTTPConfig struct {
		Host               string        `mapstructure:"host"`
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}
)

func Init(configsDir string) (*Config, error) {
	viper.AddConfigPath(configsDir)
	viper.SetConfigName("main")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}
	setFromEnv(&cfg)
	return &cfg, nil

}
func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("postgres", &cfg.Postgres); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}
	return nil
}
func setFromEnv(cfg *Config) {
	cfg.Postgres.Username = os.Getenv("POSTGRES_USER")
	cfg.Postgres.Name = os.Getenv("POSTGRES_DB")
	cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
}
