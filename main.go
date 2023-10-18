package main

import (
	"fmt"
	"os"

	"github.com/vcokltfre/stop/stop"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <run|build|explain> <file>\n", os.Args[0])
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		data, err := os.ReadFile(os.Args[2])
		if err != nil {
			fmt.Printf("Error reading file: %s\n", err.Error())
			os.Exit(1)
		}

		vm := stop.VM{}
		vm.Run(data)
	case "build":
		data, err := os.ReadFile(os.Args[2])
		if err != nil {
			fmt.Printf("Error reading file: %s\n", err.Error())
			os.Exit(1)
		}

		parsed, err := stop.Parse(string(data))
		if err != nil {
			fmt.Printf("Error parsing file: %s\n", err.Error())
			os.Exit(1)
		}

		err = os.WriteFile(os.Args[2]+".bc", stop.Compile(parsed), 0644)
		if err != nil {
			fmt.Printf("Error writing file: %s\n", err.Error())
			os.Exit(1)
		}
	case "explain":
		data, err := os.ReadFile(os.Args[2])
		if err != nil {
			fmt.Printf("Error reading file: %s\n", err.Error())
			os.Exit(1)
		}

		stop.Explain(data)
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
