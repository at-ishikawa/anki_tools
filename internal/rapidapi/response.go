// https://rapidapi.com/dpventures/api/wordsapi
package rapidapi

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
