package main

// Copyright (c) 2024 Julian Müller (ChaoticByte)
// License: MIT

import (
	"flag"
)

var configFilepath string
var privateKeyFilepath string
var logFlag bool

func ParseCommandline() bool {
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
