package main

import (
	"flag"
	"fmt"
	"os"
)

const VERSION = "prealpha-2024-11-04"

func printHelp() {
	fmt.Println(`Usage: baobud -f <file>
Creates OpenBao/Vault policies from Consul Template templates
-- (っ˘ڡ˘ς) --
baobud -f <template file> Creates OpenBao/Vault policies from Consul Template templates
baobud version Show the version
baobud help Show this help message`)
}

func printVersion() {
	fmt.Printf("baobud %s (https://github.com/danielgormly/baobud)\n", VERSION)
}

func main() {
	// TODO: Accept stdout input
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}
	fileCmd := flag.NewFlagSet("baobud", flag.ExitOnError)
	fileName := fileCmd.String("f", "", "Template file to generate policy")

	switch os.Args[1] {
	case "-f":
		fileCmd.Parse(os.Args[1:])
		if *fileName == "" {
			fmt.Println("Error: -f requires a file argument")
			os.Exit(1)
		}
		handleFile(*fileName)
	case "version":
		printVersion()
	case "help":
		printHelp()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printHelp()
		os.Exit(1)
	}
}

func handleFile(filename string) {
	fmt.Printf("Processing file \"%s\"", filename)
}
