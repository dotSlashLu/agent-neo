package main

import (
	"flag"
)

type flags struct {
	configFile string
}

func parseFlags() *flags {
	f := flags{}
	defaultConfigFile := "/etc/agent-neo/config.toml"
	configFile := flag.String("c", defaultConfigFile, "path to config file")

	flag.Parse()

	f.configFile = *configFile
	return &f
}
