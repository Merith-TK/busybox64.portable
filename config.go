package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	toml "github.com/pelletier/go-toml"
)

var (
	conf          config
	defaultConfig = `# The program to run, can be defined as a direct
	# file path, or as an EXE in the $PATH variable
	program = "busybox.exe"
	# Arguments for the program
	programArgs = "bash"

	# Wether to run as admin or not
	runAsAdmin = false
	
	# Global Environment variables to be set
	# to aid with portability,
	[environment]
		# Please do not edit these values unless you
		# know what you are doing
		APPDATA = "{data}/opt"
		LOCALAPPDATA = "{data}/opt"
		HOME = "{data}/home"
		USERPROFILE = "{data}/home"
	
		# PWD is used to tell the program where to
		# run inside of, this default is good for
		# most things, but if you need to edit it
		# feel free to
		PWD = "{data}/home"
		
		# This tells it where to find programs in 
		# PATH, feel free to edit, but leave the
		# default value for busybox
		PATH = "{data}/bin"
`
)

type config struct {
	Program     string            `toml:"program"`
	ProgramArgs string            `toml:"programArgs"`
	RunAsAdmin  bool              `toml:"runAsAdmin"`
	Environment map[string]string `toml:"environment"`
}

func setupConfig() error {
	// Check if file is readable, if not, make the file
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

// Setup the Environment part of the config,
func setupEnvironment() string {
	// Variables for env replacement
	drivePath, _ := filepath.Abs("/")
	drivePath = filepath.ToSlash(drivePath)
	drivePath = strings.TrimSuffix(drivePath, "/")
	var output string
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
	}

	// Replace Environment Variables
	for k, v := range conf.Environment {
		for key, value := range configEnvReplace {
			if strings.Contains(v, key) {
				v = strings.ReplaceAll(v, key, value)
				v = filepath.ToSlash(v)
			}
			// Set Working Directory for programs that are
			// picky about it, usually not needed so disabled
			// in default config
			if k == "PWD" {
				output = v
			}
		}
		os.Setenv(k, v)
	}

	return output
}
