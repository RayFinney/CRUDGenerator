package repository

import (
	"crud-generator/models"
	"crud-generator/utility"
	"fmt"
	"log"
	"strings"
	"sync"
)

const fileName string = "repository.go"
const getTemplate string = `
func (r *Repository) GetByID(uuid string) (models.%s, error) {
	query := "SELECT %s FROM %s WHERE %s.uuid = $1"
	return r.getOne(query, uuid)
}
`
const getAllTemplate string = `
func (r *Repository) GetAll() ([]models.%s, error) {
	query := "SELECT %s FROM %s"
	return r.fetch(query)
}
`
const postTemplate string = `
func (r *Repository) Store(%s *models.%s) error {
	query := "INSERT INTO %s (%s) VALUES (%s)"
	_, err := r.dbClient.Exec(query, %s)
	return err
}
`
const putTemplate string = `
func (r *Repository) Update(uuid string, %s *models.%s) error {
	query := "UPDATE %s SET %s WHERE %s.uuid = $%d"
	_, err := r.dbClient.Exec(query, %s, uuid)
	return err
}
`
const deleteTemplate string = `
func (r *Repository) Delete(uuid string) error {
	query := "DELETE FROM %s WHERE %s.uuid = $1"
	_, err := r.dbClient.Exec(query, uuid)
	return err
}
`
const factoryTemplate string = `
func NewRepository(dbClient *sql.DB) Repository {
	return Repository{dbClient: dbClient}
}`
const fetchTemplate string = `
func (r *Repository) fetch(query string, args ...interface{}) ([]models.%s, error) {
	rows, err := r.dbClient.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Warningf("db connection could not be closed: %s", err)
		}
	}()
	result := make([]models.%s, 0)
	for rows.Next() {
		%sDB := models.%sDB{}
		err := rows.Scan(%s)
		if err != nil && err != sql.ErrNoRows {
			log.Infof("error while loading data: %s", err)
			return result, err
		}
		result = append(result, %sDB.Get%s())
	}
	return result, nil
}
`
const getOneTemplate string = `
func (r *Repository) getOne(query string, args ...interface{}) (models.%s, error) {
	%sDB := models.%sDB{}
	err := r.dbClient.QueryRow(query, args...).Scan(%s)
	if err != nil && err != sql.ErrNoRows {
		log.Infof("error while loading data: %s", err)
	}
	return %sDB.Get%s(), err
}
`

