package main

import (
	"crud-generator/generators"
	"crud-generator/utility"
	"flag"
	"log"
	"os"
	"path/filepath"
)

func main() {
	sourceConfigPath := flag.String("config", "", "path to config yaml")
	flag.Parse()
	if *sourceConfigPath == "" {
		log.Fatal("config path is required")
	}

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	generators.DefaultTemplatePath = exPath + generators.DefaultTemplatePath

	sourceConfig, err := generators.LoadSource(*sourceConfigPath)
	if err != nil {
		os.Exit(1)
	}
	sourceConfig.PrepareForTemplate()
	generateStructure(sourceConfig)
}

func generateStructure(sourceConfig generators.GeneratorSource) {
	err := os.Mkdir(utility.ToSnakeCase(sourceConfig.Name), os.ModePerm)
	if err != nil {
		log.Fatalf("[generateStructure] unable to create folder(%s): %v", sourceConfig.Name, err)
	}

	generators.GenerateModel(sourceConfig)
	generators.GenerateOpenApiSchema(sourceConfig)
	if sourceConfig.Delivery {
		generators.GenerateDelivery(sourceConfig)
		generators.GenerateRoutes(sourceConfig)
		generators.GeneratePinia(sourceConfig)
		generators.GenerateOpenApiPaths(sourceConfig)
	}

	generators.GenerateService(sourceConfig)
	generators.GenerateRepository(sourceConfig)
	generators.GenerateMigration(sourceConfig)
	generators.GenerateSetup(sourceConfig)
	generators.GenerateUtilities(sourceConfig)
}
