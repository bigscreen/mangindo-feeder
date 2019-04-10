package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	port               int
	logLevel           string
	environment        string
	redisHost          string
	redisPort          int
	redisPool          int
	workerRedisAddress string
}

var appConfig *Config

func Load() {
	viper.SetDefault("APP_PORT", "3000")
	viper.SetDefault("LOG_LEVEL", "debug")
	viper.AutomaticEnv()

	viper.SetConfigName("application")
	viper.AddConfigPath("./")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")
	viper.SetConfigType("yaml")

	viper.ReadInConfig()

	appConfig = &Config{
		port:               getIntOrPanic("APP_PORT"),
		logLevel:           fatalGetString("LOG_LEVEL"),
		environment:        fatalGetString("ENVIRONMENT"),
		redisHost:          fatalGetString("REDIS_HOST"),
		redisPort:          getIntOrPanic("REDIS_PORT"),
		redisPool:          getIntOrPanic("REDIS_POOL"),
		workerRedisAddress: fatalGetString("WORKER_REDIS_ADDRESS"),
	}
}

func Port() int {
	return appConfig.port
}

func LogLevel() string {
	return appConfig.logLevel
}

func Environment() string {
	return appConfig.environment
}

func RedisHost() string {
	return appConfig.redisHost
}

func RedisPort() int {
	return appConfig.redisPort
}

func RedisPool() int {
	return appConfig.redisPool
}

func WorkerRedisAddress() string {
	return appConfig.workerRedisAddress
}
