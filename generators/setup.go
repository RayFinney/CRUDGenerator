package generators

import (
	"bytes"
	"crud-generator/utility"
	"fmt"
	"log"
	"path/filepath"
	"text/template"
)

func GenerateSetup(sourceConfig GeneratorSource) {
	log.Println("start setup generation")

	var fileContent bytes.Buffer
	path, _ := filepath.Abs(fmt.Sprintf("%s%s", DefaultTemplatePath, "setup.tmpl"))
	tmpl, err := template.New("setup.tmpl").ParseFiles(path)
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
