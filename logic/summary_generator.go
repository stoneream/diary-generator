package logic

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/stoneream/diary-generator/v2/data"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type Heading struct {
	Level int
	Text  string
}

type MarkdownWithTOC struct {
	TargetFile data.TargetFile
	Headings   []Heading
}

type SummaryMetadata struct {
	Title     string
	CreatedAt time.Time
}

func GenerateSummary(targetPrefix, outputDir string) error {
	// 対象ファイル(markdown)の抽出
	markdownFiles, err := GetTargetMarkdownFiles(targetPrefix)
	if err != nil {
		return err
	}

	// 目次の抽出
	var markdownWithTOCs []MarkdownWithTOC
	for _, markdownFile := range markdownFiles {
		headings, err := extractHeadings(markdownFile)
		if err != nil {
			return err
		}
		markdownWithTOC := MarkdownWithTOC{
			TargetFile: markdownFile,
			Headings:   headings,
		}
		markdownWithTOCs = append(markdownWithTOCs, markdownWithTOC)
	}

	// メタデータ
	summaryMetadata := SummaryMetadata{
		Title:     "Summary",
		CreatedAt: time.Now(),
	}

	// サマリ行の生成
	var summaryLines []string
	for _, markdownWithTOC := range markdownWithTOCs {
		fileLink := fmt.Sprintf("- [%s](%s)", markdownWithTOC.TargetFile.Info.Name(), markdownWithTOC.TargetFile.Path)
		summaryLines = append(summaryLines, fileLink)
		for _, heading := range markdownWithTOC.Headings {
			indent := strings.Repeat("  ", heading.Level-1)
			heading := fmt.Sprintf("  - %s", heading.Text)
			summaryLines = append(summaryLines, indent+heading)
		}
	}

	// テンプレート
	summaryText, err := templating(summaryLines, summaryMetadata)
	if err != nil {
		return err
	}

	// ファイルの作成
	summaryFilePath := filepath.Join(outputDir, "summary.md")
	err = os.WriteFile(summaryFilePath, []byte(summaryText), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

func extractHeadings(markdownFile data.TargetFile) ([]Heading, error) {
	// ファイルの読み込み
	content, err := os.ReadFile(markdownFile.Path)
	if err != nil {
		return nil, err
	}

	// YAML Front Matterの削除
	re := regexp.MustCompile(`(?s)^---\n.*?\n---\n`)
	stringContent := re.ReplaceAllString(string(content), "")

	markdown := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)
	reader := text.NewReader([]byte(stringContent))
	document := markdown.Parser().Parse(reader)

	// 見出しの抽出
	var headings []Heading
	err = ast.Walk(document, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if heading, ok := n.(*ast.Heading); ok {
			var headingText string

			ast.Walk(heading, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
				if entering {
					if text, ok := n.(*ast.Text); ok {
						headingText += string(text.Segment.Value(reader.Source()))
					}
				}
				return ast.WalkContinue, nil
			})

			heading := Heading{
				Level: heading.Level,
				Text:  headingText,
			}
			headings = append(headings, heading)
		}

		return ast.WalkContinue, nil
	})

	return headings, err
}

func templating(
	summaryLines []string,
	summaryMetadata SummaryMetadata,
) (string, error) {
	data := map[string]interface{}{
		"title":      summaryMetadata.Title,
		"created_at": summaryMetadata.CreatedAt.Format("2006-01-02 15:04:05"),
		"content":    strings.Join(summaryLines, "\n"),
	}

	tmpl, err := template.New("").Parse(`---
title: "{{ .title }}"
created_at: "{{ .created_at }}"
---

{{ .content }}`)

	if err != nil {
		return "", err
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
