package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"gitlab.com/golang-commonmark/markdown"
)

var exitWithError = func(e error) {
	fmt.Fprintf(os.Stderr, "error: %s\n", e)
	os.Exit(1)
	return
}

// !(path/to/filename.md)
func ruleFileTransclude(basePath string) func(s *markdown.StateInline, silent bool) bool {
	memoizedTokens := make(map[string][]markdown.Token)
	return func(s *markdown.StateInline, silent bool) bool {
		pos := s.Pos
		max := s.PosMax

		if pos+2 >= max {
			return false
		}

		path := make([]byte, 0)
		src := s.Src
		if src[pos] != '!' || src[pos+1] != '(' {
			return false
		}
		for i := pos + 2; i < max; i++ {
			if src[i] == ')' {
				break
			}
			path = append(path, src[i])
		}
		if len(path) == 0 {
			return false
		}
		tokens, exists := memoizedTokens[string(path)]
		if !exists {
			newFile := filepath.Join(basePath, string(path))
			newDoc, err := os.ReadFile(newFile)
			if err != nil {
				return false
			}
			tokens = s.Md.Block.Parse(newDoc, s.Md, s.Env)
			memoizedTokens[string(path)] = tokens
		}
		for _, tok := range tokens {
			if tt, ok := tok.(*markdown.Inline); ok {
				for _, tt := range s.Md.Inline.Parse(tt.Content, s.Md, s.Env) {
					s.PushToken(tt)
				}
				continue
			}
			s.PushToken(tok)
		}
		s.Pos = s.PosMax

		return true
	}
}

func main() {
	var status int = 0
	flag.Parse()

	for _, fileName := range flag.Args() {
		data, err := os.ReadFile(fileName)
		if err != nil {
			exitWithError(err)
		}
		data = bytes.Trim(data, "\x00")

		markdown.RegisterInlineRule(10, ruleFileTransclude(filepath.Dir(fileName)))
		md := markdown.New(
			markdown.HTML(true),
			markdown.Linkify(true),
			markdown.Typographer(true),
			markdown.Breaks(true),
		)
		if err = md.Render(os.Stdout, data); err != nil {
			exitWithError(err)
		}
	}

	os.Exit(status)
}
