package generators

import (
	"bytes"
	"crud-generator/utility"
	"fmt"
	"log"
	"path/filepath"
	"text/template"
)

func GenerateDelivery(sourceConfig GeneratorSource) {
	log.Println("start delivery generation")

	var fileContent bytes.Buffer
	path, _ := filepath.Abs(fmt.Sprintf("%s%s", DefaultTemplatePath, "delivery.tmpl"))
	tmpl, err := template.New("delivery.tmpl").ParseFiles(path)
	if err != nil {
		log.Println("unable to get delivery template:", err.Error())
	}
	err = tmpl.Execute(&fileContent, sourceConfig)
	if err != nil {
		log.Println("unable to parse delivery template:", err.Error())
	}

	log.Println("delivery generated!")
	utility.WriteFile(fmt.Sprintf("%s/%s", utility.ToSnakeCase(sourceConfig.Package), deliveryFileName), fileContent.Bytes())
}
