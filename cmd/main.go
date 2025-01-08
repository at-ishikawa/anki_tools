package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/at-ishikawa/anki_tools/internal/free_dictionary_api"
	"github.com/go-resty/resty/v2"
)

func runMain(word string) error {
	dir, err := filepath.Abs(filepath.Join(".", "dictionaries"))
	if err != nil {
		return fmt.Errorf("filepath.Abs > %w", err)
	}
	localFilePath := filepath.Join(dir, fmt.Sprintf("%s.json", word))
	if _, err := os.Stat(localFilePath); err == nil {
		// dictionary api response is already stored
		return fmt.Errorf("word '%s' is already stored", word)
	}

	var response free_dictionary_api.Response
	client := resty.New()
	res, err := client.R().
		EnableTrace().
		SetResult(&response).
		Get(
			fmt.Sprintf("https://api.dictionaryapi.dev/api/v2/entries/en/%s", word),
		)
	if err != nil {
		return fmt.Errorf("client.R.Get > %w", err)
	}

	file, err := os.Create(localFilePath)
	if err != nil {
		return fmt.Errorf("os.Create > %w", err)
	}
	defer file.Close()

	if _, err := file.Write(res.Body()); err != nil {
		return fmt.Errorf("file.Write > %w", err)
	}

	return nil
}

func main() {
	word := os.Args[1]
	if err := runMain(word); err != nil {
		slog.Error("failed to run main",
			"word", word,
			"error", err,
		)
	}
}
