package main

//nolint:depguard
import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Logger     LoggerConf
	SQL        SQLConf
	HTTPServer HTTPServerConf
	App        AppConfig
}

type LoggerConf struct {
	Level string `mapstructure:"level" default:"INFO"`
}

type SQLConf struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Database string `mapstructure:"database"`
}

type HTTPServerConf struct {
	Host string `mapstructure:"host" default:"0.0.0.0"`
	Port string `mapstructure:"port" default:"8080"`
}

type AppConfig struct {
	LoginLimitAttempts    int
	PasswordLimitAttempts int
	IPLimitAttempts       int
	Frequency             int
}

func NewConfig(path string) (Config, error) {
	var conf Config
	err := viper.BindEnv("SQL.Host", "sqlHost")
	if err != nil {
		return conf, err
	}

	err = viper.BindEnv("SQL.Port", "sqlPort")
	if err != nil {
		return conf, err
	}

	err = viper.BindEnv("SQL.Database", "sqlDatabase")
	if err != nil {
		return conf, err
	}

	err = viper.BindEnv("App.LoginLimitAttempts", "loginLimit")
	if err != nil {
		return conf, err
	}

	err = viper.BindEnv("App.PasswordLimitAttempts", "passwordLimit")
	if err != nil {
		return conf, err
	}

	err = viper.BindEnv("App.IPLimitAttempts", "ipLimit")
	if err != nil {
		return conf, err
	}

	viper.SetDefault("SQL.Username", "postgres")
	viper.SetDefault("SQL.Password", "password")
	viper.SetDefault("SQL.Host", "0.0.0.0")
	viper.SetDefault("SQL.Port", "5435")
	viper.SetDefault("SQL.Database", "backend")
	viper.SetDefault("App.LoginLimitAttempts", 10)
	viper.SetDefault("App.PasswordLimitAttempts", 100)
	viper.SetDefault("App.IPLimitAttempts", 1000)
	viper.SetDefault("App.Frequency", 60)
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return conf, fmt.Errorf("error while reading config file: %w", err)
	}

	if err := viper.Unmarshal(&conf); err != nil {
		return conf, fmt.Errorf("error while unmarshaling config: %w", err)
	}

	return conf, nil
}
