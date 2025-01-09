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
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Word: %s\n", r.Word))
	builder.WriteString(sideSeparator)
	builder.WriteString(fmt.Sprintf("Pronunciation: %s\n", r.Pronunciation.All))

	for _, result := range r.Results {
		builder.WriteString(fmt.Sprintf("Definition: %s\n", result.Definition))
		builder.WriteString(fmt.Sprintf("Part of Speech: %s\n", result.PartOfSpeech))
		if len(result.Examples) > 0 {
			builder.WriteString(fmt.Sprintf("Examples: %s\n", strings.Join(result.Examples, ", ")))
		}
		if len(result.Synonyms) > 0 {
			builder.WriteString(fmt.Sprintf("Synonyms: %s\n", strings.Join(result.Synonyms, ", ")))
		}
		if len(result.SimilarTo) > 0 {
			builder.WriteString(fmt.Sprintf("Similar to: %s\n", strings.Join(result.SimilarTo, ", ")))
		}
		if len(result.Derivation) > 0 {
			builder.WriteString(fmt.Sprintf("Derivation: %s\n", strings.Join(result.Derivation, ", ")))
		}
		builder.WriteString(strings.Repeat("-", 50) + "\n")
	}

	return builder.String()
}
