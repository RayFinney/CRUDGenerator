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
	fileContent += utility.GenerateImports([]string{fmt.Sprintf("%s/utility", sourceConfig.Service)}, []string{"database/sql"}, []string{"errors"})

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

	// VALIDATOR
	// DB STRUCT MAPPER
	fileContent += fmt.Sprintf("func (v *%s) Validate() error {\n", sourceConfig.Name)
	for _, a := range sourceConfig.Attributes {
		name := strings.Title(a.Name)
		if name == "Uuid" {
			name = "UUID"
		}
		if a.Required && CanApplyRequired(a.Type) {
			fileContent += fmt.Sprintf("\t%s\n", GetRequireIf(name, a.Type))
			fileContent += fmt.Sprintf("\t\treturn errors.New(\"%s required\")\n", name)
			fileContent += "\t}\n"
		}
		if a.Limit > 0 && CanApplyLimit(a.Type) {
			fileContent += fmt.Sprintf("\t%s\n", GetLimitIf(name, a.Limit, a.Type))
			fileContent += fmt.Sprintf("\t\treturn errors.New(\"%s to large (max %d)\")\n", name, a.Limit)
			fileContent += "\t}\n"
		}
	}
	fileContent += "\treturn nil\n"
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

func CanApplyRequired(aType string) bool {
	switch models.AttributeToType(aType) {
	case "string", "int64", "float64":
		return true
	}
	return false
}

func CanApplyLimit(aType string) bool {
	switch models.AttributeToType(aType) {
	case "string", "int64", "float64":
		return true
	}
	return false
}

func GetRequireIf(name string, aType string) string {
	switch models.AttributeToType(aType) {
	case "string":
		return fmt.Sprintf("if v.%s == \"\" {", name)
	case "int64":
		return fmt.Sprintf("if v.%s == 0 {", name)
	case "float64":
		return fmt.Sprintf("if v.%s == 0.0 {", name)
	}
	return ""
}

func GetLimitIf(name string, limit int64, aType string) string {
	switch models.AttributeToType(aType) {
	case "string":
		return fmt.Sprintf("if len(v.%s) > %d {", name, limit)
	case "int64":
		return fmt.Sprintf("if v.%s  > %d {", name, limit)
	case "float64":
		return fmt.Sprintf("if v.%s  > %f {", name, float64(limit))
	}
	return ""
}
