package generators

import (
	"bytes"
	"crud-generator/utility"
	"fmt"
	"log"
	"text/template"
)

func GenerateService(sourceConfig GeneratorSource) {
	log.Println("start service generation")

	var fileContent bytes.Buffer
	tmpl, err := template.New("service.tmpl").ParseFiles(fmt.Sprintf("%s%s", defaultTemplatePath, "service.tmpl"))
	if err != nil {
		log.Println("unable to get service template:", err.Error())
	}
	err = tmpl.Execute(&fileContent, sourceConfig)
	if err != nil {
		log.Println("unable to parse service template:", err.Error())
	}

	log.Println("service generated!")
	utility.WriteFile(fmt.Sprintf("%s/%s", utility.ToSnakeCase(sourceConfig.Package), serviceFileName), fileContent.Bytes())
}
