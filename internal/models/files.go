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
	Files []File
}

type File struct {
	Filename    string
	Password    string
	Size        int64
	Pieces      []Piece
	Compress    bool
	CreatedDate time.Time
}

type Piece struct {
	URL string
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

func (f *File) SlicePieces() []string {
	var urls []string
	for _, p := range f.Pieces {
		urls = append(urls, p.URL)
	}
	return urls
}

func (f *Filestore) ID() int {
	return len(f.Files)
}

func (f *Filestore) Add(
	id int,
	filename string,
	size int64,
	compess bool,
	url string,
) {
	if id >= len(f.Files) {
		f.Files = append(f.Files,
			File{
				Filename:    filename,
				CreatedDate: time.Now(),
				Size:        size,
				Compress:    compess,
				Pieces: []Piece{
					{URL: url},
				},
			},
		)
	} else {
		file := f.Files[id]
		file.Pieces = append(file.Pieces, Piece{URL: url})
		f.Files[id] = file
	}
	sort.Slice(f.Files, func(i, j int) bool {
		return f.Files[i].CreatedDate.Before(f.Files[j].CreatedDate)
	})
}

func (f *Filestore) Delete(id int) error {
	if id < 0 || id >= len(f.Files) {
		return fmt.Errorf("index out of range %d", len(f.Files))
	}
	f.Files = append(f.Files[:id], f.Files[id+1:]...)
	return nil
}
