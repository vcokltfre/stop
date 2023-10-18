package main

import (
	"fmt"
	"os"

	"github.com/vcokltfre/stop/stop"
)

func build(file string) {
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err.Error())
		os.Exit(1)
	}

	parsed, err := stop.Parse(string(data))
	if err != nil {
		fmt.Printf("Error parsing file: %s\n", err.Error())
		os.Exit(1)
	}

	err = os.WriteFile(file+".bc", stop.Compile(parsed), 0644)
	if err != nil {
		fmt.Printf("Error writing file: %s\n", err.Error())
		os.Exit(1)
	}
}

func run(file string) {
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err.Error())
		os.Exit(1)
	}

	vm := stop.VM{}
	vm.Run(data)
}

func explain(file string) {
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err.Error())
		os.Exit(1)
	}

	stop.Explain(data)
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <build|run|explain> <file>\n", os.Args[0])
		os.Exit(1)
	}

	switch os.Args[1] {
	case "build":
		build(os.Args[2])
	case "run":
		if os.Getenv("STOP_DEV") == "1" {
			build(os.Args[2])
			run(os.Args[2] + ".bc")
			os.Exit(0)
		}
		run(os.Args[2])
	case "explain":
		if os.Getenv("STOP_DEV") == "1" {
			build(os.Args[2])
			explain(os.Args[2] + ".bc")
			os.Exit(0)
		}
		explain(os.Args[2])
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
