package main

import "github.com/mansonxasthur/book-collection/pkg/env"

type dbConfig struct {
	Driver string
	Host   string
	Port   string
	User   string
	Pass   string
	DB     string
}

const postgresDriver = "postgres"
const envDev = "development"

type Config struct {
	ENV string
	DB  dbConfig
}

func NewConfig() *Config {
	environment := env.GetString("APP_ENV", envDev)
	return &Config{
		ENV: environment,
		DB:  dbConfig{},
	}
}
