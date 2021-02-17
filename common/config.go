package common

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

var Instance Config

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Bot    BotConfig    `mapstructure:"bot"`
}

type ServerConfig struct {
	Production    bool   `mapstructure:"production"`
	MongoURL      string `mapstructure:"mongo_url"`
	MongoUsername string `mapstructure:"mongo_username"`
	MongoPassword string `mapstructure:"mongo_password"`
	RedisURL      string `mapstructure:"redis_url"`
}

type BotConfig struct {
	Token string `mapstructure:"token"`
}

func Init() error {
	logrus.Info("Reading and initializing configurations...")
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("common")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	bindEnvironments()
	if err := viper.ReadInConfig(); err != nil {
		logrus.Info(viper.AllSettings())
	}

	if err := viper.Unmarshal(&Instance); err != nil {
		return err
	}
	logrus.Info("Production:", Instance.Server.Production)
	return nil
}

func bindEnvironments() {
	_ = viper.BindEnv("server.production")
	_ = viper.BindEnv("server.mongo_url")
	_ = viper.BindEnv("server.mongo_username")
	_ = viper.BindEnv("server.mongo_password")
	_ = viper.BindEnv("server.redis_url")
	_ = viper.BindEnv("bot.token")
}
