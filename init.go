package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/subcommands"
)

type InitCmd struct {
	baseDirectory string
	templateFile  string
}

func (*InitCmd) Name() string     { return "init" }
func (*InitCmd) Synopsis() string { return "Initialize a diary" }
func (*InitCmd) Usage() string {
	return `init:
	Initialize a diary.
		--base-directory: base directory path
		--template-file: template file path
`
}
func (p *InitCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.baseDirectory, "base-directory", "", "base directory")
	f.StringVar(&p.templateFile, "template-file", "", "template file path")
}

func (p *InitCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if p.baseDirectory == "" || p.templateFile == "" {
		log.Println("Error: --base-directory and --template-file are required")
		return subcommands.ExitUsageError
	}

	now := time.Now()
	ymdNow := now.Format("2006-01-02")
	targetFilePath := filepath.Join(p.baseDirectory, ymdNow, filepath.Base(p.templateFile))

	// テンプレートファイルの存在チェック
	_, err := os.Stat(p.templateFile)
	if err != nil {
		log.Println("Error: template file not found:", p.templateFile)
		return subcommands.ExitFailure
	}

	// 出力先ディレクトリの存在チェック
	_, err = os.Stat(filepath.Dir(targetFilePath))
	if err != nil {
		// ディレクトリが存在しない場合は作成
		err = os.MkdirAll(filepath.Dir(targetFilePath), 0755)
		if err != nil {
			log.Println("Error: failed to create directory:", err)
			return subcommands.ExitFailure
		}
	} else {
		log.Println("Error: directory already exists:", filepath.Dir(targetFilePath))
		return subcommands.ExitFailure
	}

	// テンプレートファイルをコピー
	input, err := os.ReadFile(p.templateFile)
	if err != nil {
		log.Println("Error: failed to read template file:", err)
		return subcommands.ExitFailure
	}

	templating := strings.ReplaceAll(string(input), "%TODAY%", ymdNow)

	output, err := os.Create(targetFilePath)
	if err != nil {
		log.Println("Error: failed to create target file:", err)
		return subcommands.ExitFailure
	}
	defer output.Close()

	_, err = output.Write([]byte(templating))
	if err != nil {
		log.Println("Error: failed to write to target file:", err)
		return subcommands.ExitFailure
	}

	log.Println("Diary initialized successfully at", targetFilePath)
	return subcommands.ExitSuccess
}
