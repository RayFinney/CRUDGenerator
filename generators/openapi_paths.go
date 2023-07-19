package generators

import (
	"bytes"
	"crud-generator/utility"
	"fmt"
	"log"
	"path/filepath"
	"text/template"
)

func GenerateOpenApiPaths(sourceConfig GeneratorSource) {
	log.Println("start openapi paths generation")

	var fileContent bytes.Buffer
	path, _ := filepath.Abs(fmt.Sprintf("%s%s", DefaultTemplatePath, "openapi_paths.tmpl"))
	tmpl, err := template.New("openapi_paths.tmpl").ParseFiles(path)
	if err != nil {
		log.Println("unable to get openapi paths template:", err.Error())
	}
	err = tmpl.Execute(&fileContent, sourceConfig)
	if err != nil {
		log.Println("unable to parse openapi paths  template:", err.Error())
	}

	log.Println("openapi paths generated!")
	utility.WriteFile(fmt.Sprintf("api/%s%s", sourceConfig.Name, openapiPathsFileName), fileContent.Bytes())
}
