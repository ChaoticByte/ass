package main

// Copyright (c) 2024 Julian Müller (ChaoticByte)
// License: MIT

import "os"

var Version = "dev"

func main() {
	if ParseCommandline() {
		ParseConfig(configFilepath)
		RunServer()
	} else {
		os.Exit(1)
	}
}
