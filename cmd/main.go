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
Creates OpenBao/Vault policies from Consul Template templates. Respects BAO_TOKEN & BAO_ADDR environment variables. N.b. Consul & Nomad requests are not currently supported. PRs welcome.

Main command:
baobud <flags> <template file>: Generates policy to stdout

Flags:
-o <output file path>: Generates policy to specified file path
-d: Debug mode (note this may not produce valid hcl in stdout mode)
--bao-addr <URL>: Address to OpenBao or Vault server
--bao-token <Token>: OpenBao Address or Token

Other Commands:
baobud version: Show the version
baobud help: Show this help message`)
	}
	debugPrint("OS args %v", os.Args[0])

	debugFlag := flag.Bool("d", false, "Debug mode (optional)")
	outputPath := flag.String("o", "", "Output file (optional)")
	baoAddr := flag.String("bao-addr", "", "Output file (optional)")
	baoToken := flag.String("bao-token", "", "Output file (optional)")
	flag.Parse()
	args := flag.Args()

	if *debugFlag {
		DEBUG = true
		debugPrint("Debug mode enabled")
	}

	if len(os.Args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	switch args[0] {
	case "version":
		printVersion()
	case "help":
		flag.Usage()
	default:
		config := core.BaobudConfig{
			BaoAddress: func() string {
				if *baoAddr == "" {
					return os.Getenv("BAO_ADDR")
				}
				return *baoAddr
			}(),
			BaoToken: func() string {
				if *baoToken == "" {
					return os.Getenv("BAO_TOKEN")
				}
				return *baoToken
			}(),
		}
		policy := generatePolicy(args[0], config)
		if *outputPath != "" {
			fmt.Printf("Writing policy to %s\n", *outputPath)
			core.WriteFile(policy, *outputPath)
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

func generatePolicy(filePath string, config core.BaobudConfig) string {
	debugPrint("Processing file \"%s\"", filePath)
	file := core.ReadFile(filePath)
	deps, err := core.Analyze(string(file), config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing template: %v\n", err)
		os.Exit(1)
	}
	policy := core.CreateVaultPolicy(deps)
	return policy
}
