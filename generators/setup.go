package generators

import (
	"bytes"
	"crud-generator/utility"
	"fmt"
	"log"
	"text/template"
)

func GenerateSetup(sourceConfig GeneratorSource) {
	log.Println("start setup generation")

	var fileContent bytes.Buffer
	tmpl, err := template.New("setup.tmpl").ParseFiles(fmt.Sprintf("%s%s", defaultTemplatePath, "setup.tmpl"))
	if err != nil {
		log.Println("unable to get setup template:", err.Error())
	}
	err = tmpl.Execute(&fileContent, sourceConfig)
	if err != nil {
		log.Println("unable to parse setup template:", err.Error())
	}

	log.Println("setup generated!")
	utility.WriteFile(fmt.Sprintf("%s/%s", utility.ToSnakeCase(sourceConfig.Package), setupFileName), fileContent.Bytes())
}
