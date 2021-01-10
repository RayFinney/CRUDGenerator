package model

import (
	"crud-generator/models"
	"crud-generator/utility"
	"fmt"
	"log"
	"strings"
	"sync"
)

func GenerateModel(sourceConfig models.GeneratorSource, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("start model generation")
	var fileContent string

	// HEADER
	fileContent += utility.GeneratePackage("models")
	fileContent += utility.GenerateImports([]string{fmt.Sprintf("%s/utility", sourceConfig.Service)}, []string{"database/sql"})

	// STRUCT
	fileContent += fmt.Sprintf("type %s struct {\n", sourceConfig.Name)
	for _, a := range sourceConfig.Attributes {
		name := strings.Title(a.Name)
		if name == "Uuid" {
			name = "UUID"
		}
		fileContent += fmt.Sprintf("\t%s %s `json:\"%s\"`\n", name, models.AttributeToType(a.Type), utility.ToSnakeCase(name))
	}
	fileContent += "}\n\n"

	// DB STRUCT
	fileContent += fmt.Sprintf("type %sDB struct {\n", sourceConfig.Name)
	for _, a := range sourceConfig.Attributes {
		name := strings.Title(a.Name)
		if name == "Uuid" {
			name = "UUID"
		}
		fileContent += fmt.Sprintf("\t%s %s\n", name, models.AttributeToSQLType(a.Type))
	}
	fileContent += "}\n\n"

	// DB STRUCT MAPPER
	fileContent += fmt.Sprintf("func (dbv *%sDB) Get%s() (v %s) {\n", sourceConfig.Name, sourceConfig.Name, sourceConfig.Name)
	fileContent += "\t utility.MapSqlValues(dbv, &v)\n"
	//for _, a := range sourceConfig.Attributes {
	//	name := strings.Title(a.Name)
	//	if name == "Uuid" {
	//		name = "UUID"
	//	}
	//	fileContent += fmt.Sprintf("\tv.%s = utility.%s(dbv.%s)\n", name, models.TypeToSqlGet(a.Type), name)
	//}
	fileContent += "\treturn v\n"
	fileContent += "}\n\n"

	log.Println("model generated!")
	utility.WriteFile(fmt.Sprintf("models/%s.go", utility.ToSnakeCase(sourceConfig.Name)), []byte(fileContent))
}
