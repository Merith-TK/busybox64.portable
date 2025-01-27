package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	// create blank globals for init()
	busybox string

	dataDir    = strings.TrimSuffix(os.Args[0], ".exe") + ".data"
	configfile = dataDir + "/config.toml"
)

func init() {
	// Set Absolute busybox Paths
	datapath, err := filepath.Abs(dataDir)
	if err != nil {
		log.Fatal("[ERROR]: Unable to locate Data Directory, ", err)
	} else {
		dataDir = datapath
	}
	busybox = filepath.Join(dataDir, "/busybox.exe")
}

func main() {
	// If busybox is not found, get it
	if _, err := os.Stat(busybox); err != nil {
		os.Mkdir(dataDir, 0755)
		err := fetchFile(busybox, "https://frippery.org/files/busybox/busybox64u.exe")
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
	}
	if _, err := os.Stat(dataDir + "/start.sh"); err != nil {
		err := os.WriteFile(dataDir+"/start.sh", fsDefaultStart, 0644)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
	}

	if _, err := os.Stat(dataDir + "/config.toml"); err != nil {
		err := os.WriteFile(dataDir+"/config.toml", fsDefaultConfig, 0644)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
	}

	setupConfig()
	setupEnvironment()
	execute(conf.Program, conf.ProgramArgs, conf.WorkingDirectory)
}

func execute(execute string, args string, pwd string) {
	if conf.IsolatedPath {
		os.Setenv("PATH", dataDir)
	} else {
		os.Setenv("PATH", strings.Join([]string{dataDir, os.Getenv("PATH")}, ";"))
	}
	cmdargs := strings.Split(args, " ")
	cmd := exec.Command(execute, cmdargs...)
	cmd.Dir = filepath.Dir(dataDir)
	if pwd != "" {
		cmd.Dir = pwd
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	fmt.Println("[BusyBox64 Portable] Running "+execute, args)
	cmd.Run()
}

// download literally any file
func fetchFile(filepath string, url string) error {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Mimic a browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to the file
	_, err = io.Copy(out, resp.Body)
	return err
}
