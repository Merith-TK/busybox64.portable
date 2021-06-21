package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var (
	busybox = dataDir + "/.busybox/busybox64.exe"
)

func main() {
	if _, err := os.Stat(busybox); err != nil {
		os.MkdirAll(dataDir+"/.busybox", 0755)
		os.MkdirAll(dataDir+"/.bin", 0755)
	}
	if _, err := os.Stat(busybox); err != nil {
		err := fetchFile(busybox, "https://frippery.org/files/busybox/busybox64.exe")
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
	}

	configErr := setupConfig()
	if configErr != nil {
		fmt.Println("Could not init config")
		log.Println(configErr)
	}

	envInit()

	exec.Command(busybox, "--install", dataBin).Run()
	cmd := exec.Command(busybox, "bash")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin
	cmdErr := cmd.Run()
	if cmdErr != nil {
		log.Println(cmdErr)
	}
}

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
