package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	toml "github.com/pelletier/go-toml"
)

var (
	conf config
)

type config struct {
	Program          string            `toml:"program"`
	ProgramArgs      string            `toml:"programArgs"`
	WorkingDirectory string            `toml:"workingDirectory"`
	IsolatedPath     bool              `toml:"isolatedPath"`
	Environment      map[string]string `toml:"environment"`
}

func setupConfig() {
	// Check if file is readable, if not, make the file
	str, fileErr := os.ReadFile(configfile)
	if fileErr != nil {
		// If the file is not found, create it from the embedded content
		err := os.WriteFile(configfile, fsDefaultConfig, 0644)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
		str = fsDefaultConfig // Use the embedded default config
	}

	// Unmarshal the config data into the struct
	err := toml.Unmarshal(str, &conf)
	if err != nil {
		log.Fatalln(err)
	}
}

// Setup the Environment part of the config
func setupEnvironment() {
	// Variables for env replacement
	drivePath, _ := filepath.Abs("/")
	drivePath = filepath.ToSlash(drivePath)
	drivePath = strings.TrimSuffix(drivePath, "/")
	configEnvReplace := map[string]string{
		"{data}":  dataDir,
		"{drive}": drivePath,
	}

	// Replace Normal Config options
	for key, value := range configEnvReplace {
		if strings.Contains(conf.Program, key) {
			conf.Program = filepath.ToSlash(strings.ReplaceAll(conf.Program, key, value))
		}
		if strings.Contains(conf.ProgramArgs, key) {
			conf.ProgramArgs = filepath.ToSlash(strings.ReplaceAll(conf.ProgramArgs, key, value))
		}
		if strings.Contains(conf.WorkingDirectory, key) {
			conf.WorkingDirectory = filepath.ToSlash(strings.ReplaceAll(conf.WorkingDirectory, key, value))
		}
	}

	// Replace Environment Variables
	for k, v := range conf.Environment {
		for key, value := range configEnvReplace {
			if strings.Contains(v, key) {
				v = strings.ReplaceAll(v, key, value)
				v = filepath.ToSlash(v)
			}
		}
		os.Setenv(k, v)
	}
}
