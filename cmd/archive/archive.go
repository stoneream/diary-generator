package archive

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/stoneream/diary-generator/v2/data"
)

type ArchiveCmd struct {
	TargetYM string
}

type targetFile struct {
	path string
	info os.FileInfo
}

func (p *ArchiveCmd) Execute() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// 対象ファイル(markdown)の抽出
	targetFiles, err := p.getTargetMarkdownFiles()
	if err != nil {
		return err
	}
	if len(targetFiles) == 0 {
		log.Println("No target files.")
		return nil
	}

	// 移動先(アーカイブディレクトリ)の作成
	archiveDirPath := filepath.Join(currentDir, "archive", p.TargetYM)
	if _, err := os.Stat(archiveDirPath); err != nil {
		err = os.MkdirAll(archiveDirPath, 0755)
		if err != nil {
			log.Fatalf("failed to create archive directory: %v", err)
		}
	}

	for _, targetFile := range targetFiles {
		// ファイルの移動
		archiveFilePath := filepath.Join(archiveDirPath, filepath.Base(targetFile.path))
		err = os.Rename(targetFile.path, archiveFilePath)
		if err != nil {
			log.Fatalf("failed to move file: %v", err)
		}
	}

	// フォルダの移動 (assets)
	// assets フォルダの存在チェック
	assetsDirPath := filepath.Join(currentDir, "assets")
	if _, err := os.Stat(assetsDirPath); err == nil {
		// 移動先ディレクトリの作成
		archiveAssetDirPath := filepath.Join(archiveDirPath, "assets")
		if _, err := os.Stat(archiveAssetDirPath); err != nil {
			err = os.MkdirAll(archiveAssetDirPath, 0755)
			if err != nil {
				log.Fatalf("failed to create archive asset directory: %v", err)
			}
		}

		//フォルダの移動
		archiveAssetsDirPath := filepath.Join(archiveDirPath, "assets")
		err = os.Rename(assetsDirPath, archiveAssetsDirPath)
		if err != nil {
			log.Fatalf("failed to move assets directory: %v", err)
		}
	}

	return nil
}

func (p *ArchiveCmd) getTargetMarkdownFiles() ([]targetFile, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current directory: %v", err)
		return nil, err
	}

	currentDirName := filepath.Base(currentDir)
	var targetDirPaths []targetFile

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

		// カレントディレクトリのプレフィックスで始まらない場合はスキップ
		if !strings.HasPrefix(filename, currentDirName) {
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

		// 対象年月のファイルのみを抽出
		if strings.HasPrefix(metadata.Date, p.TargetYM) {
			log.Println("Target:", path)
			targetDirPaths = append(
				targetDirPaths,
				targetFile{
					path: path,
					info: info,
				},
			)
		} else {
			log.Printf("Skip: Not target YearMonth: %s", path)
		}
	}

	return targetDirPaths, err
}
