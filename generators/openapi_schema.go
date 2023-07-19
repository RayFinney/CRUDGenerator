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
		for k, v := range getType(a.Type) {
			properties[a.Name].(map[string]interface{})[k] = v
		}
		if a.Limit != 0 && a.Type == "string" {
			properties[a.Name].(map[string]interface{})["maxLength"] = a.Limit
		}

	}
	return properties
}

func getType(aType string) map[string]interface{} {
	switch aType {
	case "timestamp":
		return map[string]interface{}{"type": "string", "format": "timestamp"}
	case "uuid":
		return map[string]interface{}{"type": "string", "format": "uuid"}
	case "float":
		return map[string]interface{}{"type": "number", "format": "float"}
	case "bool":
		return map[string]interface{}{"type": "boolean"}
	default:
		return map[string]interface{}{"type": aType}
	}
}
