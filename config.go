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
	dataBin       = dataDir + "/.busybox"
	defaultConfig = `
	launcher = "busybox64.exe"
	launcherArgs = "bash"
	[environment]
	  APPDATA = "{data}/"
	  HOME = "{data}/"
	  PATH = "{data}/.busybox;{data}/.bin
`
)

type config struct {
	Input        string            `toml:"input"`
	LauncherArgs string            `toml:"launcherArgs"`
	Environment  map[string]string `toml:"environment"`
}

func setupConfig() error {
	str, err := os.ReadFile(configfile)
	if err != nil {
		f, err := os.OpenFile(configfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		if _, err := f.WriteString(defaultConfig); err != nil {
			log.Println(err)
		}
		str, _ = os.ReadFile(configfile)
		err = nil
	}
	_ = toml.Unmarshal([]byte(str), &conf)
	return err
}
