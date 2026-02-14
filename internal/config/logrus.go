package config

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	log := logrus.New()

	level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	log.SetLevel(level)
	log.SetFormatter(&logrus.JSONFormatter{})

	return log
}
