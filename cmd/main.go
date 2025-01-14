package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/at-ishikawa/anki_tools/internal/free_dictionary_api"
	"github.com/at-ishikawa/anki_tools/internal/rapidapi"
	"github.com/go-resty/resty/v2"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
	if res.StatusCode() != http.StatusOK {
		return fmt.Errorf("f > status code: %d, body: %s", res.StatusCode(), string(res.Body()))
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

func runMain(api API, word string) error {
	var config Env
	envconfig.MustProcess("", &config)

	dir, err := filepath.Abs(filepath.Join(".", "dictionaries"))
	if err != nil {
		return fmt.Errorf("filepath.Abs > %w", err)
	}

	if api == APIFreeDictionaryAPI {
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
	} else if api == APIWordsAPIInRapidAPI {
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
	}

	return nil
}

type Env struct {
	RapidAPIHost string `envconfig:"RAPID_API_HOST"`
	RapidAPIKey  string `envconfig:"RAPID_API_KEY"`
}

type API string

func (a *API) Set(val string) error {
	for _, api := range allAPIs {
		if val == string(api) {
			*a = api
			return nil
		}
	}
	return fmt.Errorf("invalid API: %s", val)
}

func (a API) String() string {
	return string(a)
}

func (a *API) Type() string {
	return "API"
}

const (
	APIFreeDictionaryAPI  API = "free_dictionary"
	APIWordsAPIInRapidAPI API = "words_api"
)

var (
	_       pflag.Value = (*API)(nil)
	allAPIs             = []API{APIFreeDictionaryAPI, APIWordsAPIInRapidAPI}
)

func main() {
	rootCommand := cobra.Command{}
	flags := rootCommand.PersistentFlags()

	api := APIFreeDictionaryAPI
	flags.Var(&api, "api", fmt.Sprintf("API to use. Possible values are %v", allAPIs))

	rootCommand.AddCommand(&cobra.Command{
		Use:  "generate",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			word := args[0]
			if err := runMain(api, word); err != nil {
				slog.Error("failed to run main",
					"word", word,
					"error", err,
				)
			}
			return nil
		},
	})
	rootCommand.AddCommand(&cobra.Command{
		Use: "dump",
		RunE: func(cmd *cobra.Command, args []string) error {
			if api != APIWordsAPIInRapidAPI {
				return fmt.Errorf("dump is only available for API: %s", APIWordsAPIInRapidAPI)
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

				// add a youglish link
				url := fmt.Sprintf("https://youglish.com/pronounce/%s", word.Word)
				fmt.Printf("%s\nYouglish: %s\n", strings.Repeat("-", 50), url)
				fmt.Println(cardSeparator)
			}
			return nil
		},
	})

	if err := rootCommand.Execute(); err != nil {
		slog.Error("failed to execute root command", "error", err)
		os.Exit(1)
	}
	os.Exit(0)
}
