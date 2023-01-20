package generators

import (
	"bytes"
	"crud-generator/utility"
	"fmt"
	"log"
	"path/filepath"
	"text/template"
)

func GenerateUtilities(sourceConfig GeneratorSource) {
	log.Println("start utilities generation")

	var fileContent bytes.Buffer
	path, _ := filepath.Abs(fmt.Sprintf("%s%s", DefaultTemplatePath, "utility_errors.tmpl"))
	tmpl, err := template.New("utility_errors.tmpl").ParseFiles(path)
	if err != nil {
		log.Println("unable to get utilities template:", err.Error())
	}
	err = tmpl.Execute(&fileContent, sourceConfig)
	if err != nil {
		log.Println("unable to parse utilities template:", err.Error())
	}
	utility.WriteFile("utility/errors.go", fileContent.Bytes())

	path, _ = filepath.Abs(fmt.Sprintf("%s%s", DefaultTemplatePath, "results_list.tmpl"))
	tmpl, err = template.New("results_list.tmpl").ParseFiles(path)
	if err != nil {
		log.Println("unable to get results_list.tmpl template:", err.Error())
	}
	fileContent = bytes.Buffer{}
	err = tmpl.Execute(&fileContent, sourceConfig)
	if err != nil {
		log.Println("unable to parse results_list.tmpl template:", err.Error())
	}
	utility.WriteFile("models/results_list.go", fileContent.Bytes())

	log.Println("utilities generated!")
}
