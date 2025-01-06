package model

import "google.golang.org/api/drive/v3"

type ListFilesResult struct {
	FileList   []*drive.File
	StartIndex int
}
