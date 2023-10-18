package main

import (
	"fmt"
	"os"
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
		fmt.Println("Not implemented yet")
		os.Exit(1)
	case "explain":
		fmt.Println("Not implemented yet")
		os.Exit(1)
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
