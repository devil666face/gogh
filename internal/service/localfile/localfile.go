package localfile

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

const (
	_chunk = 20 * 1024 * 1024 // 25MB
)

type LocalFile struct {
	Id       int
	Filname  string
	Size     int64
	Pieces   []string
	Compress bool
	tempDir  string
	tempId   string
	chunk    int
}

func New(
	id int,
	path string,
	compress bool,
) (*LocalFile, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed get info about file %s: %w", path, err)
	}
	f := LocalFile{
		Id:       id,
		Filname:  filepath.Base(path),
		Size:     stat.Size(),
		Compress: compress,
		chunk:    _chunk,
		tempId:   uuid.NewString(),
	}
	f.tempDir = filepath.Join(os.TempDir(), f.tempId)
	if err := os.MkdirAll(f.tempDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp dir %s: %w", f.tempDir, err)
	}
	if err := f.split(path); err != nil {
		return nil, fmt.Errorf("failed to spilt file %s on chunks: %w", path, err)
	}
	return &f, nil
}

func (f *LocalFile) Clear() error {
	if err := os.RemoveAll(f.tempDir); err != nil {
		return fmt.Errorf("failed to remove temp %s: %w", f.tempDir, err)
	}
	return nil
}

func (f *LocalFile) split(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	var (
		buff = make([]byte, f.chunk)
		num  = 1
	)

	for {
		n, err := file.Read(buff)
		if err != nil && err != io.EOF {
			return fmt.Errorf("could not read file: %w", err)
		}
		if n == 0 {
			break
		}

		filename := fmt.Sprintf("%s.%d.gz", filepath.Join(f.tempDir, f.tempId), num)
		if err := f.writePiece(filename, buff[:n]); err != nil {
			return err
		}

		f.Pieces = append(f.Pieces, filename)
		num++
	}
	return nil
}

func (f *LocalFile) writePiece(filename string, data []byte) error {
	piece, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create part file: %w", err)
	}
	defer piece.Close()

	if f.Compress {
		gw := gzip.NewWriter(piece)
		defer gw.Close()

		if _, err = gw.Write(data); err != nil {
			return fmt.Errorf("could not write to gzip part file: %w", err)
		}
		if err := gw.Close(); err != nil {
			return fmt.Errorf("could not close gzip writer: %w", err)
		}
	} else {
		if _, err = piece.Write(data); err != nil {
			return fmt.Errorf("could not write to part file: %w", err)
		}
	}

	return nil
}

// func (f *LocalFile) split(path string) error {
// 	file, err := os.Open(path)
// 	if err != nil {
// 		return fmt.Errorf("could not open file: %v", err)
// 	}
// 	defer file.Close()

// 	var (
// 		buff = make([]byte, f.chunk)
// 		num  = 1
// 	)

// 	for {
// 		n, err := file.Read(buff)
// 		if err != nil && err != io.EOF {
// 			return fmt.Errorf("could not read file: %w", err)
// 		}
// 		if n == 0 {
// 			break
// 		}

// 		filename := fmt.Sprintf("%s.%d.gz", filepath.Join(f.tempDir, f.tempId), num)
// 		piece, err := os.Create(filename)
// 		if err != nil {
// 			return fmt.Errorf("could not create part file: %w", err)
// 		}

// 		switch {
// 		case f.Compress:
// 			gw := gzip.NewWriter(piece)
// 			defer gw.Close()

// 			if _, err = gw.Write(buff[:n]); err != nil {
// 				return fmt.Errorf("could not write to gzip part file: %w", err)
// 			}
// 			if err := gw.Close(); err != nil {
// 				return fmt.Errorf("could not close gzip writer: %w", err)
// 			}
// 		default:
// 			if _, err = piece.Write(buff[:n]); err != nil {
// 				return fmt.Errorf("could not write to part file: %w", err)
// 			}
// 		}
// 		if err := piece.Close(); err != nil {
// 			return fmt.Errorf("could not close part file: %w", err)
// 		}

// 		f.Pieces = append(f.Pieces, filename)
// 		num++
// 	}
// 	return nil
// }

// func (f *LocalFile) split(path string) error {
// 	file, err := os.Open(path)
// 	if err != nil {
// 		return fmt.Errorf("could not open file: %v", err)
// 	}
// 	defer file.Close()

// 	var (
// 		buff = make([]byte, f.chunk)
// 		num  = 1
// 	)

// 	for {
// 		n, err := file.Read(buff)
// 		if err != nil && err != io.EOF {
// 			return fmt.Errorf("could not read file: %w", err)
// 		}
// 		if n == 0 {
// 			break
// 		}

// 		filename := fmt.Sprintf("%s.%d.zip", filepath.Join(f.tempDir, f.tempId), num)
// 		piece, err := os.Create(filename)
// 		if err != nil {
// 			return fmt.Errorf("could not create part file: %w", err)
// 		}

// 		if _, err = piece.Write(buff[:n]); err != nil {
// 			return fmt.Errorf("could not write to part file: %w", err)
// 		}
// 		if err := piece.Close(); err != nil {
// 			return fmt.Errorf("cloud not close path of file: %w", err)
// 		}

// 		f.Pieces = append(f.Pieces, filename)
// 		num++
// 	}
// 	return nil
// }
