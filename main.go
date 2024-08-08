package main

import (
	"context"
	"diary-generator/cmd/init"
	"flag"
	"os"

	"github.com/google/subcommands"
)

func main() {
	subcommands.Register(&init.InitCmd{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
