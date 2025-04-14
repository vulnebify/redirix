package main

import (
	"context"
	"fmt"
	"os"

	"github.com/vulnebify/redirix/internal/app"
)

func main() {
	if err := app.Run(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "[Redirix] Error: %v\n", err)
		os.Exit(1)
	}
}
