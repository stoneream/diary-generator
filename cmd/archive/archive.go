package archive

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	assetsDirPath := filepath.Join(currentDir, "assets")
	if _, err := os.Stat(assetsDirPath); err == nil {
		archiveAssetsDirPath := filepath.Join(archiveDirPath, "assets")
		err = os.Rename(assetsDirPath, archiveAssetsDirPath)
		if err != nil {
			log.Fatalf("failed to move assets directory: %v", err)
		}
	}

	// サマリーファイルの作成
	err = createSummaryFile(targetFiles, archiveDirPath)
	if err != nil {
		log.Fatalf("failed to create summary file: %v", err)
	}

	return nil
}

func createSummaryFile(targetFiles []data.TargetFile, archiveDirPath string) error {
	// 目次の抽出
	var markdownWithTOCs []logic.MarkdownWithTOC
	for _, targetFile := range targetFiles {
		// アーカイブディレクトリ内の移動後のパスを計算
		archiveFilePath := filepath.Join(archiveDirPath, filepath.Base(targetFile.Path))

		content, err := os.ReadFile(archiveFilePath)
		if err != nil {
			return err
		}

		headings, err := logic.ExtractHeadingsFromBytes(content)
		if err != nil {
			return err
		}

		markdownWithTOCs = append(markdownWithTOCs, logic.MarkdownWithTOC{
			TargetFile: data.TargetFile{
				Path:     archiveFilePath,
				Info:     targetFile.Info,
				Metadata: targetFile.Metadata,
			},
			Headings: headings,
		})
	}

	// サマリ行の生成
	summaryLines, err := logic.GenerateSummaryContent(markdownWithTOCs, archiveDirPath)
	if err != nil {
		return err
	}

	// メタデータ
	summaryMetadata := logic.SummaryMetadata{
		Title:     "Summary",
		CreatedAt: time.Now(),
	}

	// テンプレート
	summaryText, err := logic.TemplatingSummary(summaryLines, summaryMetadata)
	if err != nil {
		return err
	}

	// サマリーファイルの作成
	summaryFilePath := filepath.Join(archiveDirPath, "summary.md")
	err = os.WriteFile(summaryFilePath, []byte(summaryText), 0644)
	if err != nil {
		return err
	}

	log.Printf("Summary file created: %s", summaryFilePath)
	return nil
}
