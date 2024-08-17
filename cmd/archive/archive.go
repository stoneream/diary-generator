package archive

import (
	"diary-generator/data"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type ArchiveCmd struct {
	BaseDirectory         string
	Name                  string
	TargetYM              string
	TemplateFile          string
	EnabledArchiveSummary bool
}

type targetFile struct {
	path string
	info os.FileInfo
}

func (p *ArchiveCmd) Execute() error {
	_, err := os.Stat(p.BaseDirectory)
	if err != nil {
		log.Println("Error: base directory not found:", p.BaseDirectory)
		return err
	}

	archiveDirPath, err := p.complateArchiveDir()
	if err != nil {
		log.Println("Error: failed to create archive directory:", err)
		return err
	}

	targetFiles, err := p.getTargetFiles()
	if err != nil {
		log.Println("Error: failed to get target directory paths:", err)
		return err
	}

	var archivedFilePaths []string
	for _, targetFile := range targetFiles {
		// ファイルの移動
		moveTo := filepath.Join(archiveDirPath, targetFile.info.Name())

		_, err := os.Stat(moveTo)
		if err == nil {
			log.Println("Skip: already exists:", moveTo)
			continue
		}

		err = os.Rename(targetFile.path, moveTo)
		if err != nil {
			log.Println("Error: failed to move directory:", targetFile.path, moveTo)
			return err
		}

		log.Println("Success: move directory:", targetFile.path, moveTo)

		archivedFilePaths = append(archivedFilePaths, moveTo)
	}

	// サマリの生成
	if p.EnabledArchiveSummary {
		// 各ファイルの目次を生成
		var tocs []string
		for _, archivedFilePath := range archivedFilePaths {
			file, err := os.ReadFile(archivedFilePath)
			if err != nil {
				log.Println("Error: failed to read file:", archivedFilePath)
				return err
			}
			filename := filepath.Base(archivedFilePath)

			content := string(file)
			toc, err := generateTOC(content)
			if err != nil {
				log.Println("Error: failed to generate TOC:", archivedFilePath)
				return err
			}

			relativePath, err := filepath.Rel(archiveDirPath, archivedFilePath)
			if err != nil {
				log.Println("Error: failed to get relative path:", archivedFilePath)
				return err
			}

			tocs = append(tocs, fmt.Sprintf("- [%s](%s)\n%s", filename, relativePath, addIntent(toc, 2)))
		}

		// サマリファイルの作成
		summaryPath := filepath.Join(p.BaseDirectory, "archive", p.TargetYM, "summary.md")

		if _, err := os.Stat(summaryPath); err == nil {
			log.Println("Skip: already exists:", summaryPath)
			return nil
		}

		outputFile, err := os.Create(summaryPath)
		if err != nil {
			log.Println("Error: failed to create summary file:", err)
			return err
		}
		defer outputFile.Close()

		// templating
		metadata := data.Metadata{
			Title: fmt.Sprintf("Archived Summary (%s)", p.Name),
			Date:  time.Now().Format("2006-01-02 15:04:05"),
		}
		metadataYaml, err := metadata.String()
		if err != nil {
			log.Println("Error: failed to convert metadata to yaml:", err)
			return err
		}
		data := map[string]interface{}{
			"metadata": metadataYaml,
			"toc":      strings.Join(tocs, "\n"),
		}
		tmpl, err := template.New("").Parse(`---
{{ .metadata }}
---
{{ .toc }}
`)
		if err != nil {
			log.Println("Error: failed to parse template:", err)
			return err
		}

		err = tmpl.Execute(outputFile, data)

		if err != nil {
			log.Println("Error: failed to execute template:", err)
			return err
		}
	}

	return nil
}

func (p *ArchiveCmd) getTargetFiles() ([]targetFile, error) {
	var targetDirPaths []targetFile

	err := filepath.Walk(p.BaseDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		log.Println("Check: file path:", path)

		// ディレクトリはスキップ
		if info.IsDir() && path != p.BaseDirectory {
			log.Println("Skip: directory:", path)
			return filepath.SkipDir
		}

		// テンプレートファイル, markdown以外のファイルはスキップ
		if path == p.TemplateFile || filepath.Ext(path) != ".md" {
			log.Println("Skip: not markdown file:", path)
			return nil
		}

		// ファイルの読み込み
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// メタデータの取得
		metadata := data.Metadata{}
		_, err = frontmatter.Parse(file, &metadata)

		// メタデータが取得できない場合はスキップ
		if err != nil {
			log.Println("Skip: failed to get metadata:", path)
			return nil
		}

		if metadata.Title == p.Name && strings.HasPrefix(metadata.Date, p.TargetYM) {
			log.Println("Target:", path)
			targetDirPaths = append(targetDirPaths, targetFile{path: path, info: info})
		} else {
			log.Println("Skip: not target file:", path)
		}

		return nil
	})

	return targetDirPaths, err
}

// ベースディレクトリ以下に `archive` ディレクトリが存在しない場合は作成する
func (p *ArchiveCmd) complateArchiveDir() (string, error) {
	archiveDirPath := filepath.Join(p.BaseDirectory, "archive", p.TargetYM)
	_, err := os.Stat(archiveDirPath)
	if err != nil {
		err = os.MkdirAll(archiveDirPath, 0755)
		if err != nil {
			return "", err
		}
	}

	return archiveDirPath, nil
}

func generateTOC(content string) (string, error) {
	md := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	reader := text.NewReader([]byte(content))
	doc := md.Parser().Parse(reader)

	var toc strings.Builder

	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
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

			toc.WriteString(fmt.Sprintf("%s- %s\n", strings.Repeat("  ", heading.Level-1), headingText))
		}

		return ast.WalkContinue, nil
	})

	if err != nil {
		return "", err
	}

	return toc.String(), nil
}

func addIntent(text string, spaces int) string {
	indent := strings.Repeat(" ", spaces)
	lines := strings.Split(text, "\n")

	for i, line := range lines {
		lines[i] = indent + line
	}

	return strings.Join(lines, "\n")
}
