package models

import "github.com/google/uuid"

type Filestore struct {
	Files map[string]File
}

type File struct {
	Filename string
	Password string
	Pieces   []Piece
}

type Piece struct {
	URL string
}

func NewFilestore() *Filestore {
	return &Filestore{
		Files: make(map[string]File),
	}
}

func (f *Filestore) AddFile(path string) {
	f.Files[uuid.NewString()] = File{
		Filename: path,
	}
}
