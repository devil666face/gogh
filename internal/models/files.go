package models

import (
	"fmt"
	"sort"
	"time"
)

const (
	_        = iota
	KB int64 = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
)

type Filestore struct {
	Files map[string]File
}

type File struct {
	Filename    string
	Password    string
	Size        int64
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

func (f *File) FormatSize() string {
	switch {
	case f.Size >= GB:
		return fmt.Sprintf("%.2f GB", float64(f.Size)/float64(GB))
	case f.Size >= MB:
		return fmt.Sprintf("%.2f MB", float64(f.Size)/float64(MB))
	case f.Size >= KB:
		return fmt.Sprintf("%.2f KB", float64(f.Size)/float64(KB))
	default:
		return fmt.Sprintf("%d B", f.Size)
	}
}

func (f *Filestore) Add(
	id, filename, url string,
	size int64,
) {
	file, ok := f.Files[id]
	if !ok {
		f.Files[id] = File{
			Filename:    filename,
			CreatedDate: time.Now(),
			Size:        size,
			Pieces: []Piece{
				{URL: url},
			},
		}
		return
	}
	file.Pieces = append(file.Pieces, Piece{URL: url})
	f.Files[id] = file
}

func (f *Filestore) FilesSlice() []File {
	var files = []File{}
	for _, file := range f.Files {
		files = append(files, file)
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].CreatedDate.Before(files[j].CreatedDate)
	})
	return files
}
