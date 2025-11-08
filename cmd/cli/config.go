package main

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

func NewConfig(env string, db dbConfig) *Config {
	return &Config{
		ENV: env,
		DB:  db,
	}
}
