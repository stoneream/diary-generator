package logic

import (
	"path/filepath"
	"strings"

	"github.com/stoneream/diary-generator/v2/data"
)

type MarkdownWithTOC struct {
	TargetFile data.TargetFile
	Headings   []Heading
}

func GenerateSummaryContent(markdownWithTOCs []MarkdownWithTOC, baseDir string) ([]string, error) {
	var summaryLines []string
	for _, markdownWithTOC := range markdownWithTOCs {
		// 相対パスの計算
		relativePath, err := filepath.Rel(baseDir, markdownWithTOC.TargetFile.Path)
		if err != nil {
			// 相対パス計算に失敗した場合はファイル名のみを使用
			relativePath = markdownWithTOC.TargetFile.Info.Name()
		}

		fileLink := strings.Join([]string{"- [", markdownWithTOC.TargetFile.Info.Name(), "](", relativePath, ")"}, "")
		summaryLines = append(summaryLines, fileLink)
		for _, heading := range markdownWithTOC.Headings {
			indent := strings.Repeat("  ", heading.Level-1)
			headingLine := strings.Join([]string{"  - ", heading.Text}, "")
			summaryLines = append(summaryLines, indent+headingLine)
		}
	}

	return summaryLines, nil
}
