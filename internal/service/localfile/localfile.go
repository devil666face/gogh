package localfile

import (
	"compress/gzip"
	"fmt"
	"gogh/internal/service/crypt"
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
	Compress bool
	Pieces   []LocalPiece
	tempDir  string
	tempId   string
	chunk    int
}

type LocalPiece struct {
	Filename string
	Key      string
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
	lf := LocalFile{
		Id:       id,
		Filname:  filepath.Base(path),
		Size:     stat.Size(),
		Compress: compress,
		chunk:    _chunk,
		tempId:   uuid.NewString(),
	}
	lf.tempDir = filepath.Join(os.TempDir(), lf.tempId)
	if err := os.MkdirAll(lf.tempDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp dir %s: %w", lf.tempDir, err)
	}
	if err := lf.split(path); err != nil {
		return nil, fmt.Errorf("failed to spilt file %s on chunks: %w", path, err)
	}
	return &lf, nil
}

func (lf *LocalFile) Clear() error {
	if err := os.RemoveAll(lf.tempDir); err != nil {
		return fmt.Errorf("failed to remove temp %s: %w", lf.tempDir, err)
	}
	return nil
}

func (lf *LocalFile) split(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	var (
		buff = make([]byte, lf.chunk)
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

		filename := fmt.Sprintf("%s.%d.gz", filepath.Join(lf.tempDir, lf.tempId), num)

		piece, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("could not create part file: %w", err)
		}
		defer piece.Close()

		if err := lf.writePiece(piece, buff[:n]); err != nil {
			return fmt.Errorf("gzip error: %w", err)
		}

		cryptor, _ := crypt.New()
		// if err := lf.cryptPiece(piece, cryptor); err != nil {
		// 	return fmt.Errorf("crypt error: %w", err)
		// }

		lf.Pieces = append(lf.Pieces, LocalPiece{
			Filename: filename,
			Key:      cryptor.B64Key(),
		})
		num++
	}
	return nil
}

func (lf *LocalFile) writePiece(piece *os.File, data []byte) error {
	switch {
	case lf.Compress:
		gw := gzip.NewWriter(piece)
		defer gw.Close()

		if _, err := gw.Write(data); err != nil {
			return fmt.Errorf("could not write to gzip part file: %w", err)
		}
		if err := gw.Close(); err != nil {
			return fmt.Errorf("could not close gzip writer: %w", err)
		}
	default:
		if _, err := piece.Write(data); err != nil {
			return fmt.Errorf("could not write to part file: %w", err)
		}
	}
	return nil
}

// func (lf *LocalFile) cryptPiece(piece *os.File, cryptor *crypt.Sync) error {
// 	body, err := io.ReadAll(piece)
// 	if err != nil {
// 		return fmt.Errorf("failed read body from piece: %w", err)
// 	}
// 	encrypt, err := cryptor.Encrypt(body)
// 	if err != nil {
// 		return fmt.Errorf("encrypt error: %w", err)
// 	}
// 	if _, err := piece.Write(encrypt); err != nil {
// 		return fmt.Errorf("could not write part of encrypt file: %w", err)
// 	}
// 	return nil
// }
