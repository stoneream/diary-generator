package initialize

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type InitializeCmd struct {
	Now time.Time
}

func (p *InitializeCmd) Execute() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	currentDirName := filepath.Base(currentDir)

	// ファイル名を生成
	// e.g. diary_2024-01-01.md
	ymdNow := p.Now.Format("2006-01-02")
	targetFileName := fmt.Sprintf("%s_%s.md", currentDirName, ymdNow)

	// ファイルの存在チェック
	targetFilePath := filepath.Join(currentDir, targetFileName)
	if _, err := os.Stat(targetFilePath); err == nil {
		log.Fatalf("file already exists: %s", targetFilePath)
		return err
	}

	// テンプレートファイルが存在するか？
	templateFilePath := filepath.Join(currentDir, "template.md")
	var templateText string

	if _, err := os.Stat(templateFilePath); err != nil {
		log.Printf("template file not found (%s)", templateFilePath)
		templateText = ""
	} else {
		// テンプレートファイル読み込み
		templateTextBytes, err := os.ReadFile(templateFilePath)
		if err != nil {
			log.Fatalf("failed to read template file: %v", err)
			return err
		}
		templateText = string(templateTextBytes)
	}

	text, err := p.templating(currentDirName, templateText)
	if err != nil {
		log.Fatalf("failed to templating: %v", err)
	}

	// ファイルを作成して書き込む
	outputFile, err := os.Create(targetFilePath)
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}

	defer outputFile.Close()

	_, err = outputFile.WriteString(text)
	if err != nil {
		log.Fatalf("failed to write file: %v", err)
	}

	log.Printf("initialized successfully at %s", targetFilePath)

	return nil
}

func (p *InitializeCmd) templating(
	current_dir_name string,
	template_text string,
) (string, error) {
	data := map[string]interface{}{
		"title":   current_dir_name,
		"date":    p.Now.Format("2006-01-02"),
		"content": string(template_text),
	}

	tmpl, err := template.New("").Parse(`---
title: "{{ .title }}"
date: "{{ .date }}"
---

{{ .content }}`)

	if err != nil {
		log.Fatalf("failed to parse template: %v", err)
		return "", err
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, data)

	if err != nil {
		log.Fatalf("failed to execute template: %v", err)
		return "", err
	}

	return buf.String(), nil
}
