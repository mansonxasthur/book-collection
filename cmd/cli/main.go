package main

import (
	"fmt"
	"os"

	"github.com/mansonxasthur/book-collection/pkg/env"
)

var commands = []Command{
	NewMigrationCommand(),
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("command is required")
		return
	}

	config := getConfig()

	for _, cmd := range commands {
		if cmd.Signature() == args[0] {
			if err := cmd.Execute(config, args[1:]); err != nil {
				fmt.Println(err)
			}
			return
		}
	}
	fmt.Println("Command not found")
}

func getConfig() *Config {
	environment := env.GetString("APP_ENV", envDev)
	dbConfig := dbConfig{
		Driver: "postgres",
		Host:   env.GetString("DB_HOST", "localhost"),
		Port:   env.GetString("DB_PORT", "5432"),
		User:   env.GetString("DB_USER", "postgres"),
		Pass:   env.GetString("DB_PASSWORD", ""),
		DB:     env.GetString("DB_NAME", "book_collection"),
	}
	config := NewConfig(environment, dbConfig)
	return config
}
