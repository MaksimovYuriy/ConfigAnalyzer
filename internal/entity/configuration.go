package entity

import "io/fs"

type Configuration struct {
	Content      string
	FileMetadata *FileMetadata
}

type FileMetadata struct {
	Path string
	Mode fs.FileMode
}
