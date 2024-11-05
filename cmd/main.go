package main

import (
	"baobud/core"
	"flag"
	"fmt"
	"os"
)

const VERSION = "prealpha-2024-11-05"
const DEBUG = false

func printVersion() {
	fmt.Printf("baobud %s (https://github.com/danielgormly/baobud)\n", VERSION)
}

func main() {
	// TODO: Accept stdout input
	flag.Usage = func() {
		fmt.Println(`Usage: baobud -f <file>
Creates OpenBao/Vault policies from Consul Template templates

Commands:
baobud -f <template file>   Generate policy to stdout
baobud -f <template file> -o <output file>   Generate policy to file
baobud version              Show the version
baobud help                 Show this help message`)
	}
	debugPrint("OS args %v", os.Args[0])

	filePath := flag.String("f", "", "Path to template file")
	outputPath := flag.String("o", "", "Output file (optional)")
	flag.Parse()

	switch {
	case os.Args[1] == "version":
		printVersion()
	case os.Args[1] == "help":
		flag.Usage()
	default:
		if *filePath == "" {
			fmt.Println("Error: -f flag is required")
			flag.Usage()
			os.Exit(1)
		}
		fmt.Println(*outputPath)
		policy := generatePolicy(*filePath)
		if *outputPath != "" {
			writeFile(policy, *outputPath)
		} else {
			fmt.Println(policy)
		}
	}
}

func debugPrint(format string, a ...any) {
	if DEBUG {
		fmt.Printf("# Debug: "+format+"\n", a...)
	}
}

func generatePolicy(filePath string) string {
	debugPrint("Processing file \"%s\"", filePath)
	file := readFile(filePath)
	deps, err := core.Analyze(string(file))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing template: %v\n", err)
		os.Exit(1)
	}
	policy := core.CreateVaultPolicy(deps)
	return policy
}
