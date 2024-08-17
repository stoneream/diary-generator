package initialize

import (
	"diary-generator/data"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

type InitializeCmd struct {
	BaseDirectory string
	TemplateFile  string
	Name          string
}

func (p *InitializeCmd) Execute() error {
	now := time.Now()
	ymdNow := now.Format("2006-01-02")
	targetFileName := fmt.Sprintf("%s_%s.md", p.Name, ymdNow) // e.g. diary_2024-01-01.md
	targetFilePath := filepath.Join(p.BaseDirectory, targetFileName)

	// アセットディレクトリの存在チェック & 作成
	assetDir := filepath.Join(p.BaseDirectory, "assets")
	if _, err := os.Stat(assetDir); os.IsNotExist(err) {
		err := os.Mkdir(assetDir, 0755)
		if err != nil {
			log.Println("Error: failed to create asset directory:", err)
			return err
		}
	}

	// テンプレートファイルの存在チェック
	_, err := os.Stat(p.TemplateFile)
	if err != nil {
		log.Println("Error: template file not found:", p.TemplateFile)
		return err
	}

	// 生成しようとしているファイルが既に存在するかチェック
	if _, err := os.Stat(targetFilePath); err == nil {
		log.Println("Error: target file already exists:", targetFilePath)
		return err
	}

	outputFile, err := os.Create(targetFilePath)
	if err != nil {
		log.Println("Error: failed to create target file:", err)
		return err
	}
	defer outputFile.Close()

	// メタデータ
	metadata := data.Metadata{
		Title: p.Name,
		Date:  ymdNow,
	}

	metadataYaml, err := metadata.String()
	if err != nil {
		log.Println("Error: failed to convert metadata to yaml:", err)
		return err
	}

	// テンプレートファイル読み込み
	input, err := os.ReadFile(p.TemplateFile)
	if err != nil {
		log.Println("Error: failed to read template file:", err)
		return err
	}

	// 内容の埋込
	data := map[string]interface{}{
		"metadata": metadataYaml,
		"content":  string(input),
	}

	tmpl, err := template.New("").Parse(`---
{{ .metadata }}
---
{{ .content }}
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

	log.Println("Diary initialized successfully at", targetFilePath)
	return nil
}
