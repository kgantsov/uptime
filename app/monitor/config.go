package monitor

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/kgantsov/uptime/app/model"
)

type Config struct {
	Services []model.Service `json:"services"`
}

func ReadConfig(configFilePath string) (*Config, error) {
	configFile, err := os.Open(configFilePath)
	if err != nil {
		fmt.Println(err)
		return &Config{}, err
	}

	fmt.Println("Successfully Opened jsonFile")

	defer configFile.Close()

	configBytes, err := io.ReadAll(configFile)
	if err != nil {
		fmt.Printf("Error reading %s %s \n", configFilePath, err)
		return &Config{}, err
	}

	var config Config

	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		fmt.Printf("Error unmarshaling %s %s \n", configFilePath, err)
		return &Config{}, err
	}

	fmt.Printf("Found %d services to monitor\n", len(config.Services))

	return &config, nil
}
