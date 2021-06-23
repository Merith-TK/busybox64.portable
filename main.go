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
	busyboxBin string
	busybox    string

	configfile string
	dataDir    string
)

func init() {
	// set Absolute System Paths
	configfile, _ = filepath.Abs(strings.TrimSuffix(os.Args[0], ".exe") + ".toml")
	configfile = strings.ReplaceAll(configfile, "\\", "/")
	dataDir, _ = filepath.Abs(strings.TrimSuffix(os.Args[0], ".exe") + ".data")
	dataDir = strings.ReplaceAll(dataDir, "\\", "/")
	// Set Absolute buysbox Paths
	busyboxBin, _ = filepath.Abs(dataDir + "/bin")
	busyboxBin = strings.ReplaceAll(busyboxBin, "\\", "/")
	busybox = busyboxBin + "/busybox.exe"
}
func main() {
	// Just Ensure that the folders are there
	if _, err := os.Stat(busybox); err != nil {
		err := nil
		os.MkdirAll(dataDir+"/bin", 0755)
		os.MkdirAll(dataDir+"/home", 0755)
		os.MkdirAll(dataDir+"/opt", 0755)
		if err != nil {
			log.Fatal("[ERROR] Could not make rootfs directories")
		}
	}
	// If busybox is not found, get it
	if _, err := os.Stat(busybox); err != nil {
		err := fetchFile(busybox, "https://frippery.org/files/busybox/busybox64.exe")
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
	}

	setupConfig()
	setupEnvironment()
	execute(conf.Program, conf.ProgramArgs)
}

func execute(execute string, args string) {
	cmdargs := strings.Split(args, " ")
	cmd := exec.Command(execute, cmdargs...)
	cmd.Dir = filepath.Dir(dataDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin
	fmt.Println("[BusyBox64 Portable] Running "+execute, args)
	cmd.Run()
}

// download litterally any file
func fetchFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
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

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
