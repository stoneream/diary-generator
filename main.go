package main

import (
	"context"
	"flag"
	"github.com/google/subcommands"
	"os"
)

func main() {
	subcommands.Register(&InitCmd{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
