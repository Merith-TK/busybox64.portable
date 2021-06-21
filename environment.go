package main

import (
	"fmt"
	"os"
	"strings"
)

func envInit() {
	envValueReplace := map[string]string{
		"{data}": dataDir,
	}
	for k, v := range conf.Environment {
		for key, value := range envValueReplace {
			if strings.Contains(v, key) {
				v = strings.ReplaceAll(v, key, value)
			}
		}
		os.Setenv(k, v)
		fmt.Println("ENV:", k, "=", v)
	}
}
