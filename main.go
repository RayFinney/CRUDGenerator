package main

import (
	"crud-generator/delivery"
	"crud-generator/migrate"
	"crud-generator/model"
	"crud-generator/models"
	"crud-generator/repository"
	"crud-generator/service"
	"crud-generator/utility"
	"flag"
	"log"
	"os"
	"sync"
)

var wg sync.WaitGroup

func main() {
	sourceConfigPath := flag.String("domain", "", "Domain name")
	flag.Parse()
	if *sourceConfigPath == "" {
		log.Fatal("config path is required")
	}

	sourceConfig, err := models.LoadSource("generator_sources/news_article.yaml")
	if err != nil {
		os.Exit(1)
	}
	generateStructure(sourceConfig)
}

func generateStructure(sourceConfig models.GeneratorSource) {
	err := os.Mkdir(utility.ToSnakeCase(sourceConfig.Name), os.ModePerm)
	if err != nil {
		log.Fatalf("[generateStructure] unable to create folder(%s): %v", sourceConfig.Name, err)
	}
	wg = sync.WaitGroup{}

	wg.Add(1)
	go model.GenerateModel(sourceConfig, &wg)

	if sourceConfig.Delivery {
		wg.Add(1)
		go delivery.GenerateDelivery(sourceConfig, &wg)
	}

	wg.Add(1)
	go service.GenerateService(sourceConfig, &wg)

	wg.Add(1)
	go repository.GenerateRepository(sourceConfig, &wg)

	wg.Add(1)
	go migrate.GenerateMigration(sourceConfig, &wg)

	wg.Wait()
}
