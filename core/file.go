package core

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func ReadFile(filePath string) []byte {
	if filePath == "" {
		flag.Usage()
		os.Exit(1)
	}
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
	}
	fileInfo, err := os.Stat(absPath)
	// Test if file can be accessed
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "File does not exist: %s\n", absPath)
		} else {
			fmt.Fprintf(os.Stderr, "Error accessing file: %v\n", err)
		}
		os.Exit(1)
	}
	// Test if file is file & not a directory
	if fileInfo.IsDir() {
		fmt.Fprintf(os.Stderr, "Path is a directory, not a file: %s\n", absPath)
		os.Exit(1)
	}
	content, err := os.ReadFile(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	return content
}

func WriteFile(content string, filePath string) {
	if filePath == "" {
		flag.Usage()
		os.Exit(1)
	}
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
		os.Exit(1)
	}
	err = os.WriteFile(absPath, []byte(content), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}
}
