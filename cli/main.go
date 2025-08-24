package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("  vbook-cli build   <scriptPath>")
		fmt.Println("  vbook-cli install <scriptPath>")
		fmt.Println("  vbook-cli test    <scriptPath>")
		os.Exit(1)
	}

	cmd := os.Args[1]
	scriptPath, _ := filepath.Abs(os.Args[2])

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
