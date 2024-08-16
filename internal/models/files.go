package models

import "time"

type Filestore struct {
	Files map[string]File
}

type File struct {
	Filename    string
	Password    string
	Pieces      []Piece
	CreatedDate time.Time
}

type Piece struct {
	URL string
}

func NewFilestore() *Filestore {
	return &Filestore{
		Files: make(map[string]File),
	}
}

func (f *Filestore) Add(
	id, filename, url string,
) {
	file, ok := f.Files[id]
	if !ok {
		f.Files[id] = File{
			Filename:    filename,
			CreatedDate: time.Now(),
			Pieces: []Piece{
				{URL: url},
			},
		}
		return
	}
	file.Pieces = append(file.Pieces, Piece{URL: url})
	f.Files[id] = file
}
