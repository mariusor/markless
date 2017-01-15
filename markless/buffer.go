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
		bytes, err := br.ReadBytes('\n')
		if err != nil {
			break
		}
		l := NewLine(bytes)
		if len(b.Lines) > 0 {
			prev := b.Lines[len(b.Lines)-1]
			l.previous = &prev
		}
		b.Lines = append(b.Lines, *l)
		b.Size += l.Size
	}
	return b.Size, err
}

func (b *Buffer) String() string {
	var y []byte
	for _, line := range b.Lines {
		y = append(y, line.Data...)
	}
	return string(y)
}

type Buffer struct {
	Path  string
	Size  int
	Lines []Line
}
