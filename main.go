package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	// get executable path
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dir := filepath.Dir(exePath)
	// dir = "W:\\Software\\CloudMusic"

	processName := "unblockneteasemusic-win-x64.exe"
	cloudMusicExe := filepath.Join(dir, "cloudmusic.exe")
	unblockExe := filepath.Join(dir, "UnblockNeteaseMusic", processName)

	_ = godotenv.Load(filepath.Join(dir, "UnblockNeteaseMusic", ".env"))
	port := os.Getenv("PORT")
	source := os.Getenv("SOURCE")
	if port == "" {
		port = "2323"
	}
	if source == "" {
		source = "pyncmd,kuwo"
	}
	sources := strings.Split(source, ",")
	// check if process is already running
	exists, err := processExists(processName)
	if err != nil {
		fmt.Println("Error checking process:", err)
		return
	}
	if exists {
		cmd := exec.Command("taskkill", "/F", "/IM", processName)
		_, _ = cmd.CombinedOutput()
	}
	args := []string{"-p", port, "-o"}
	args = append(args, sources...)
	cmd := exec.Command(unblockExe, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	err = cmd.Start()
	if err != nil {
		fmt.Println("Error starting process:", err)
		return
	}

	// start cloudmusic.exe
	cmd = exec.Command(cloudMusicExe)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	cmd.Start()
}

func processExists(processName string) (bool, error) {
	cmd := exec.Command("wmic", "process", "get", "Name")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return false, err
	}
	if bytes.Contains(out.Bytes(), []byte(processName)) {
		return true, nil
	}
	return false, nil
}
