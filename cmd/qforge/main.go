package main

import (
	"context"
	"fmt"
	"os"

	"qforge/internal/cli"
)

func main() {
	ctx := context.Background()
	if err := cli.Run(ctx, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
