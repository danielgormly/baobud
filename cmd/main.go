package main

import (
	"baobud/core"
	"flag"
	"fmt"
	"os"
)

const VERSION = "prealpha-2024-11-05"

var DEBUG = os.Getenv("BAOBUD_DEBUG") == "true"

func printVersion() {
	fmt.Printf("baobud %s (https://github.com/danielgormly/baobud)\n", VERSION)
}

func main() {
	// TODO: Accept stdout input
	flag.Usage = func() {
		fmt.Println(`Usage: baobud -f <file>
Creates OpenBao/Vault policies after evaluating Consul Template templates. Respects BAO_TOKEN & BAO_ADDR environment variables. N.b. Only Vault (and OpenBao) requests are evaluated.

Main command:
baobud <template file>: Generates policy to stdout

Flags:
-o <output file path>: Generates policy to specified file path
-bao-addr <URL>: Address to OpenBao or Vault server
-bao-token <Token>: OpenBao Address or Token
-d: Debug mode (note this may not produce valid hcl in stdout mode)

Other Commands:
baobud version: Show the version
baobud help: Show this help message`)
	}
	debugPrint("OS args %v", os.Args[0])

	filePath := flag.String("f", "", "Path to template file")
	outputPath := flag.String("o", "", "Output file (optional)")
	flag.Parse()

	if len(os.Args) <= 1 {
		os.Exit(1)
	}

	switch {
	case len(os.Args) <= 1:
		flag.Usage()
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
		policy := generatePolicy(*filePath)
		if *outputPath != "" {
			fmt.Printf("Writing policy to %s\n", *outputPath)
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
