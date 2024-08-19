package localfile

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gogh/internal/service/crypt"

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

		data := buff[:n]
		if lf.Compress {
			if data, err = gzipData(data); err != nil {
				return fmt.Errorf("gzip error: %w", err)
			}
		}

		cryptor, _ := crypt.New()
		if data, err = cryptor.Encrypt(data); err != nil {
			return fmt.Errorf("encrypt error: %w", err)
		}

		if _, err := piece.Write(data); err != nil {
			return fmt.Errorf("could not write to part file: %w", err)
		}

		lf.Pieces = append(lf.Pieces, LocalPiece{
			Filename: filename,
			Key:      cryptor.B64Key(),
		})
		num++
	}
	return nil
}

func gzipData(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	defer gw.Close()
	if _, err := gw.Write(data); err != nil {
		return nil, fmt.Errorf("could not write to gzip part file: %w", err)
	}
	if err := gw.Close(); err != nil {
		return nil, fmt.Errorf("could not close gzip writer: %w", err)
	}
	return buf.Bytes(), nil
}
