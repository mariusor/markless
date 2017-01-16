package markless

import (
	"fmt"
	"github.com/golang-commonmark/markdown"
)

const (
	exit_success = 0
	exit_error   = 1
)

type app struct {
	buffer *Buffer
	path   string
	follow bool
}

type option func(a *app)

func WithBuffer(file string) option {
	return func(a *app) {
		a.buffer, _ = NewBuffer(file)
	}
}

func Follow(f bool) option {
	return func(a *app) {
		a.follow = f
		// will follow
	}
}

func Init(opts ...option) *app {
	a := &app{}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func (a *app) Run() (int, error) {
	_, err := a.buffer.Read(default_buffer_height, 0)
	checkErr(err)

	md := markdown.New()
	output := md.RenderToString(a.buffer.Data())

	fmt.Printf("%s", output)

	fmt.Printf("Read %d bytes", a.buffer.Size)

	return exit_success, nil
}
