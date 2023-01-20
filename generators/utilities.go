package generators

import (
	"bytes"
	"crud-generator/utility"
	"fmt"
	"log"
	"text/template"
)

func GenerateUtilities(sourceConfig GeneratorSource) {
	log.Println("start utilities generation")

	var fileContent bytes.Buffer
	tmpl, err := template.New("utility_errors.tmpl").ParseFiles(fmt.Sprintf("%s%s", defaultTemplatePath, "utility_errors.tmpl"))
	if err != nil {
		log.Println("unable to get utilities template:", err.Error())
	}
	err = tmpl.Execute(&fileContent, sourceConfig)
	if err != nil {
		log.Println("unable to parse utilities template:", err.Error())
	}

	log.Println("utilities generated!")
	utility.WriteFile("utility/errors.go", fileContent.Bytes())
}
