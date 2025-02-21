package archive

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/stoneream/diary-generator/v2/data"
	"github.com/stoneream/diary-generator/v2/logic"
)

type ArchiveCmd struct {
	TargetYM string
}

func (p *ArchiveCmd) Execute() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// 対象ファイル(markdown)の抽出
	currentDirName := filepath.Base(currentDir)
	markdownFiles, err := logic.GetTargetMarkdownFiles(currentDirName)
	if err != nil {
		return err
	}
	// 対象日時で絞り込み
	var targetFiles []data.TargetFile
	for _, markdownFile := range markdownFiles {
		if strings.HasPrefix(markdownFile.Metadata.Date, p.TargetYM) {
			targetFiles = append(targetFiles, markdownFile)
		} else {
			log.Printf("Skip: Not target YearMonth: %s", markdownFile.Path)
		}
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
		archiveFilePath := filepath.Join(archiveDirPath, filepath.Base(targetFile.Path))
		err = os.Rename(targetFile.Path, archiveFilePath)
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
