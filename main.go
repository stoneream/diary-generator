package main

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/google/subcommands"
)

type initCmd struct {
	baseDirectory string
	templateFile  string
}

func (*initCmd) Name() string     { return "init" }
func (*initCmd) Synopsis() string { return "Initialize a diary" }
func (*initCmd) Usage() string {
	return `init:
	Initialize a diary.
`
}
func (p *initCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.baseDirectory, "base-directory", "", "base directory")
	f.StringVar(&p.templateFile, "template-file", "", "template file path")
}
func (p *initCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if p.baseDirectory == "" || p.templateFile == "" {
		return subcommands.ExitUsageError
	}

	now := time.Now()
	targetDirName := now.Format("2006-01-02")
	templateFileName := filepath.Base(p.templateFile)
	targetFilePath := filepath.Join(p.baseDirectory, targetDirName, templateFileName)

	// テンプレートファイルの存在チェック
	_, err := os.Stat(templateFileName)
	if err != nil {
		return subcommands.ExitFailure
	}

	// 出力先ディレクトリの存在チェック
	_, err = os.Stat(filepath.Dir(targetFilePath))
	if err != nil {
		// ディレクトリが存在しない場合は作成
		err = os.MkdirAll(filepath.Dir(targetFilePath), 0755)
		if err != nil {
			return subcommands.ExitFailure
		}
	}

	// テンプレートファイルをコピー
	input, err := os.ReadFile(p.templateFile)
	if err != nil {
		return subcommands.ExitFailure
	}

	output, err := os.Create(targetFilePath)
	if err != nil {
		return subcommands.ExitFailure
	}
	defer output.Close()

	_, err = output.Write(input)
	if err != nil {
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

func main() {
	subcommands.Register(&initCmd{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
