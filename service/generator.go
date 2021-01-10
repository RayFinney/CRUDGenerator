package service

import (
	"crud-generator/models"
	"crud-generator/utility"
	"fmt"
	"log"
	"sync"
)

const fileName string = "service.go"
const getTemplate string = `
func (s *Service) GetByID(ctx context.Context, uuid string) (%s models.%s, err error) {
	%s, err = s.%s.GetByID(uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return %s, utility.NOT_FOUND
		}
		log.Warningf("[%sService.GetByID()] unable to load %s: %s", err)
		return %s, utility.DATABASE_ERROR
	}
	return %s, nil
}
`
const postTemplate string = `
func (s *Service) Store(ctx context.Context, %s *models.%s) error {
	%s
	if err := s.%s.Store(%s); err != nil {
		log.Warningf("[%sService.Store()] unable to store %s: %s", err)
		return utility.DATABASE_ERROR
	}
	*%s, _ = s.GetByID(ctx, %s.UUID)
	return nil
}
`
const putTemplate string = `
func (s *Service) Update(ctx context.Context, uuid string, %s *models.%s) error {
	if err := s.%s.Update(uuid, %s); err != nil {
		log.Warningf("[%sService.Update()] unable to update %s: %s", err)
		return utility.DATABASE_ERROR
	}
	*%s, _ = s.GetByID(ctx, %s.UUID)
	return nil
}
`
const deleteTemplate string = `
func (s *Service) Delete(ctx context.Context, uuid string) error {
	if err := s.%s.Delete(uuid); err != nil {
		if err == sql.ErrNoRows {
			return utility.NOT_FOUND
		}
		log.Warningf("[%sService.Delete()] unable to delete %s(%s): %s", uuid, err)
		return utility.DATABASE_ERROR
	}
	return nil
}
`
const factoryTemplate string = `
func NewService(%s Repository) Service {
	return Service{%s: %s}
}`

func GenerateService(sourceConfig models.GeneratorSource, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("start service generation")
	var fileContent string
	lowercaseName := utility.ReverseTitle(sourceConfig.Name)
	serviceName := fmt.Sprintf("%sService", lowercaseName)
	repoName := fmt.Sprintf("%sRepo", lowercaseName)

	// HEADER
	generateDeliveryHeader(&fileContent, sourceConfig)

	// STRUCT
	fileContent += "type Service struct {\n"
	fileContent += fmt.Sprintf("\t %s Repository\n", repoName)
	fileContent += "}\n\n"

	// GET
	fileContent += fmt.Sprintf(getTemplate, lowercaseName, sourceConfig.Name, lowercaseName, repoName, lowercaseName, lowercaseName, sourceConfig.Name, "%v", lowercaseName, lowercaseName)

	// POST
	uuidGenTemp := ""
	if sourceConfig.HasUUIDAsPKey() {
		uuidGenTemp = fmt.Sprintf(`if %s.UUID == "" {
		%s.UUID = uuid.NewV4().String()
	}`, lowercaseName, lowercaseName)
	}
	fileContent += fmt.Sprintf(postTemplate, lowercaseName, sourceConfig.Name, uuidGenTemp, repoName, lowercaseName,
		lowercaseName, sourceConfig.Name, "%v", lowercaseName, lowercaseName)

	// PUT
	fileContent += fmt.Sprintf(putTemplate, lowercaseName, sourceConfig.Name, repoName, lowercaseName, lowercaseName, sourceConfig.Name, "%v",
		lowercaseName, lowercaseName)

	// DELETE
	fileContent += fmt.Sprintf(deleteTemplate, repoName, lowercaseName, sourceConfig.Name, "%s", "%v")

	// FACTORY
	fileContent += fmt.Sprintf(factoryTemplate, serviceName, serviceName, serviceName)

	log.Println("service generated!")
	utility.WriteFile(fmt.Sprintf("%s/%s", utility.ToSnakeCase(sourceConfig.Name), fileName), []byte(fileContent))
}

func generateDeliveryHeader(fileContent *string, sourceConfig models.GeneratorSource) {
	*fileContent += utility.GeneratePackage(utility.ToSnakeCase(sourceConfig.Name))
	*fileContent += utility.GenerateImports([]string{"context"}, []string{"database/sql"}, []string{"github.com/satori/go.uuid"},
		[]string{"log", "github.com/sirupsen/logrus"}, []string{fmt.Sprintf("%s/models", sourceConfig.Service)},
		[]string{fmt.Sprintf("%s/utility", sourceConfig.Service)})
}
