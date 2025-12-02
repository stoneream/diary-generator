package logic

import (
	"regexp"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type Heading struct {
	Level int
	Text  string
}

func ExtractHeadingsFromBytes(content []byte) ([]Heading, error) {
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
	err := ast.Walk(document, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
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
