// https://rapidapi.com/dpventures/api/wordsapi
package rapidapi

import (
	"fmt"
	"strings"
)

type Response struct {
	Word          string        `json:"word"`
	Syllables     Syllable      `json:"syllables"`
	Frequency     float64       `json:"frequency"`
	Pronunciation Pronunciation `json:"pronunciation"`
	Results       []Result      `json:"results"`
}

type Syllable struct {
	Count int      `json:"count"`
	List  []string `json:"list"`
}

type Pronunciation struct {
	All string `json:"all"`
}

type Result struct {
	Definition   string   `json:"definition"`
	Derivation   []string `json:"derivation,omitempty"`
	PartOfSpeech string   `json:"partOfSpeech"`
	Synonyms     []string `json:"synonyms"`
	SimilarTo    []string `json:"similarTo,omitempty"`
	TypeOf       []string `json:"typeOf,omitempty"`
	Examples     []string `json:"examples"`
}

func (r Response) ToFlashCard(sideSeparator string) string {
	meanings := make([]string, 0, len(r.Results))
	for _, result := range r.Results {
		lines := make([]string, 0)
		lines = append(lines, fmt.Sprintf("[%s]: %s", result.PartOfSpeech, result.Definition))
		if len(result.Examples) > 0 {
			lines = append(lines, fmt.Sprintf("Examples: %s", strings.Join(result.Examples, ", ")))
		}
		if len(result.Synonyms) > 0 {
			lines = append(lines, fmt.Sprintf("Synonyms: %s", strings.Join(result.Synonyms, ", ")))
		}
		if len(result.SimilarTo) > 0 {
			lines = append(lines, fmt.Sprintf("Similar to: %s", strings.Join(result.SimilarTo, ", ")))
		}
		if len(result.Derivation) > 0 {
			lines = append(lines, fmt.Sprintf("Derivation: %s", strings.Join(result.Derivation, ", ")))
		}
		meanings = append(meanings, strings.Join(lines, "\n"))
	}

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Word: %s /%s/\n", r.Word, r.Pronunciation.All))
	builder.WriteString(sideSeparator)
	builder.WriteString(strings.Join(meanings, "\n"+strings.Repeat("-", 50)+"\n"))

	return builder.String()
}
