package main

import (
	"encoding/json"
	"os"
)

var config = Config{
	ImAddr: "http://43.128.72.19:10002",
}

type Config struct {
	ImAddr string `json:"imAddr"`
}

func initConfig() {
	f, err := os.Open("config.json")
	if err != nil {
		return
	}
	defer f.Close()

	var c Config
	err = json.NewDecoder(f).Decode(&c)
	if err != nil {
		return
	}
	if c.ImAddr != "" {
		config.ImAddr = c.ImAddr
	}
}
