package inpx

import (
	"archive/zip"
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	DELIMITER = ""
)

type Parser interface {
	Parse(inpFileName string) (chan Inp, error)
}

type InpxParser struct {
	Basedir       string
	IndexFilename string
	Delimiter     string
	Err           error
}

type Inp struct {
	Author   string
	Genre    string
	Title    string
	Series   string
	Serno    string
	File     string
	Size     string
	Libid    string
	Del      string
	Ext      string
	Date     string
	Lang     string
	Librate  string
	Keywords string
}

func NewInpx(basedir, indexFilename string) (*InpxParser, error) {
	fi, err := os.Stat(basedir)
	if err != nil {
		return nil, fmt.Errorf("Can't open inpx library dir: %w", err)
	}
	if !fi.IsDir() {
		return nil, fmt.Errorf("Inpx lib path is not a directory")
	}

	inpx := InpxParser{
		Basedir:       basedir,
		IndexFilename: indexFilename,
		Delimiter:     DELIMITER,
	}

	return &inpx, nil
}

func (inpx *InpxParser) Parse(ctx context.Context) <-chan *Inp {
	path := filepath.Join(inpx.Basedir, inpx.IndexFilename)

	ch := make(chan *Inp)
	wg := new(sync.WaitGroup)

	zipReader, err := zip.OpenReader(path)
	if err != nil {
		inpx.Err = fmt.Errorf("Can't open inpx index file: %w", err)
		close(ch)
		return ch
	}

	go func() {
		for _, f := range zipReader.File {
			fd, err := f.Open()
			if err != nil {
				// Ignoring non-inp files
				continue
			}
			wg.Add(1)

			go func(f *io.ReadCloser) {
				defer wg.Done()
				defer fd.Close()

				sc := bufio.NewScanner(fd)
				for sc.Scan() {
					if err := sc.Err(); err != nil {
						inpx.Err = fmt.Errorf("Can't parse inp file")
						return
					}
					inp, err := fillUpInp(strings.Split(sc.Text(), inpx.Delimiter))
					if err != nil {
						continue
					}
					select {
					case <-ctx.Done():
						return
					default:
						ch <- inp
					}
				}
			}(&fd)

		}

		go func() {
			wg.Wait()
			zipReader.Close()
			close(ch)
		}()

	}()

	return ch
}

func fillUpInp(fields []string) (*Inp, error) {
	if len(fields) < 14 {
		return nil, fmt.Errorf("Wrong number of fields: %#v", fields)
	}

	inp := Inp{}

	inp.Author = fields[0]
	inp.Genre = fields[1]
	inp.Title = fields[2]
	inp.Series = fields[3]
	inp.Serno = fields[4]
	inp.File = fields[5]
	inp.Size = fields[6]
	inp.Libid = fields[7]
	inp.Del = fields[8]
	inp.Ext = fields[9]
	inp.Date = fields[10]
	inp.Lang = fields[11]
	inp.Librate = fields[12]
	inp.Keywords = fields[13]

	return &inp, nil
}
