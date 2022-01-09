package internal

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func LoadENVFiles() error {
	for _, filename := range [...]string{".env"} {
		if err := godotenv.Load(filename); err != nil && !strings.Contains(err.Error(), "no such file or directory") {
			return Wrap(err, "godotenv.Load")
		}
	}

	return nil
}

func GetBrokers() []string {
	var brokers []string

	if envVal := os.Getenv("BROKERS"); envVal != "" {
		brokers = strings.Split(envVal, ",")
	}

	return brokers
}
