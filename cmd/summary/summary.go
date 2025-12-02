package summary

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/stoneream/diary-generator/v2/logic"
)

type SummaryCmd struct {
	TargetPrefixOpt string
}

func (p *SummaryCmd) Execute() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// サマリする対象ファイルのプレフィックスが指定されている場合はその値を使う
	// 指定されていない場合はカレントディレクトリ名を使う
	// これはarchiveディレクトリなどで実行することを考慮している
	var targetPrefix string
	if p.TargetPrefixOpt == "" {
		targetPrefix = filepath.Base(currentDir)
	} else {
		targetPrefix = p.TargetPrefixOpt
	}

	// 対象ファイル(markdown)の抽出
	markdownFiles, err := logic.GetTargetMarkdownFiles(targetPrefix)
	if err != nil {
		return err
	}

	// 目次の抽出
	var markdownWithTOCs []logic.MarkdownWithTOC
	for _, markdownFile := range markdownFiles {
		content, err := os.ReadFile(markdownFile.Path)
		if err != nil {
			return err
		}

		headings, err := logic.ExtractHeadingsFromBytes(content)
		if err != nil {
			return err
		}

		markdownWithTOCs = append(markdownWithTOCs, logic.MarkdownWithTOC{
			TargetFile: markdownFile,
			Headings:   headings,
		})
	}

	// サマリ行の生成
	summaryLines, err := logic.GenerateSummaryContent(markdownWithTOCs, currentDir)
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

	// ファイルの作成
	// 上書きされるので注意
	summaryFilePath := filepath.Join(currentDir, "summary.md")
	err = os.WriteFile(summaryFilePath, []byte(summaryText), 0644)
	if err != nil {
		log.Fatalf("failed to write file: %v", err)
	}

	return nil
}
