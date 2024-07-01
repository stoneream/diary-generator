package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
)

type initCmd struct {
	config string
}

func (*initCmd) Name() string     { return "init" }
func (*initCmd) Synopsis() string { return "Initialize a diary" }
func (*initCmd) Usage() string {
	return `init:
	Initialize a diary.
`
}
func (p *initCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.config, "config", "", "config file")
}
func (*initCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// TODO implement
	return subcommands.ExitSuccess
}

type archiveCmd struct {
	config string
}

func (*archiveCmd) Name() string     { return "archive" }
func (*archiveCmd) Synopsis() string { return "Archive a diary" }
func (*archiveCmd) Usage() string {
	return `archive:
	Archive a diary.
`
}
func (p *archiveCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.config, "config", "", "config file")
}
func (*archiveCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// TODO implement
	return subcommands.ExitSuccess
}

func main() {
	subcommands.Register(&initCmd{}, "")
	subcommands.Register(&archiveCmd{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
