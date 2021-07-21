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
	"syscall"

	"golang.org/x/sys/windows"
)

var (
	// create blank globals for init()
	busyboxBin string
	busybox    string

	configfile = strings.TrimSuffix(os.Args[0], ".exe") + ".toml"
	dataDir    = strings.TrimSuffix(os.Args[0], ".exe") + ".data"
	pwd        string
)

func init() {
	// Set Absolute buysbox Paths
	datapath, err := filepath.Abs(dataDir)
	if err != nil {
		log.Fatal("[ERROR]: Unable to locate Data Directory, ", err)
	} else {
		dataDir = datapath
	}
	busyboxBin = filepath.ToSlash(datapath + "/bin")
	busybox = filepath.Join(busyboxBin, "/busybox.exe")
}

func main() {
	// Just Ensure that the folders are there
	if _, err := os.Stat(busybox); err != nil {
		var fileerr error
		os.MkdirAll(dataDir+"/bin", 0755)
		os.MkdirAll(dataDir+"/home", 0755)
		os.MkdirAll(dataDir+"/opt", 0755)
		if fileerr != nil {
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
	pwd = setupEnvironment()
	if conf.RunAsAdmin {
		_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
		if err != nil {
			executeAsAdmin()
			println("ELEVATING")
			os.Exit(1)
		}
	}
	execute(conf.Program, conf.ProgramArgs, pwd)
}

func execute(execute string, args string, pwd string) {
	cmdargs := strings.Split(args, " ")
	cmd := exec.Command(execute, cmdargs...)
	cmd.Dir = filepath.Dir(dataDir)
	if pwd != "" {
		cmd.Dir = pwd
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin
	fmt.Println("[BusyBox64 Portable] Running "+execute, args)
	cmd.Run()
}

func executeAsAdmin() {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	var showCmd int32 = 1 //SW_NORMAL

	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		fmt.Println(err)
	}
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
