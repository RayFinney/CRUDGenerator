package generators

import (
	"crud-generator/utility"
	"fmt"
	"github.com/go-yaml/yaml"
	"log"
)

func GenerateOpenApiSchema(sourceConfig GeneratorSource) {
	log.Println("start openapi schema generation")

	schema := make(map[string]interface{})
	schema[sourceConfig.Name] = make(map[string]interface{})
	schema[sourceConfig.Name].(map[string]interface{})["type"] = "object"
	schema[sourceConfig.Name].(map[string]interface{})["required"] = getRequiredAttrs(sourceConfig)
	schema[sourceConfig.Name].(map[string]interface{})["properties"] = getProperties(sourceConfig)

	fileContent, err := yaml.Marshal(&schema)
	if err != nil {
		log.Println("unable to create yaml schema:", err.Error())
	}

	log.Println("openapi schema generated!")
	utility.WriteFile(fmt.Sprintf("api/%s%s", sourceConfig.Name, openapiSchemaFileName), fileContent)
}

func getRequiredAttrs(sourceConfig GeneratorSource) []string {
	requiredAttrs := make([]string, 0)
	for _, a := range sourceConfig.Attributes {
		if a.Required {
			requiredAttrs = append(requiredAttrs, a.Name)
		}
	}
	return requiredAttrs
}

func getProperties(sourceConfig GeneratorSource) map[string]interface{} {
	properties := make(map[string]interface{})
	for _, a := range sourceConfig.Attributes {
		properties[a.Name] = make(map[string]interface{})
		properties[a.Name].(map[string]interface{})["type"] = a.Type
		if a.Limit != 0 && a.Type == "string" {
			properties[a.Name].(map[string]interface{})["maxLength"] = a.Limit
		}

	}
	return properties
}
