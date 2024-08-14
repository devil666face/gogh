package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

const (
	_chunk = 20 * 1024 * 1024 // 25MB
	// _chunk = 5 * 1024 * 1024 // 25MB
)

type File struct {
	Id      string
	Pieces  []string
	path    string
	tempDir string
	chunk   int
}

func New(_path string) (*File, error) {
	f := File{
		Id:    uuid.NewString(),
		path:  _path,
		chunk: _chunk,
	}
	f.tempDir = filepath.Join(os.TempDir(), f.Id)
	if err := os.MkdirAll(f.tempDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp dir %s: %w", f.tempDir, err)
	}
	if err := f.split(); err != nil {
		return nil, fmt.Errorf("failed to spilt file %s on chunks: %w", _path, err)
	}
	return &f, nil
}

func (f *File) Clear() error {
	if err := os.RemoveAll(f.tempDir); err != nil {
		return fmt.Errorf("failed to remove temp %s: %w", f.tempDir, err)
	}
	return nil
}

func (f *File) split() error {
	file, err := os.Open(f.path)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
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

		filename := fmt.Sprintf("%s.%d.zip", filepath.Join(f.tempDir, f.Id), num)
		piece, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("could not create part file: %w", err)
		}

		if _, err = piece.Write(buff[:n]); err != nil {
			return fmt.Errorf("could not write to part file: %w", err)
		}
		if err := piece.Close(); err != nil {
			return fmt.Errorf("cloud not close path of file: %w", err)
		}

		f.Pieces = append(f.Pieces, filename)
		num++
	}
	return nil
}
