package rapidapi

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Reader struct {
}

func NewReader() *Reader {
	return &Reader{}
}

func (r *Reader) Read(dir string) (map[string]Response, error) {
	lookedUpWords, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("os.ReadDir > %w", err)
	}

	dictionaries := make(map[string]Response)
	for _, word := range lookedUpWords {
		if word.Name() == ".gitignore" {
			continue
		}

		file := filepath.Join(dir, word.Name())
		f, err := os.Open(file)
		if err != nil {
			return nil, fmt.Errorf("word: %s, os.Open > %w", word, err)
		}
		defer f.Close()

		contents, err := io.ReadAll(f)
		if err != nil {
			return nil, fmt.Errorf("word: %s. io.ReadAll > %w", word, err)
		}

		var res Response
		if err := json.Unmarshal(contents, &res); err != nil {
			return nil, fmt.Errorf("word: %s. json.Unmarshal > %w", word, err)
		}
		dictionaries[res.Word] = res
	}
	return dictionaries, nil
}
