package main

import "fmt"

type DbConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

func DefaultDbConfig() DbConfig {
	return DbConfig{
		Host:     "localhost",
		Port:     "5432",
		Database: "mydatabase",
		User:     "myuser",
		Password: "mypassword",
	}
}

func DbConnectionString(c DbConfig) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Database)
}
