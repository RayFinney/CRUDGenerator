package generators

import (
	"crud-generator/utility"
	"fmt"
	"log"
	"strings"
)

func GenerateModel(sourceConfig GeneratorSource) {
	log.Println("start model generation")
	var fileContent string

	// HEADER
	fileContent += utility.GeneratePackage("models")
	fileContent += utility.GenerateImports([]string{"errors"})

	// STRUCT
	fileContent += fmt.Sprintf("type %s struct {\n", sourceConfig.Name)
	for _, a := range sourceConfig.Attributes {
		name := strings.Title(a.Name)
		if name == "Uuid" {
			name = "UUID"
		}
		fileContent += fmt.Sprintf("\t%s %s `json:\"%s\"`\n", name, AttributeToType(a.Type), utility.ToSnakeCase(name))
	}
	fileContent += "}\n\n"

	// VALIDATOR
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

	log.Println("model generated!")
	utility.WriteFile(fmt.Sprintf("models/%s.go", utility.ToSnakeCase(sourceConfig.Name)), []byte(fileContent))
}

func CanApplyRequired(aType string) bool {
	switch AttributeToType(aType) {
	case "string", "int64", "float64":
		return true
	}
	return false
}

func CanApplyLimit(aType string) bool {
	switch AttributeToType(aType) {
	case "string", "int64", "float64":
		return true
	}
	return false
}

func GetRequireIf(name string, aType string) string {
	switch AttributeToType(aType) {
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
	switch AttributeToType(aType) {
	case "string":
		return fmt.Sprintf("if len(v.%s) > %d {", name, limit)
	case "int64":
		return fmt.Sprintf("if v.%s  > %d {", name, limit)
	case "float64":
		return fmt.Sprintf("if v.%s  > %f {", name, float64(limit))
	}
	return ""
}
