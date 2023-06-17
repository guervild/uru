package _go

import (
	"embed"
	"fmt"
	"math/rand"
	"strings"

	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/models"

	"golang.org/x/exp/slices"
)

type EnglishWordsEncoder struct {
	Name          string
	Description   string
	Debug         bool
	Dictionary    []string
	DictionaryStr string
	LenShellcode  int
}

func NewEnglishWordsEncoder() models.ObjectModel {
	englishWords := common.GetEnglishWords()

	dictionary := make([]string, 0)

	for len(dictionary) < 256 {
		curr_word := strings.ToLower(strings.TrimSpace(englishWords[rand.Intn(len(englishWords))]))

		if len(curr_word) < 3 {
			continue
		}

		if slices.Contains(dictionary, curr_word) {
			continue
		}

		if strings.Contains(curr_word, "'") {
			continue
		}

		dictionary = append(dictionary, curr_word)
	}

	for i, w := range dictionary {
		dictionary[i] = fmt.Sprintf("\"%s\"", w)
	}

	dictionaryStr := fmt.Sprintf("[]string{%s}", strings.Join(dictionary, ", "))

	return &EnglishWordsEncoder{
		Name:          "english_words",
		Description:   "Use english word to encode given data",
		Debug:         false,
		Dictionary:    dictionary,
		DictionaryStr: dictionaryStr,
	}
}

func (e *EnglishWordsEncoder) Encode(shellcode []byte) ([]byte, error) {
	encoded := make([]string, len(shellcode))

	for i, b := range shellcode {
		encoded[i] = e.Dictionary[int(b)]
	}

	str := strings.Join(encoded, ",")

	decoded := make([]byte, 0)

	str2 := string([]byte(str))

	sliceStr2 := strings.Split(str2, ",")

	for i, w := range sliceStr2 {
		sliceStr2[i] = strings.Trim(w, "\"")
	}

	for _, b := range sliceStr2 {
		for j, db := range e.Dictionary {
			if db == b {
				decoded = append(decoded, byte(j))
				break
			}
		}
	}

	b := []byte(str)
	return b, nil
}

func (e *EnglishWordsEncoder) GetImports() []string {
	return []string{
		`"strings"`,
	}
}

func (e *EnglishWordsEncoder) RenderInstanciationCode(data embed.FS) (string, error) {
	return common.CommonRendering(data, "templates/go/encoders/english_words/instanciation.go.tmpl", e)
}

func (e *EnglishWordsEncoder) RenderFunctionCode(data embed.FS) (string, error) {
	return common.CommonRendering(data, "templates/go/encoders/english_words/functions.go.tmpl", e)
}
