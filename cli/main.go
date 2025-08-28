package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args

	switch {
	case len(args) < 2:
		// No arguments → interactive mode
		cmd, scriptPath := promptInteractive()
		runCommand(cmd, scriptPath)

	case len(args) < 3:
		// Invalid usage
		showUsage()
		os.Exit(1)

	default:
		// Normal CLI usage
		cmd := args[1]
		scriptPath, _ := filepath.Abs(args[2])
		runCommand(cmd, scriptPath)
	}
}

func showUsage() {
	fmt.Println("Usage:")
	fmt.Println("  vbook-cli build   <scriptPath>")
	fmt.Println("  vbook-cli install <scriptPath>")
	fmt.Println("  vbook-cli test    <scriptPath>")
}

func promptInteractive() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to vbook-cli")
	fmt.Println("Choose a command:")
	fmt.Println("  1) build")
	fmt.Println("  2) install")
	fmt.Println("  3) test")
	fmt.Print("Enter choice: ")

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var cmd string
	switch choice {
	case "1":
		cmd = "build"
	case "2":
		cmd = "install"
	case "3":
		cmd = "test"
	default:
		fmt.Println("Invalid choice")
		os.Exit(1)
	}

	fmt.Print("Enter script path: ")
	scriptPath, _ := reader.ReadString('\n')
	scriptPath = filepath.Clean(strings.Trim(strings.TrimSpace(scriptPath), `"'`))

	absPath, err := filepath.Abs(scriptPath)
	if err != nil {
		fmt.Println("Invalid path:", err)
		os.Exit(1)
	}

	return cmd, absPath
}

func runCommand(cmd, scriptPath string) {
	switch cmd {
	case "build":
		if err := BuildExtension(scriptPath); err != nil {
			log.Fatal(err)
		}
	case "install":
		if err := InstallExtension(scriptPath); err != nil {
			log.Fatal(err)
		}
	case "test":
		name, content, err := GetFileContent(scriptPath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Testing file:", name)
		if err := TestScript(scriptPath, content); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	default:
		fmt.Println("Unknown command:", cmd)
		os.Exit(1)
	}
}
