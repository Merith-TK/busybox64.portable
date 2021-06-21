package main

import (
	"fmt"
	"os"
	"strings"
)

func envInit() {
	for k, v := range conf.Environment {

		// dataDirectory
		v = strings.ReplaceAll(v, "{data}", dataDir)

		// Current PWD
		pwd, _ := os.Getwd()
		v = strings.ReplaceAll(v, "{pwd}", pwd)
		os.Setenv(k, v)
		fmt.Println("ENV:", k, "=", v)
	}
}
