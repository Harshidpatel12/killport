package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: killport <port>  |  killport --container <name|id|port>")
		os.Exit(1)
	}

	arg1 := os.Args[1]

	// --container / --containor (typo kept for compat) / -c
	if arg1 == "--container" || arg1 == "--containor" || arg1 == "-c" {
		if len(os.Args) < 3 {
			fmt.Printf("Usage: killport %s <container_name_or_id_or_port>\n", arg1)
			os.Exit(1)
		}
		target := os.Args[2]

		// If target is a valid port number, resolve containers on that port first.
		// ponytail: fall through to treating target as a container name/ID if no containers found.
		if port, err := strconv.Atoi(target); err == nil && port > 0 && port <= 65535 {
			if containers, err := findDockerContainers(port); err == nil && len(containers) > 0 {
				for _, c := range containers {
					fmt.Printf("Killing Docker container '%s' (ID: %s) occupying port %d...\n", c.Name, c.ID, port)
					killContainer(c.ID, c.Name)
				}
				os.Exit(0)
			}
		}

		fmt.Printf("Killing Docker container %s...\n", target)
		killContainer(target, target)
		os.Exit(0)
	}

	// Plain port mode
	port, err := strconv.Atoi(arg1)
	if err != nil || port <= 0 || port > 65535 {
		fmt.Printf("Invalid port: %s. Port must be an integer between 1 and 65535.\n", arg1)
		os.Exit(1)
	}

	// ponytail: docker check is best-effort; ignored if docker isn't available.
	if containers, err := findDockerContainers(port); err == nil && len(containers) > 0 {
		fmt.Printf("Port %d is occupied by Docker container(s):\n", port)
		for _, c := range containers {
			fmt.Printf("  - '%s' (ID: %s)\n", c.Name, c.ID)
		}
		fmt.Println("\nTo kill the container(s), run:")
		for _, c := range containers {
			fmt.Printf("  killport --container %s\n", c.Name)
		}
		os.Exit(0)
	}

	fmt.Printf("Searching for processes on port %d...\n", port)

	// lsof -t: terse PID-only output; non-zero exit = nothing found.
	out, err := exec.Command("lsof", "-t", fmt.Sprintf("-i:%d", port)).Output()
	if err != nil || strings.TrimSpace(string(out)) == "" {
		fmt.Printf("No process found running on port %d.\n", port)
		os.Exit(0)
	}

	pids := strings.Fields(string(out)) // Fields splits on any whitespace and ignores empty tokens
	fmt.Printf("Found %d process(es) on port %d: %s\n", len(pids), port, strings.Join(pids, ", "))

	self := os.Getpid()
	for _, pidStr := range pids {
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			fmt.Printf("Failed to parse PID: %s\n", pidStr)
			continue
		}
		if pid == self {
			fmt.Println("Skipping current process.")
			continue
		}

		fmt.Printf("Killing process %d...\n", pid)
		// os.FindProcess on Linux never errors; the kill call is the real check.
		process, _ := os.FindProcess(pid)
		if err := process.Signal(syscall.SIGKILL); err != nil {
			fmt.Printf("Failed to kill process %d: %v\n", pid, err)
		} else {
			fmt.Printf("Successfully killed process %d.\n", pid)
		}
	}
}

// killContainer runs `docker kill <id>` and exits 1 on failure.
func killContainer(id, displayName string) {
	out, err := exec.Command("docker", "kill", id).CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to kill Docker container %s: %s\nError: %v\n", displayName, strings.TrimSpace(string(out)), err)
		os.Exit(1)
	}
	fmt.Printf("Successfully killed Docker container %s.\n", displayName)
}

type DockerContainer struct {
	ID, Name string
}

func findDockerContainers(port int) ([]DockerContainer, error) {
	if _, err := exec.LookPath("docker"); err != nil {
		return nil, err
	}

	// 1. Get running container IDs
	cmdIds := exec.Command("docker", "ps", "-q")
	idsBytes, err := cmdIds.Output()
	if err != nil {
		return nil, err
	}
	idsStr := strings.TrimSpace(string(idsBytes))
	if idsStr == "" {
		return nil, nil
	}
	ids := strings.Fields(idsStr)

	// 2. Inspect all running containers to get detailed networking info
	args := append([]string{"inspect"}, ids...)
	cmdInspect := exec.Command("docker", args...)
	inspectBytes, err := cmdInspect.Output()
	if err != nil {
		return nil, err
	}

	var rawContainers []struct {
		ID         string `json:"Id"`
		Name       string `json:"Name"`
		HostConfig struct {
			NetworkMode string `json:"NetworkMode"`
		} `json:"HostConfig"`
		Config struct {
			ExposedPorts map[string]interface{} `json:"ExposedPorts"`
		} `json:"Config"`
		NetworkSettings struct {
			Ports map[string][]struct {
				HostPort string `json:"HostPort"`
			} `json:"Ports"`
		} `json:"NetworkSettings"`
	}

	if err := json.Unmarshal(inspectBytes, &rawContainers); err != nil {
		return nil, err
	}

	portStr := strconv.Itoa(port)
	targetTcp := fmt.Sprintf("%d/tcp", port)
	targetUdp := fmt.Sprintf("%d/udp", port)

	var containers []DockerContainer
	for _, c := range rawContainers {
		matched := false

		// Case A: Port binding (standard published port)
		for _, bindings := range c.NetworkSettings.Ports {
			for _, b := range bindings {
				if b.HostPort == portStr {
					matched = true
					break
				}
			}
			if matched {
				break
			}
		}

		// Case B: Host network mode + exposed port
		if !matched && c.HostConfig.NetworkMode == "host" {
			if c.Config.ExposedPorts != nil {
				_, hasTcp := c.Config.ExposedPorts[targetTcp]
				_, hasUdp := c.Config.ExposedPorts[targetUdp]
				if hasTcp || hasUdp {
					matched = true
				}
			}
		}

		if matched {
			name := strings.TrimPrefix(c.Name, "/")
			// Use 12 character short ID
			shortID := c.ID
			if len(shortID) > 12 {
				shortID = shortID[:12]
			}
			containers = append(containers, DockerContainer{ID: shortID, Name: name})
		}
	}
	return containers, nil
}