func GenerateRepository(sourceConfig models.GeneratorSource, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("start repository generation")
	var fileContent string
	lowercaseName := utility.ReverseTitle(sourceConfig.Name)
	snakeCaseName := fmt.Sprintf("%ss", utility.ToSnakeCase(sourceConfig.Name))

	// HEADER
	generateDeliveryHeader(&fileContent, sourceConfig)

	// STRUCT
	fileContent += "type Repository struct {\n"
	fileContent += "\t dbClient *sql.DB\n"
	fileContent += "}\n\n"

	var selectString string
	for index, a := range sourceConfig.Attributes {
		selectString += fmt.Sprintf("%s.%s", snakeCaseName, utility.ToSnakeCase(a.Name))
		if len(sourceConfig.Attributes)-1 > index {
			selectString += ", "
		}
	}

	// GET
	fileContent += fmt.Sprintf(getTemplate, sourceConfig.Name, selectString, snakeCaseName, snakeCaseName)

	// GET_ALL
	fileContent += fmt.Sprintf(getAllTemplate, sourceConfig.Name, selectString, snakeCaseName)

	// POST
	var dbvalues string
	var ivalues string
	var values string

	for index, a := range sourceConfig.Attributes {
		name := strings.Title(a.Name)
		if name == "Uuid" {
			name = "UUID"
		}
		dbvalues += utility.ToSnakeCase(a.Name)
		ivalues += fmt.Sprintf("$%d", index+1)
		if models.TypeToSqlSet(a.Type) != "" {
			values += fmt.Sprintf("utility.%s(%s.%s)", models.TypeToSqlSet(a.Type), lowercaseName, name)
		} else {
			values += fmt.Sprintf("%s.%s", lowercaseName, name)
		}

		if len(sourceConfig.Attributes)-1 > index {
			dbvalues += ", "
			ivalues += ", "
			values += ", "
		}
	}

	fileContent += fmt.Sprintf(postTemplate, lowercaseName, sourceConfig.Name, snakeCaseName, dbvalues, ivalues, values)

	// PUT
	var setString string
	var uValues string
	counter := 1
	for index, a := range sourceConfig.Attributes {
		if utility.ToSnakeCase(a.Name) == "uuid" || utility.ToSnakeCase(a.Name) == "created_by" || utility.ToSnakeCase(a.Name) == "deleted_by" ||
			utility.ToSnakeCase(a.Name) == "created_date" || utility.ToSnakeCase(a.Name) == "deleted_date" {
			continue
		}
		name := strings.Title(a.Name)
		if name == "Uuid" {
			name = "UUID"
		}
		if utility.ToSnakeCase(a.Name) == "last_modified" {
			setString += fmt.Sprintf("%s = CURRENT_TIMESTAMP", utility.ToSnakeCase(a.Name))
		} else {
			setString += fmt.Sprintf("%s = $%d", utility.ToSnakeCase(a.Name), counter)
			if models.TypeToSqlSet(a.Type) != "" {
				uValues += fmt.Sprintf("utility.%s(%s.%s)", models.TypeToSqlSet(a.Type), lowercaseName, name)
			} else {
				uValues += fmt.Sprintf("%s.%s", lowercaseName, name)
			}
			if len(sourceConfig.Attributes)-1 > index {
				setString += ", "
				uValues += ", "
			}
			counter++
		}
	}
	if string(uValues[len(uValues)-2]) == "," {
		uValues = uValues[0 : len(uValues)-2]
	}

	fileContent += fmt.Sprintf(putTemplate, lowercaseName, sourceConfig.Name, snakeCaseName, setString, snakeCaseName, counter, uValues)

	// DELETE
	fileContent += fmt.Sprintf(deleteTemplate, snakeCaseName, snakeCaseName)

	var scanString string
	for index, a := range sourceConfig.Attributes {
		name := strings.Title(a.Name)
		if name == "Uuid" {
			name = "UUID"
		}
		if index%5 == 0 && index > 0 {
			scanString += "\n\t\t"
		}
		scanString += fmt.Sprintf("&%sDB.%s", lowercaseName, name)
		if len(sourceConfig.Attributes)-1 > index {
			scanString += ", "
		}
	}

	// GETONE
	fileContent += fmt.Sprintf(getOneTemplate, sourceConfig.Name, lowercaseName, sourceConfig.Name, scanString, "%v", lowercaseName, sourceConfig.Name)

	// FETCH
	fileContent += fmt.Sprintf(fetchTemplate, sourceConfig.Name, "%v", sourceConfig.Name, lowercaseName, sourceConfig.Name, scanString, "%v", lowercaseName, sourceConfig.Name)

	// FACTORY
	fileContent += factoryTemplate

	log.Println("repository generated!")
	utility.WriteFile(fmt.Sprintf("%s/%s", utility.ToSnakeCase(sourceConfig.Name), fileName), []byte(fileContent))
}

func generateDeliveryHeader(fileContent *string, sourceConfig models.GeneratorSource) {
	*fileContent += utility.GeneratePackage(utility.ToSnakeCase(sourceConfig.Name))
	*fileContent += utility.GenerateImports([]string{fmt.Sprintf("%s/models", sourceConfig.Service)},
		[]string{fmt.Sprintf("%s/utility", sourceConfig.Service)}, []string{"database/sql"}, []string{"log", "github.com/sirupsen/logrus"})
}
