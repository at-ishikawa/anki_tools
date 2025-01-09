package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/at-ishikawa/anki_tools/internal/free_dictionary_api"
	"github.com/at-ishikawa/anki_tools/internal/rapidapi"
	"github.com/go-resty/resty/v2"
	"github.com/kelseyhightower/envconfig"
)

func storeWord(word string, rootDir string, f func() (*resty.Response, error)) error {
	localFilePath := filepath.Join(rootDir, fmt.Sprintf("%s.json", word))
	if _, err := os.Stat(localFilePath); err == nil {
		// dictionary api response is already stored
		return fmt.Errorf("word '%s' is already stored", word)
	}

	res, err := f()
	if err != nil {
		return fmt.Errorf("f > %w", err)
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

func runMain(word string) error {
	var config Env
	envconfig.MustProcess("", &config)

	dir, err := filepath.Abs(filepath.Join(".", "dictionaries"))
	if err != nil {
		return fmt.Errorf("filepath.Abs > %w", err)
	}

	if err := storeWord(word, filepath.Join(dir, "rapidapi"), func() (*resty.Response, error) {
		var response rapidapi.Response
		client := resty.New()
		res, err := client.R().
			EnableTrace().
			SetHeader("x-rapidapi-host", config.RapidAPIHost).
			SetHeader("x-rapidapi-key", config.RapidAPIKey).
			SetResult(&response).
			Get(
				fmt.Sprintf("https://%s/words/%s", config.RapidAPIHost, word),
			)
		if err != nil {
			return nil, fmt.Errorf("client.R.Get > %w", err)
		}

		return res, nil
	}); err != nil {
		return fmt.Errorf("storeWord for RapidAPI > %w", err)
	}

	if err := storeWord(word, filepath.Join(dir, "free_dictionary_api"), func() (*resty.Response, error) {
		var response free_dictionary_api.Response
		client := resty.New()
		res, err := client.R().
			EnableTrace().
			SetResult(&response).
			Get(
				fmt.Sprintf("https://api.dictionaryapi.dev/api/v2/entries/en/%s", word),
			)
		if err != nil {
			return nil, fmt.Errorf("client.R.Get > %w", err)
		}

		return res, nil
	}); err != nil {
		return fmt.Errorf("storeWord for Free Dictionary API > %w", err)
	}

	return nil
}

type Env struct {
	RapidAPIHost string `envconfig:"RAPID_API_HOST"`
	RapidAPIKey  string `envconfig:"RAPID_API_KEY"`
}

func main() {
	command := os.Args[1]
	if command == "generate" {
		word := os.Args[2]
		if err := runMain(word); err != nil {
			slog.Error("failed to run main",
				"word", word,
				"error", err,
			)
		}
		os.Exit(0)
	}

	r := rapidapi.NewReader()
	res, err := r.Read(filepath.Join("dictionaries", "rapidapi"))
	if err != nil {
		slog.Error("failed to read rapidapi", "error", err)
	}
	sideSeparator := strings.Repeat("-", 50) + "\n"
	cardSeparator := strings.Repeat("=", 50) + "\n\n"
	fmt.Printf("Side separator: %s\n", sideSeparator)
	fmt.Printf("Card separator: %s\n", cardSeparator)
	for _, word := range res {
		fmt.Println(word.ToFlashCard(sideSeparator))
		fmt.Println(cardSeparator)
	}
	os.Exit(0)
}
