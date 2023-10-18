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
		fmt.Println("Not implemented yet")
		os.Exit(1)
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

		os.Remove(os.Args[2] + ".bc")

		file, err := os.Create(os.Args[2] + ".bc")
		if err != nil {
			fmt.Printf("Error creating file: %s\n", err.Error())
			os.Exit(1)
		}

		for _, val := range parsed {
			_, err := file.Write(val.Emit())
			if err != nil {
				fmt.Printf("Error writing to file: %s\n", err.Error())
				os.Exit(1)
			}
		}

		file.Close()
	case "explain":
		fmt.Println("Not implemented yet")
		os.Exit(1)
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
