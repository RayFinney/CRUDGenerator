package migrate

import (
	"crud-generator/models"
	"crud-generator/utility"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

func GenerateMigration(sourceConfig models.GeneratorSource, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("start migration generation")
	var fileContent string
	snakeCaseName := utility.ToSnakeCase(utility.Pluralize(sourceConfig.Name))

	fileContent += fmt.Sprintf("CREATE TABLE %s\n", snakeCaseName)
	fileContent += "(\n"
	primaryKey := ""
	for _, a := range sourceConfig.Attributes {
		null := "NULL"
		if a.Required {
			null = "NOT NULL"
		}
		fileContent += fmt.Sprintf("\t%s %s %s,\n", utility.ToSnakeCase(a.Name), strings.ToUpper(models.AttributeToPostgresType(a.Type, a.Limit)), null)
		if a.Pkey {
			primaryKey = fmt.Sprintf("\tPRIMARY KEY (%s)\n", strings.ToLower(a.Name))
		}
	}
	if primaryKey != "" {
		fileContent += primaryKey
	}
	fileContent += ");"
	log.Println("migration generated!")
	fileName := fmt.Sprintf("%s-create-table-%s.sql", time.Now().Format("2006-01-02"), snakeCaseName)
	utility.WriteFile(fmt.Sprintf("migrations/%s", fileName), []byte(fileContent))
}
