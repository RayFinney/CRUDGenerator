package generators

import (
	"bytes"
	"crud-generator/utility"
	"fmt"
	"log"
	"path/filepath"
	"text/template"
)

func GeneratePinia(sourceConfig GeneratorSource) {
	log.Println("start pinia generation")

	var fileContent bytes.Buffer
	path, _ := filepath.Abs(fmt.Sprintf("%s%s", DefaultTemplatePath, "pinia.tmpl"))
	tmpl, err := template.New("pinia.tmpl").ParseFiles(path)
	if err != nil {
		log.Println("unable to get pinia template:", err.Error())
	}
	err = tmpl.Execute(&fileContent, sourceConfig)
	if err != nil {
		log.Println("unable to parse pinia template:", err.Error())
	}

	log.Println("pinia generated!")
	utility.WriteFile(fmt.Sprintf("%s/%s", utility.ToSnakeCase(sourceConfig.Package), fmt.Sprintf("%s.js", sourceConfig.PackageVarLower)), fileContent.Bytes())
}
