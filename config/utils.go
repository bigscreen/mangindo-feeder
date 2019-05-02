package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

func getIntOrPanic(key string) int {
	checkKey(key)
	v, err := strconv.Atoi(fatalGetString(key))
	if err != nil {
		v, err = strconv.Atoi(os.Getenv(key))
	}
	panicIfErrorForKey(err, key)
	return v
}

func fatalGetString(key string) string {
	checkKey(key)
	value := os.Getenv(key)
	if value == "" {
		value = viper.GetString(key)
	}
	return value
}

func fatalGetStringArray(key, sep string) []string {
	checkKey(key)
	value := fatalGetString(key)
	return strings.Split(value, sep)
}

func checkKey(key string) {
	if !viper.IsSet(key) && os.Getenv(key) == "" {
		log.Fatalf("%s key is not set", key)
	}
}

func panicIfErrorForKey(err error, key string) {
	if err != nil {
		log.Fatalf("Could not parse key: %s, Error: %s", key, err)
	}
}
