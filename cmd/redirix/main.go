package main

import (
	"fmt"
	"os"
)

func main() {
	if err := Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "[Redirix] Error: %v\n", err)
		os.Exit(1)
	}
}
