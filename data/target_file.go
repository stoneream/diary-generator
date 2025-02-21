package data

import "os"

type TargetFile struct {
	Path     string
	Info     os.FileInfo
	Metadata Metadata
}
