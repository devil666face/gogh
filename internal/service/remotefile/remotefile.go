package remotefile

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gogh/internal/models"
	"gogh/internal/service/crypt"
)

type RemoteFile struct {
	Id       int
	Filname  string
	Compress bool
	Pieces   []models.Piece
	tempDir  string
	tempId   string
}

func New(
	id int,
	filename string,
	compress bool,
	pieces []models.Piece,
) (*RemoteFile, error) {
	parsed, err := url.Parse(pieces[0].URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}
	s := strings.Split(path.Base(parsed.Path), ".")
	if len(s) == 0 {
		return nil, fmt.Errorf("failed to get uuid from url")
	}

	rf := RemoteFile{
		Id:       id,
		Filname:  filename,
		Compress: compress,
		Pieces:   pieces,
		tempId:   s[0],
	}

	rf.tempDir = filepath.Join(os.TempDir(), rf.tempId)
	if err := os.MkdirAll(rf.tempDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp dir %s: %w", rf.tempDir, err)
	}
	return &rf, nil
}

func (rf *RemoteFile) Clear() error {
	if err := os.RemoveAll(rf.tempDir); err != nil {
		return fmt.Errorf("failed to remove temp %s: %w", rf.tempDir, err)
	}
	return nil
}

func (rf *RemoteFile) Download() error {
	for _, u := range rf.Pieces {
		if err := rf.downloadPiece(u.URL); err != nil {
			return err
		}
	}
	return rf.join()
}

func ungzipData(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	buf.Write(data)

	gr, err := gzip.NewReader(&buf)
	if err != nil {
		return nil, fmt.Errorf("could not create gzip reader: %w", err)
	}
	defer gr.Close()

	return io.ReadAll(gr)
}

func (rf *RemoteFile) join() error {
	file, err := os.Create(rf.Filname)
	if err != nil {
		return fmt.Errorf("could not create output file: %v", err)
	}
	defer file.Close()

	for num := 1; ; num++ {
		filename := fmt.Sprintf("%s.%d.gz", filepath.Join(rf.tempDir, filepath.Base(rf.tempId)), num)
		piece, err := os.Open(filename)

		if err != nil {
			if os.IsNotExist(err) {
				break
			}
			return fmt.Errorf("could not open part file: %v", err)
		}
		defer piece.Close()

		data, err := io.ReadAll(piece)
		if err != nil {
			return fmt.Errorf("failed read piece: %w", err)
		}

		cryptor, err := crypt.New(rf.Pieces[num-1].Key)
		if err != nil {
			return fmt.Errorf("init decryptor error: %w", err)
		}
		if data, err = cryptor.Decrypt(data); err != nil {
			return fmt.Errorf("decrypt error: %w", err)
		}

		if rf.Compress {
			if data, err = ungzipData(data); err != nil {
				return fmt.Errorf("failed decompress piece: %w", err)
			}
		}
		if _, err := file.Write(data); err != nil {
			return fmt.Errorf("could not copy part file to output file: %w", err)
		}

	}
	return nil
}

func (rf *RemoteFile) downloadPiece(u string) error {
	parsed, err := url.Parse(u)
	if err != nil {
		return fmt.Errorf("failed to parse url: %w", err)
	}
	resp, err := http.Get(u)
	if err != nil {
		fmt.Println("Error:", err)
		return fmt.Errorf("failed to send get request on %s: %w", u, err)
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath.Join(rf.tempDir, path.Base(parsed.Path)))
	if err != nil {
		return fmt.Errorf("failed to create local file %w", err)
	}
	defer out.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to write response to file: %w", err)
	}
	return nil
}
