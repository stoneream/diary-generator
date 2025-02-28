package logic

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/stoneream/diary-generator/v2/data"
)

func GetTargetMarkdownFiles(
	targetPrefix string,
) ([]data.TargetFile, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current directory: %v", err)
		return nil, err
	}

	var targetDirPaths []data.TargetFile

	entries, err := os.ReadDir(currentDir)
	if err != nil {
		log.Fatalf("failed to read directory: %v", err)
		return nil, err
	}

	for _, entry := range entries {
		path := filepath.Join(currentDir, entry.Name())
		info, err := entry.Info()
		if err != nil {
			log.Fatalf("failed to get file info: %v", err)
			return nil, err
		}

		log.Printf("Check: file info: %s", path)

		// ディレクトリはスキップ
		if info.IsDir() && path != currentDir {
			log.Printf("Skip: directory: %s", path)
			continue
		}

		filename := filepath.Base(path)

		// template.md はスキップする
		if filename == "template.md" {
			log.Printf("Skip: Template file: %s", path)
			continue
		}

		// Markdownファイル以外はスキップ
		if filepath.Ext(path) != ".md" {
			log.Printf("Skip: Not markdown file: %s", path)
			continue
		}

		// 指定のプレフィックスで始まらない場合はスキップ
		if !strings.HasPrefix(filename, targetPrefix) {
			log.Printf("Skip: Not start with current directory name: %s", path)
			continue
		}

		// ファイルの読み込み
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		// メタデータの取得
		metadata := data.Metadata{}
		_, err = frontmatter.Parse(file, &metadata)

		// メタデータが取得できない場合はスキップ
		if err != nil {
			log.Printf("Skip: Failed to get metadata: %v", err)
			continue
		}

		// 対象ファイルとして追加
		targetDirPaths = append(
			targetDirPaths,
			data.TargetFile{
				Path:     path,
				Info:     info,
				Metadata: metadata,
			},
		)
	}

	return targetDirPaths, err
}
