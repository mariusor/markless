package markless

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	max_buffer_height     = 1024
	default_buffer_height = 256
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func NewBuffer(path string) (*Buffer, error) {
	fullPath, _ := filepath.Abs(path)
	_, err := os.Stat(fullPath)
	return &Buffer{path: fullPath}, err
}

func (b *Buffer) Read(height int, offset int) (int, error) {
	f, err := os.Open(b.path)
	checkErr(err)

	br := bufio.NewReader(f)
	lineCount := 0
	for {
		bytes, err := br.ReadBytes('\n')
		if err != nil {
			break
		}
		if lineCount++; lineCount <= offset {
			continue
		}

		l := NewLine(bytes)
		if len(b.lines) > 0 {
			prev := b.lines[len(b.lines)-1]
			l.previous = &prev
		}
		b.lines = append(b.lines, *l)
		b.Size += l.Size
		if b.Height() > height {
			return b.Size, nil
		}
		if b.Height() > max_buffer_height {
			return b.Size, errors.New(fmt.Sprintf("Max buffer height exceeded", max_buffer_height))
		}
	}
	return b.Size, err
}

func (b *Buffer) Data() []byte {
	var y []byte
	for _, line := range b.lines {
		y = append(y, line.Data...)
	}
	return y
}

func (b *Buffer) String() string {
	return string(b.Data())
}

func (b *Buffer) Height() int {
	return len(b.lines)
}

type Buffer struct {
	path   string
	offset int
	lines  []Line
	Size   int
}
