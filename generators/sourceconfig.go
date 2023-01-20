package generators

import (
	"crud-generator/utility"
	"fmt"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
	"strings"
)

const defaultTemplatePath = "templates/"
const deliveryFileName string = "delivery.go"
const repositoryFileName string = "repository.go"
const serviceFileName string = "service.go"
const routesFileName string = "routes.go"
const setupFileName string = "setup.go"

type GeneratorSource struct {
	Service             string                      `yaml:"service"`
	Name                string                      `yaml:"name"`
	Delivery            bool                        `yaml:"delivery"`
	EchoVersion         string                      `yaml:"echoVersion"`
	Template            map[string]bool             `yaml:"template"`
	Attributes          []GeneratorSourceAttributes `yaml:"attributes"`
	Package             string
	PackageSlug         string
	PackageVarTitle     string
	PackageVarLower     string
	InsertValuesString  string
	InsertValuesCounter string
	InsertValues        string
	UpdateSetString     string
	UpdateWhereString   string
	UpdateValues        string
	SelectString        string
	SelectScan          string
	PrimaryKey          string
}

func (gs *GeneratorSource) PrepareForTemplate() {
	gs.Package = utility.ToSnakeCase(gs.Name)
	gs.PackageSlug = utility.ToSlug(gs.Name)
	gs.PackageVarTitle = gs.Name
	gs.PackageVarLower = utility.ReverseTitle(gs.Name)
	gs.PrimaryKey = strings.ToUpper(gs.GetPKeyName())
}

func (gs *GeneratorSource) HasUUIDAsPKey() bool {
	for _, a := range gs.Attributes {
		if strings.ToLower(a.Name) == "uuid" && a.Pkey == true {
			return true
		}
	}
	return false
}

func (gs *GeneratorSource) GetPKeyName() string {
	for _, a := range gs.Attributes {
		if a.Pkey == true {
			return a.Name
		}
	}
	return ""
}

func AttributeToType(aType string) string {
	switch aType {
	case "uuid", "timestamp":
		return "string"
	case "integer":
		return "int64"
	case "float":
		return "float64"
	default:
		return aType
	}
}

func AttributeToPostgresType(aType string, limit int64) string {
	switch aType {
	case "uuid":
		return "UUID"
	case "string":
		if limit > 0 {
			return fmt.Sprintf("VARCHAR(%d)", limit)
		} else {
			return "TEXT"
		}
	case "timestamp":
		return "TIMESTAMP"
	case "integer":
		return "INT"
	case "float":
		return "REAL"
	default:
		return strings.ToUpper(aType)
	}
}

func AttributeToSQLType(aType string) string {
	switch aType {
	case "string", "uuid", "timestamp":
		return "sql.NullString"
	case "integer":
		return "sql.NullInt64"
	case "float":
		return "sql.NullFloat64"
	case "bool":
		return "sql.NullBool"
	default:
		return aType
	}
}

func TypeToSqlGet(aType string) string {
	switch aType {
	case "string", "uuid", "timestamp":
		return "GetStringValue"
	case "integer":
		return "GetIntValue"
	case "float":
		return "GetFloatValue"
	case "bool":
		return "GetBoolValue"
	default:
		return ""
	}
}

func TypeToSqlSet(aType string) string {
	switch aType {
	case "string", "uuid", "timestamp":
		return "NewNullString"
	case "integer":
		return "NewNullInt"
	case "float":
		return "NewNullFloat"
	default:
		return ""
	}
}

type GeneratorSourceAttributes struct {
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Limit       int64  `yaml:"limit,omitempty"`
	Pkey        bool   `yaml:"pkey,omitempty"`
	Required    bool   `yaml:"required"`
	IsReference bool   `yaml:"isReference"`
}

func LoadSource(path string) (sourceConfig GeneratorSource, err error) {
	log.Println("loading source config file")
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err.Error())
		return sourceConfig, err
	}

	log.Println("unmarshal yaml")
	err = yaml.Unmarshal(configFile, &sourceConfig)
	if err != nil {
		log.Println(err.Error())
		return sourceConfig, err
	}

	if val, ok := sourceConfig.Template["generateModifyUserReferences"]; ok && val {
		attr := GeneratorSourceAttributes{
			Name:     "createdBy",
			Type:     "uuid",
			Limit:    0,
			Pkey:     false,
			Required: true,
		}
		sourceConfig.Attributes = append(sourceConfig.Attributes, attr)
		attr = GeneratorSourceAttributes{
			Name:     "modifiedBy",
			Type:     "uuid",
			Limit:    0,
			Pkey:     false,
			Required: false,
		}
		sourceConfig.Attributes = append(sourceConfig.Attributes, attr)
		attr = GeneratorSourceAttributes{
			Name:     "deletedBy",
			Type:     "uuid",
			Limit:    0,
			Pkey:     false,
			Required: false,
		}
		sourceConfig.Attributes = append(sourceConfig.Attributes, attr)
	}

	if val, ok := sourceConfig.Template["generateTimestamps"]; ok && val {
		attr := GeneratorSourceAttributes{
			Name:     "createdDate",
			Type:     "timestamp",
			Limit:    0,
			Pkey:     false,
			Required: true,
		}
		sourceConfig.Attributes = append(sourceConfig.Attributes, attr)
		attr = GeneratorSourceAttributes{
			Name:     "lastModified",
			Type:     "timestamp",
			Limit:    0,
			Pkey:     false,
			Required: false,
		}
		sourceConfig.Attributes = append(sourceConfig.Attributes, attr)
		attr = GeneratorSourceAttributes{
			Name:     "deletedDate",
			Type:     "timestamp",
			Limit:    0,
			Pkey:     false,
			Required: false,
		}
		sourceConfig.Attributes = append(sourceConfig.Attributes, attr)
	}

	log.Println("source config loaded")
	return sourceConfig, nil
}
