package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	UrlSeller       string `json:"urlSeller"`
	DeliveryAddress struct {
		City   string `json:"city"`
		Street string `json:"street"`
	} `json:"deliveryAddress"`
	Categories   []string `json:"categories"`
	Proxy        string   `json:"proxy"`
	Headless     bool     `json:"headless"`
	Port         int      `json:"port"`
	ChromeDriver string   `json:"chromedriver"`
	Export       string   `json:"export"`
}

func ReadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
