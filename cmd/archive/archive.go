package archive

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type ArchiveCmd struct {
	BaseDirectory string
	StartsWith    string
}

type targetDir struct {
	path string
	info os.FileInfo
}

func (p *ArchiveCmd) Execute() error {
	_, err := os.Stat(p.BaseDirectory)
	if err != nil {
		log.Println("Error: base directory not found:", p.BaseDirectory)
		return err
	}

	err = complateArchiveDir(p.BaseDirectory)
	if err != nil {
		log.Println("Error: failed to create archive directory:", err)
		return err
	}

	targetDirs, err := getTargetDirPaths(p.BaseDirectory, p.StartsWith)
	if err != nil {
		log.Println("Error: failed to get target directory paths:", err)
		return err
	}

	for _, targetDir := range targetDirs {
		archiveDirPath := filepath.Join(p.BaseDirectory, "archive", targetDir.info.Name())
		_, err := os.Stat(archiveDirPath)
		if err == nil {
			log.Println("Skip: already exists:", archiveDirPath)
			continue
		}

		err = os.Rename(targetDir.path, archiveDirPath)
		if err != nil {
			log.Println("Error: failed to move directory:", targetDir.path, archiveDirPath)
			return err
		}

		log.Println("Success: move directory:", targetDir.path, archiveDirPath)
	}

	return nil
}

// `--base-directory` 以下に存在する、`--starts-with` で指定された文字列で始まるディレクトリを取得する
func getTargetDirPaths(BaseDirectory, StartsWith string) ([]targetDir, error) {
	var targetDirPaths []targetDir
	err := filepath.Walk(BaseDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if strings.HasPrefix(info.Name(), StartsWith) {
				targetDirPaths = append(targetDirPaths, targetDir{path: path, info: info})
			}
		}

		return nil
	})

	return targetDirPaths, err
}

// `--base-directory` 以下に `archive` ディレクトリが存在しない場合は作成する
func complateArchiveDir(BaseDirectory string) error {
	archiveDirPath := filepath.Join(BaseDirectory, "archive")
	_, err := os.Stat(archiveDirPath)
	if err != nil {
		err = os.Mkdir(archiveDirPath, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}
