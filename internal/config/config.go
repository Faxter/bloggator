package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	FILENAME string = ".config/bloggator.json"
)

type Config struct {
	DbUrl       string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

func Read() (Config, error) {
	filepath := getFilePath()
	data, err := os.ReadFile(filepath)
	if err != nil {
		return Config{}, fmt.Errorf("error while reading from file %s", filepath)
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, fmt.Errorf("error while decoding into json")
	}
	return config, nil
}

func (c Config) SetUser(userName string) error {
	c.CurrentUser = userName
	err := write(c)
	if err != nil {
		return err
	}
	return nil
}

func getFilePath() string {
	file, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("could not find user home directory")
		return ""
	}
	file += string(os.PathSeparator)
	file += FILENAME
	return file
}

func write(c Config) error {
	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("could not convert to json")
	}
	err = os.WriteFile(getFilePath(), data, 0644)
	if err != nil {
		return fmt.Errorf("error while writing file")
	}
	return nil
}
