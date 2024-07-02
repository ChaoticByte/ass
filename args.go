package main

// Copyright (c) 2024 Julian Müller (ChaoticByte)
// License: MIT

import (
	"flag"
	"fmt"
)

var configFilepath string
var privateKeyFilepath string
var logFlag bool

func ParseCommandline() bool {
	flag.Usage = func () {
		fmt.Println(`Usage: ass -config <string> -pkey <string> [-log]

  -config <string>   The path to the config file (required)
  -pkey <string>     The path to the private key file (required)
  -log               Enable logging of messages

version:`, Version)
	}
	flag.StringVar(&configFilepath, "config", "", "The path to the config file (required)")
	flag.StringVar(&privateKeyFilepath, "pkey", "", "The path to the private key file (required)")
	flag.BoolVar(&logFlag, "log", false, "Enable logging of messages")
	flag.Parse()
	missing := false
	if configFilepath == "" {
		missing = true
	}
	if privateKeyFilepath == "" {
		missing = true
	}
	if missing {
		flag.Usage()
	}
	return !missing
}
