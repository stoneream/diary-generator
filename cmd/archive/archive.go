package archive

import (
	"diary-generator/data"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
)

type ArchiveCmd struct {
	BaseDirectory string
	Name          string
	TargetYM      string
	TemplateFile  string
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

	err = p.complateArchiveDir()
	if err != nil {
		log.Println("Error: failed to create archive directory:", err)
		return err
	}

	targetFiles, err := p.getTargetFiles()
	if err != nil {
		log.Println("Error: failed to get target directory paths:", err)
		return err
	}

	for _, targetFile := range targetFiles {
		moveTo := filepath.Join(p.BaseDirectory, "archive", p.TargetYM, targetFile.info.Name())

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
func (p *ArchiveCmd) complateArchiveDir() error {
	archiveDirPath := filepath.Join(p.BaseDirectory, "archive", p.TargetYM)
	_, err := os.Stat(archiveDirPath)
	if err != nil {
		err = os.MkdirAll(archiveDirPath, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}
