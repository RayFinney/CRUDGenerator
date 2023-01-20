package generators

import (
	"bytes"
	"crud-generator/utility"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"text/template"
)

func GenerateRepository(sourceConfig GeneratorSource) {
	log.Println("start repository generation")

	for index, a := range sourceConfig.Attributes {
		sourceConfig.SelectString += utility.ToSnakeCase(a.Name)
		if len(sourceConfig.Attributes)-1 > index {
			sourceConfig.SelectString += ", "
		}
	}

	counter := 1
	for index, a := range sourceConfig.Attributes {
		name := strings.Title(a.Name)
		if name == "Uuid" {
			name = "UUID"
		}
		sourceConfig.SelectScan += fmt.Sprintf("%s.%s", sourceConfig.PackageVarLower, name)

		// INSERT
		sourceConfig.InsertValuesString += utility.ToSnakeCase(a.Name)
		sourceConfig.InsertValuesCounter += fmt.Sprintf("$%d", index+1)
		sourceConfig.InsertValues += fmt.Sprintf("%s.%s", sourceConfig.PackageVarLower, name)

		if len(sourceConfig.Attributes)-1 > index {
			sourceConfig.InsertValuesString += ", "
			sourceConfig.InsertValuesCounter += ", "
			sourceConfig.InsertValues += ", "
			sourceConfig.SelectScan += ", "
		}

		// UPDATE
		if utility.ToSnakeCase(a.Name) == "uuid" || utility.ToSnakeCase(a.Name) == "created_by" || utility.ToSnakeCase(a.Name) == "deleted_by" ||
			utility.ToSnakeCase(a.Name) == "created_date" || utility.ToSnakeCase(a.Name) == "deleted_date" {
			continue
		}
		if utility.ToSnakeCase(a.Name) == "last_modified" {
			sourceConfig.UpdateSetString += fmt.Sprintf("%s = CURRENT_TIMESTAMP", utility.ToSnakeCase(a.Name))
		} else {
			sourceConfig.UpdateSetString += fmt.Sprintf("%s = $%d", utility.ToSnakeCase(a.Name), counter)
			sourceConfig.UpdateValues += fmt.Sprintf("%s.%s", sourceConfig.PackageVarLower, name)
			if len(sourceConfig.Attributes)-1 > index {
				sourceConfig.UpdateSetString += ", "
				sourceConfig.UpdateValues += ", "
			}
			counter++
		}
	}
	if string(sourceConfig.UpdateValues[len(sourceConfig.UpdateValues)-2]) == "," {
		sourceConfig.UpdateValues = sourceConfig.UpdateValues[0 : len(sourceConfig.UpdateValues)-2]
	}
	sourceConfig.UpdateWhereString = fmt.Sprintf("%s = $%d", utility.ToSnakeCase(sourceConfig.GetPKeyName()), counter)

	var fileContent bytes.Buffer
	path, _ := filepath.Abs(fmt.Sprintf("%s%s", DefaultTemplatePath, "repository.tmpl"))
	tmpl, err := template.New("repository.tmpl").ParseFiles(path)
	if err != nil {
		log.Println("unable to get repository template:", err.Error())
	}
	err = tmpl.Execute(&fileContent, sourceConfig)
	if err != nil {
		log.Println("unable to parse repository template:", err.Error())
	}

	log.Println("repository generated!")
	utility.WriteFile(fmt.Sprintf("%s/%s", utility.ToSnakeCase(sourceConfig.Package), repositoryFileName), fileContent.Bytes())
}
