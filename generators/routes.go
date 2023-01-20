package generators

import (
	"bytes"
	"crud-generator/utility"
	"fmt"
	"log"
	"text/template"
)

func GenerateRoutes(sourceConfig GeneratorSource) {
	log.Println("start routes generation")

	var fileContent bytes.Buffer
	tmpl, err := template.New("routes.tmpl").ParseFiles(fmt.Sprintf("%s%s", defaultTemplatePath, "routes.tmpl"))
	if err != nil {
		log.Println("unable to get routes template:", err.Error())
	}
	err = tmpl.Execute(&fileContent, sourceConfig)
	if err != nil {
		log.Println("unable to parse routes template:", err.Error())
	}

	log.Println("routes generated!")
	utility.WriteFile(fmt.Sprintf("%s/%s", utility.ToSnakeCase(sourceConfig.Package), routesFileName), fileContent.Bytes())
}
