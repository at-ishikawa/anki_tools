// https://dictionaryapi.dev/
package free_dictionary_api

type Definition struct {
	Definition string   `json:"definition"`
	Example    string   `json:"example"`
	Synonyms   []string `json:"synonyms,omitempty"`
	Antonyms   []string `json:"antonyms,omitempty"`
}

type Meaning struct {
	PartOfSpeech string       `json:"partOfSpeech"`
	Definitions  []Definition `json:"definitions"`
	Synonyms     []string     `json:"synonyms,omitempty"`
	Antonyms     []string     `json:"antonyms,omitempty"`
}

type License struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Phonetics struct {
	Text    string  `json:"text"`
	Audio   string  `json:"audio,omitempty"`
	Source  string  `json:"source,omitempty"`
	License License `json:"license,omitempty"`
}

type Word struct {
	Word      string      `json:"word"`
	Phonetic  string      `json:"phonetic"`
	Phonetics []Phonetics `json:"phonetics"`
	Origin    string      `json:"origin"`
	Meanings  []Meaning   `json:"meanings"`

	License    License  `json:"license,omitempty"`
	SourceURLs []string `json:"sourceUrls,omitempty"`
}

type Response []Word
