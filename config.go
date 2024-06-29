package main

// Copyright (c) 2024 Julian Müller (ChaoticByte)
// License: MIT

import (
	"log"
	"os"

	"github.com/goccy/go-yaml"
)

var config struct {
	Host string `yaml:"host"`
	Port int `yaml:"port"`
	Clients map[string]string `yaml:"clients"`
}

func ParseConfig(filepath string) {
	data, err := os.ReadFile(filepath)
	if err != nil { log.Fatal(err) }
	err = yaml.Unmarshal(data, &config)
	if err != nil { log.Fatal(err) }
}
