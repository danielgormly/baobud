package main

import (
	"baobud/core"
	"flag"
	"fmt"
	"os"
)

const VERSION = "prealpha-2024-11-05"
const DEBUG = true

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
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}
	fileCmd := flag.NewFlagSet("baobud", flag.ExitOnError)
	filePath := fileCmd.String("f", "", "Path to template file")

	switch os.Args[1] {
	case "-f":
		fileCmd.Parse(os.Args[1:])
		if *filePath == "" {
			fmt.Println("Error: -f requires a file argument")
			os.Exit(1)
		}
		handleFile(*filePath)
	case "version":
		printVersion()
	case "help":
		flag.Usage()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		flag.Usage()
		os.Exit(1)
	}
}

func debugPrint(format string, a ...any) {
	if DEBUG {
		fmt.Printf("# Debug: "+format+"\n", a...)
	}
}

func handleFile(filePath string) {
	debugPrint("Processing file \"%s\"", filePath)
	file := readFile(filePath)
	deps, err := core.Analyze(string(file))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing template: %v\n", err)
		os.Exit(1)
	}
	policy := core.CreateVaultPolicy(deps)
	fmt.Println(policy)
}
