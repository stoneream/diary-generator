package initialize

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type InitializeCmd struct {
	BaseDirectory string
	TemplateFile  string
}

func (p *InitializeCmd) Execute() error {
	now := time.Now()
	ymdNow := now.Format("2006-01-02")
	targetFilePath := filepath.Join(p.BaseDirectory, ymdNow, filepath.Base(p.TemplateFile))

	// テンプレートファイルの存在チェック
	_, err := os.Stat(p.TemplateFile)
	if err != nil {
		log.Println("Error: template file not found:", p.TemplateFile)
		return err
	}

	// 出力先ディレクトリの存在チェック
	_, err = os.Stat(filepath.Dir(targetFilePath))
	if err != nil {
		// ディレクトリが存在しない場合は作成
		err = os.MkdirAll(filepath.Dir(targetFilePath), 0755)
		if err != nil {
			log.Println("Error: failed to create directory:", err)
			return err
		}
	} else {
		log.Println("Error: directory already exists:", filepath.Dir(targetFilePath))
		return err
	}

	// テンプレートファイルをコピー
	input, err := os.ReadFile(p.TemplateFile)
	if err != nil {
		log.Println("Error: failed to read template file:", err)
		return err
	}

	templating := strings.ReplaceAll(string(input), "%TODAY%", ymdNow)

	output, err := os.Create(targetFilePath)
	if err != nil {
		log.Println("Error: failed to create target file:", err)
		return err
	}
	defer output.Close()

	_, err = output.Write([]byte(templating))
	if err != nil {
		log.Println("Error: failed to write to target file:", err)
		return err
	}

	log.Println("Diary initialized successfully at", targetFilePath)
	return nil
}
