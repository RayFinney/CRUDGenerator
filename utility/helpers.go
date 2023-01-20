package utility

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	str = strings.TrimSpace(str)
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func ToSlug(str string) string {
	str = strings.TrimSpace(str)
	snake := matchFirstCap.ReplaceAllString(str, "${1}-${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}-${2}")
	return strings.ToLower(snake)
}

func Pluralize(str string) string {
	str = strings.TrimSpace(str)
	if string(str[len(str)-1]) == "y" {
		str = fmt.Sprintf("%sies", str[:len(str)-1])
	} else if string(str[len(str)-1]) != "s" {
		str = fmt.Sprintf("%ss", str)
	}
	return str
}

func WriteFile(path string, data []byte) {
	dirName := filepath.Dir(path)
	log.Println("make folder path (", path, ") if not existing")
	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		log.Fatalf("[writeFile] unable to create folder(%s): %v", dirName, err)
	}
	log.Println("write file", path)
	err = ioutil.WriteFile(path, data, os.ModePerm)
	if err != nil {
		log.Fatalf("[writeFile] unable to write file(%s): %v", path, err)
	}
}

func ReverseTitle(str string) string {
	str = strings.TrimSpace(str)
	firstLetter := strings.ToLower(string(str[0:1]))
	return fmt.Sprintf("%s%s", firstLetter, string(str[1:]))
}
