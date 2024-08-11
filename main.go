package main

import (
	"context"
	"flag"
	"github.com/google/subcommands"
	"os"
)

func main() {
	subcommands.Register(&InitCmd{}, "")
	subcommands.Register(&ArchiveCmd{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
