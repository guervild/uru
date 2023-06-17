package _go

import (
	"embed"
	"fmt"
	"strconv"

	"github.com/guervild/uru/pkg/common"
	"github.com/guervild/uru/pkg/models"
)

type EnglishWordsEvasion struct {
	Name         string
	Description  string
	NumberOfWord string
}

func NewEnglishWordsEvasion() models.ObjectModel {
	randomNumber := common.RandomInt(1, 999)
	return &EnglishWordsEvasion{
		Name: "english-words",
		Description: `Add a random number of english words to the binary.
  Argument(s):
    NumberOfWord: define the number of english words to add to the binary between 1 and 1000.`,
		NumberOfWord: strconv.Itoa(randomNumber),
	}
}

func (e *EnglishWordsEvasion) GetImports() []string {
	return []string{
		`"os"`,
		`"fmt"`,
	}
}

func (e *EnglishWordsEvasion) RenderInstanciationCode(data embed.FS) (string, error) {
	numberOfWord, err := strconv.Atoi(e.NumberOfWord)

	if err != nil {
		return "", err
	}

	if numberOfWord > 1000 {
		numberOfWord = 999
	}

	if numberOfWord < 0 {
		numberOfWord = common.RandomInt(1, 999)
	}

	words := `
temp := os.Stdout
os.Stdout = nil
`
	var currentRandomInt int
	englishWords := common.GetEnglishWords()

	for i := 0; i <= numberOfWord; i++ {
		currentRandomInt = common.RandomInt(1, 999)
		words += fmt.Sprintf("const Word_%d = \"%s\"\nfmt.Println(Word_%d)\n", i, englishWords[currentRandomInt], i)
	}

	words += "os.Stdout = temp"

	return words, nil
}

func (e *EnglishWordsEvasion) RenderFunctionCode(data embed.FS) (string, error) {
	return "", nil
}
