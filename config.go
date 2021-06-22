package main

import (
	"log"
	"os"
	"strings"

	toml "github.com/pelletier/go-toml"
)

var (
	conf          config
	configfile    = strings.TrimSuffix(os.Args[0], ".exe") + ".toml"
	dataDir       = strings.TrimSuffix(os.Args[0], ".exe") + ".data"
	defaultConfig = `
program = "bash"
programArgs = ""
[environment]
  APPDATA = "{data}/opt"
  LOCALAPPDATA = "{data}/opt"
  HOME = "{data}/home"
  PATH = "{data}/bin;{data}/bin/busybox"
`
)

type config struct {
	Program     string            `toml:"program"`
	ProgramArgs string            `toml:"programArgs"`
	Environment map[string]string `toml:"environment"`
}

func setupConfig() error {
	str, fileErr := os.ReadFile(configfile)
	if fileErr != nil {
		f, err := os.OpenFile(configfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		if _, err := f.WriteString(defaultConfig); err != nil {
			log.Println(err)
		}
		str, fileErr = os.ReadFile(configfile)
	}
	_ = toml.Unmarshal([]byte(str), &conf)
	return fileErr
}
