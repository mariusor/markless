package markless

import (
	"bufio"
	"os"
	"path/filepath"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func NewBuffer(path string) (*Buffer, error) {
	fullpath, _ := filepath.Abs(path)
	_, err := os.Stat(fullpath)
	return &Buffer{Path: fullpath}, err
}

func (b *Buffer) Read() (int, error) {
	f, err := os.Open(b.Path)
	checkErr(err)

	br := bufio.NewReader(f)
	for {
		line, err := br.ReadBytes('\n')
		if err != nil {
			break
		}
		b.Data = append(b.Data, line...)
		b.Size += len(line)
	}
	return b.Size, err
}

type Buffer struct {
	Path string
	Size int
	Data []byte
}
