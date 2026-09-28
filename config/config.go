package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Api_key      string `json:"API_KEY"`
	Email        string `json:"email"`
	Token        string
	Amount       int      `json:"amount"`
	Note         string   `json:"note"`
	Date_from    string   `json:"date_from"`
	Date_to      string   `json:"date_to"`
	Names_tutors []string `json:"names_tutors"`
	Id_tutors    []int
	Id_groups    []int
}
func LoadConfig() (*Config, error) {

	filePath := "job.json"

	// 1. Читаем содержимое файла
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// 2. Создаем экземпляр структуры
	var cfg Config

	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
