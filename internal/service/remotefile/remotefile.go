package remotefile

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type RemoteFile struct {
	Id       int
	Filname  string
	Compress bool
	Pieces   []string
	tempDir  string
	tempId   string
}

func New(
	id int,
	filename string,
	compress bool,
	pieces []string,
) (*RemoteFile, error) {
	parsed, err := url.Parse(pieces[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}
	s := strings.Split(path.Base(parsed.Path), ".")
	if len(s) == 0 {
		return nil, fmt.Errorf("failed to get uuid from url")
	}

	f := RemoteFile{
		Id:       id,
		Filname:  filename,
		Compress: compress,
		Pieces:   pieces,
		tempId:   s[0],
	}

	f.tempDir = filepath.Join(os.TempDir(), f.tempId)
	if err := os.MkdirAll(f.tempDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp dir %s: %w", f.tempDir, err)
	}
	return &f, nil
}

func (f *RemoteFile) Clear() error {
	if err := os.RemoveAll(f.tempDir); err != nil {
		return fmt.Errorf("failed to remove temp %s: %w", f.tempDir, err)
	}
	return nil
}

func (f *RemoteFile) download(u string) error {
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

	out, err := os.Create(filepath.Join(f.tempDir, path.Base(parsed.Path)))
	if err != nil {
		return fmt.Errorf("failed to create local file %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write response to file: %w", err)
	}
	return nil
}

func (f *RemoteFile) Download() error {
	for _, u := range f.Pieces {
		if err := f.download(u); err != nil {
			return err
		}
	}
	return f.join()
}

// func (f *RemoteFile) join() error {
// 	file, err := os.Create(f.Filname)
// 	if err != nil {
// 		return fmt.Errorf("could not create output file: %w", err)
// 	}
// 	defer file.Close()

// 	for num := 1; ; num++ {
// 		filename := fmt.Sprintf("%s.%d.gz", filepath.Join(f.tempDir, filepath.Base(f.tempId)), num)
// 		if err := f.processPiece(file, filename); err != nil {
// 			if os.IsNotExist(err) {
// 				break
// 			}
// 			return err
// 		}
// 	}

// 	return nil
// }

func (f *RemoteFile) processPiece(output *os.File, piece *os.File, filename string) error {
	defer piece.Close()

	switch {
	case f.Compress:
		gr, err := gzip.NewReader(piece)
		if err != nil {
			return fmt.Errorf("could not create gzip reader: %w", err)
		}
		defer gr.Close()

		if _, err = io.Copy(output, gr); err != nil {
			return fmt.Errorf("could not copy decompressed data to output file: %w", err)
		}
	default:
		if _, err := io.Copy(output, piece); err != nil {
			return fmt.Errorf("could not copy part file to output file: %w", err)
		}
	}

	return nil
}

func (f *RemoteFile) join() error {
	file, err := os.Create(f.Filname)
	if err != nil {
		return fmt.Errorf("could not create output file: %v", err)
	}
	defer file.Close()

	// var (
	// 	num = 1
	// )

	for num := 1; ; num++ {
		filename := fmt.Sprintf("%s.%d.gz", filepath.Join(f.tempDir, filepath.Base(f.tempId)), num)
		piece, err := os.Open(filename)
		if err != nil {
			if os.IsNotExist(err) {
				break
			}
			return fmt.Errorf("could not open part file: %v", err)
		}

		// switch {
		// case f.Compress:
		// 	gr, err := gzip.NewReader(piece)
		// 	if err != nil {
		// 		piece.Close()
		// 		return fmt.Errorf("could not create gzip reader: %v", err)
		// 	}

		// 	_, err = io.Copy(file, gr)
		// 	if err != nil {
		// 		gr.Close()
		// 		piece.Close()
		// 		return fmt.Errorf("could not copy decompressed data to output file: %v", err)
		// 	}
		// 	gr.Close()
		// 	piece.Close()

		// default:
		// 	_, err = io.Copy(file, piece)
		// 	if err != nil {
		// 		return fmt.Errorf("could not copy part file to output file: %v", err)
		// 	}
		// 	piece.Close()
		// }
		if err := f.processPiece(file, piece, filename); err != nil {
			return err
		}
		// num++
	}

	return nil
}

// func (f *RemoteFile) join() error {
// 	file, err := os.Create(f.Filname)
// 	if err != nil {
// 		return fmt.Errorf("could not create output file: %v", err)
// 	}
// 	defer file.Close()

// 	var (
// 		num = 1
// 	)

// 	for {
// 		filename := fmt.Sprintf("%s.%d.zip", filepath.Join(f.tempDir, filepath.Base(f.tempId)), num)
// 		piece, err := os.Open(filename)
// 		if err != nil {
// 			if os.IsNotExist(err) {
// 				break
// 			}
// 			return fmt.Errorf("could not open part file: %v", err)
// 		}

// 		_, err = io.Copy(file, piece)
// 		if err != nil {
// 			return fmt.Errorf("could not copy part file to output file: %v", err)
// 		}

// 		piece.Close()
// 		num++
// 	}

// 	return nil
// }
