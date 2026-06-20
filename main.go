package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: killport <port>")
		os.Exit(1)
	}

	portStr := os.Args[1]
	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 {
		fmt.Printf("Invalid port: %s. Port must be a positive integer.\n", portStr)
		os.Exit(1)
	}

	fmt.Printf("Searching for processes on port %d...\n", port)

	// Run lsof to get the PID(s) using the port
	// -t: terse output (PID only)
	// -i: select IPv[46] files
	cmd := exec.Command("lsof", "-t", fmt.Sprintf("-i:%d", port))
	outputBytes, err := cmd.Output()
	if err != nil {
		// If lsof exits with a non-zero status (e.g. no process found), we handle it
		fmt.Printf("No process found running on port %d.\n", port)
		os.Exit(0)
	}

	output := strings.TrimSpace(string(outputBytes))
	if output == "" {
		fmt.Printf("No process found running on port %d.\n", port)
		os.Exit(0)
	}

	pids := strings.Split(output, "\n")
	fmt.Printf("Found %d process(es) on port %d: %s\n", len(pids), port, strings.Join(pids, ", "))

	for _, pidStr := range pids {
		pidStr = strings.TrimSpace(pidStr)
		if pidStr == "" {
			continue
		}
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			fmt.Printf("Failed to parse PID: %s\n", pidStr)
			continue
		}

		// Don't accidentally kill the current process
		if pid == os.Getpid() {
			fmt.Println("Skipping killing current process.")
			continue
		}

		fmt.Printf("Killing process with PID %d...\n", pid)
		process, err := os.FindProcess(pid)
		if err != nil {
			fmt.Printf("Failed to find process with PID %d: %v\n", pid, err)
			continue
		}

		err = process.Signal(syscall.SIGKILL)
		if err != nil {
			fmt.Printf("Failed to kill process with PID %d: %v\n", pid, err)
		} else {
			fmt.Printf("Successfully killed process %d.\n", pid)
		}
	}
}
