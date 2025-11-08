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
	config := NewConfig()
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("command is required")
		return
	}
	config.DB.Driver = "postgres"
	config.DB.Host = env.GetString("DB_HOST", "localhost")
	config.DB.Port = env.GetString("DB_PORT", "5432")
	config.DB.User = env.GetString("DB_USER", "postgres")
	config.DB.Pass = env.GetString("DB_PASSWORD", "")
	config.DB.DB = env.GetString("DB_NAME", "book_collection")
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
