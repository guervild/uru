package common

import (
	"bytes"
	"embed"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
	"time"

	"github.com/guervild/uru/data"
)

func RandomInt(start, end int) int {
	rand.Seed(time.Now().UnixNano())
	min := start
	max := end
	return rand.Intn(max-min+1) + min
}

func CommonRendering(data embed.FS, pathtorender string, i interface{}) (string, error) {

	t, err := template.ParseFS(data, pathtorender)
	if err != nil {
		return "", err
	}

	var tplBuffer bytes.Buffer

	err = t.Execute(&tplBuffer, i)
	if err != nil {
		return "", err
	}

	return tplBuffer.String(), nil
}

func GetGolangByteArray(data []byte, lang string) string {

	var newData []string

	if lang == "go" {
		for _, v := range data {
			newData = append(newData, fmt.Sprintf("%d", v))
		}
		return fmt.Sprintf("[]byte { %s }", strings.Join(newData, ","))
	} else if lang == "c" {
		for _, v := range data {
			newData = append(newData, fmt.Sprintf("\\x%x", v))
		}
		return fmt.Sprintf("\"%s\";", strings.Join(newData, ""))
	}

	return fmt.Sprintf("[]byte { %s }", strings.Join(newData, ","))
}

func CheckIfFileExists(filepath string) error {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return err
	}
	return nil
}

// https://stackoverflow.com/questions/66643946/how-to-remove-duplicates-strings-or-int-from-slice-in-go
func RemoveDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// https://stackoverflow.com/questions/34816489/reverse-slice-of-strings
func ReverseSlice(ss []interface{}) {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
}

func ContainsStringInSlice(s []string, toFind string) bool {
	for _, a := range s {
		if a == toFind {
			return true
		}
	}
	return false
}

func ContainsStringInSliceIgnoreCase(s []string, toFind string) bool {
	for _, a := range s {
		if strings.ToLower(a) == strings.ToLower(toFind) {
			return true
		}
	}
	return false
}

func HasField(v interface{}, name string) bool {
	name = strings.ToLower(name)

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return false
	}
	return rv.FieldByNameFunc(func(n string) bool { return strings.ToLower(n) == strings.ToLower(name) }).IsValid()
}

func GetField(v interface{}, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByNameFunc(func(n string) bool { return strings.ToLower(n) == strings.ToLower(field) })
	return f.String()
}

// https://stackoverflow.com/questions/44255344/using-reflection-setstring/44255582
func SetField(source interface{}, fieldName string, fieldValue string) {
	v := reflect.ValueOf(source).Elem()

	if v.FieldByNameFunc(func(n string) bool { return strings.ToLower(n) == strings.ToLower(fieldName) }).CanSet() {
		v.FieldByNameFunc(func(n string) bool { return strings.ToLower(n) == strings.ToLower(fieldName) }).SetString(fieldValue)
	}
}

func SetDebug(source interface{}, fieldName string, debugValue bool) {
	if HasField(source, "debug") {
		v := reflect.ValueOf(source).Elem()

		if v.FieldByNameFunc(func(n string) bool { return strings.ToLower(n) == strings.ToLower(fieldName) }).CanSet() {
			v.FieldByNameFunc(func(n string) bool { return strings.ToLower(n) == strings.ToLower(fieldName) }).SetBool(debugValue)
		}
	}

}

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(n int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

const letterBytesOnlychar = "abcdefghijklmnopqrstuvwxyz"

func RandomStringOnlyChar(n int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytesOnlychar[rand.Intn(len(letterBytesOnlychar))]
	}
	return string(b)
}

func GetCurrentDate() string {
	// Use layout string for time format.
	const layout = "20060102"
	// Place now in the string.
	t := time.Now()
	return t.Format(layout)
}

// TOD: Rework that function
func CreatePayloadFile(name, ext, source string) (*os.File, error) {

	var path string
	var file *os.File

	extension := ext

	if extension == "" {
		return nil, fmt.Errorf("Unsupported/unspecified language\n")
	}

	if name == "" || len(name) == 0 {

		rand.Seed(time.Now().UnixNano())
		path = fmt.Sprintf("%s_%s_main.%s", GetCurrentDate(), RandomString(4), extension)

	} else {
		path = fmt.Sprintf("%s.%s", name, extension)
	}

	if source != "" || len(source) > 0 {
		path = filepath.Join(source, filepath.Base(path))
	}

	if _, err := os.Stat(path); err == nil {
		//log.Printf("Error file \"%s\" already exists\n", path)
		return nil, err

	} else if os.IsNotExist(err) {
		file, err = os.Create(path)
		if err != nil {
			//log.Printf("create file: ", err)
			return nil, err
		}
	}

	return file, nil
}

func CreateDir(path string) error {
	return os.MkdirAll(path, 0700)
}

func RemoveExt(filename string) string {

	var extension = filepath.Ext(filename)

	if extension != "" {
		return filename[0 : len(filename)-len(extension)]
	}

	return filename
}

func GetProperArch(arch string, lang string) (string, error) {

	switch lang {
	case "go":
		if arch == "x64" {
			return "amd64", nil
		} else if arch == "x86" {
			return "386", nil
		}
	case "c":
		if arch == "x64" {
			return "x64", nil
		}
	case "rust":
		//
	}

	return "", fmt.Errorf("golang arch value must either x86 either x64.")
}

func GetCoreFile(lang string) (string, error) {
	switch lang {
	case "go":
		return "templates/go/core.go.tmpl", nil
	case "c":
		return "templates/c/core.c.tmpl", nil
	case "rust":
		return "", nil
	}
	return "", fmt.Errorf("golang, rust and c are the only supported languages")
}

func GetEnglishWords() []string {
	rawEnglish, err := data.GetTemplates().ReadFile("templates/go/common/english.txt")
	if err != nil {
		return []string{}
	}
	englishWords := strings.Split(string(rawEnglish), "\n")
	return englishWords
}
