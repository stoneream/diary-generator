package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/subcommands"
)

type ArchiveCmd struct {
	baseDirectory string
	startsWith    string
}

type targetDir struct {
	path string
	info os.FileInfo
}

func (*ArchiveCmd) Name() string     { return "archive" }
func (*ArchiveCmd) Synopsis() string { return "Archive a diary" }
func (*ArchiveCmd) Usage() string {
	return `archive:
	Archive a diary.
		--base-directory: base directory path
		--starts-with: starts with string
`
}
func (p *ArchiveCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.baseDirectory, "base-directory", "", "base directory")
	f.StringVar(&p.startsWith, "starts-with", "", "starts with string")
}

func (p *ArchiveCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if p.baseDirectory == "" || p.startsWith == "" {
		log.Println("Error: --base-directory and --starts-with are required")
		return subcommands.ExitUsageError
	}

	_, err := os.Stat(p.baseDirectory)
	if err != nil {
		log.Println("Error: base directory not found:", p.baseDirectory)
		return subcommands.ExitFailure
	}

	err = complateArchiveDir(p.baseDirectory)
	if err != nil {
		log.Println("Error: failed to create archive directory:", err)
		return subcommands.ExitFailure
	}

	targetDirs, err := getTargetDirPaths(p.baseDirectory, p.startsWith)
	if err != nil {
		log.Println("Error: failed to get target directory paths:", err)
		return subcommands.ExitFailure
	}

	for _, targetDir := range targetDirs {
		archiveDirPath := filepath.Join(p.baseDirectory, "archive", targetDir.info.Name())
		_, err := os.Stat(archiveDirPath)
		if err == nil {
			log.Println("Skip: already exists:", archiveDirPath)
			continue
		}

		err = os.Rename(targetDir.path, archiveDirPath)
		if err != nil {
			log.Println("Error: failed to move directory:", targetDir.path, archiveDirPath)
			return subcommands.ExitFailure
		}

		log.Println("Success: move directory:", targetDir.path, archiveDirPath)
	}

	return subcommands.ExitSuccess
}

// `--base-directory` 以下に存在する、`--starts-with` で指定された文字列で始まるディレクトリを取得する
func getTargetDirPaths(baseDirectory, startsWith string) ([]targetDir, error) {
	var targetDirPaths []targetDir
	err := filepath.Walk(baseDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if strings.HasPrefix(info.Name(), startsWith) {
				targetDirPaths = append(targetDirPaths, targetDir{path: path, info: info})
			}
		}

		return nil
	})

	return targetDirPaths, err
}

// `--base-directory` 以下に `archive` ディレクトリが存在しない場合は作成する
func complateArchiveDir(baseDirectory string) error {
	archiveDirPath := filepath.Join(baseDirectory, "archive")
	_, err := os.Stat(archiveDirPath)
	if err != nil {
		err = os.Mkdir(archiveDirPath, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}
