package delivery

import (
	"crud-generator/models"
	"crud-generator/utility"
	"fmt"
	"log"
	"strings"
	"sync"
)

const fileName string = "delivery.go"
const getTemplate string = `
func (d *Delivery) GetByID(c echo.Context) error {
	uuid := c.Param("uuid")
	%s, err := d.%s.GetByID(c.Request().Context(), uuid)
	if err != nil {
		switch err {
		case utility.FORBIDDEN:
			return c.NoContent(http.StatusForbidden)
		case utility.NOT_FOUND:
			return c.NoContent(http.StatusNotFound)
		default:
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, %s)
}
`
const getAllTemplate string = `
func (d *Delivery) GetAll(c echo.Context) error {
	%ss, err := d.%s.GetAll(c.Request().Context())
	if err != nil {
		switch err {
		case utility.FORBIDDEN:
			return c.NoContent(http.StatusForbidden)
		case utility.NOT_FOUND:
			return c.NoContent(http.StatusNotFound)
		default:
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, %ss)
}
`
const getByReferenceTemplate string = `
func (d *Delivery) GetBy%s(c echo.Context) error {
	uuid := c.Param("uuid")
	result, err := d.%s.GetBy%s(c.Request().Context(), uuid)
	if err != nil {
		switch err {
		case utility.FORBIDDEN:
			return c.NoContent(http.StatusForbidden)
		case utility.NOT_FOUND:
			return c.NoContent(http.StatusNotFound)
		default:
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, result)
}
`
const postTemplate string = `
func (d *Delivery) Store(c echo.Context) error {
	%s := new(models.%s)
	if err := c.Bind(%s); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if err := d.%s.Store(c.Request().Context(), %s); err != nil {
		switch err {
		case utility.FORBIDDEN:
			return c.NoContent(http.StatusForbidden)
		case utility.NOT_FOUND:
			return c.NoContent(http.StatusNotFound)
		default:
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, %s)
}
`
const putTemplate string = `
func (d *Delivery) Update(c echo.Context) error {
	uuid := c.Param("uuid")
	%s := new(models.%s)
	if err := c.Bind(%s); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if err := d.%s.Update(c.Request().Context(), uuid, %s); err != nil {
		switch err {
		case utility.FORBIDDEN:
			return c.NoContent(http.StatusForbidden)
		case utility.NOT_FOUND:
			return c.NoContent(http.StatusNotFound)
		default:
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, %s)
}
`
const deleteTemplate string = `
func (d *Delivery) Delete(c echo.Context) error {
	uuid := c.Param("uuid")
	err := d.%s.Delete(c.Request().Context(), uuid)
	if err != nil {
		switch err {
		case utility.FORBIDDEN:
			return c.NoContent(http.StatusForbidden)
		case utility.NOT_FOUND:
			return c.NoContent(http.StatusNotFound)
		default:
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}
	return c.NoContent(http.StatusNoContent)
}
`
const factoryTemplate string = `
func NewDelivery(%s Service) Delivery {
	return Delivery{%s: %s}
}`

func GenerateDelivery(sourceConfig models.GeneratorSource, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("start delivery generation")
	var fileContent string
	lowercaseName := utility.ReverseTitle(sourceConfig.Name)
	serviceName := fmt.Sprintf("%sService", lowercaseName)

	// HEADER
	generateDeliveryHeader(&fileContent, sourceConfig)

	// STRUCT
	fileContent += "type Delivery struct {\n"
	fileContent += fmt.Sprintf("\t %s Service\n", serviceName)
	fileContent += "}\n\n"

	// GET
	fileContent += fmt.Sprintf(getTemplate, lowercaseName, serviceName, lowercaseName)

	// GET_ALL
	fileContent += fmt.Sprintf(getAllTemplate, lowercaseName, serviceName, lowercaseName)

	for _, a := range sourceConfig.Attributes {
		if a.IsReference {
			name := strings.Title(a.Name)
			// GET_BY_REFERENCE
			fileContent += fmt.Sprintf(getByReferenceTemplate, name, serviceName, name)
		}
	}

	// POST
	fileContent += fmt.Sprintf(postTemplate, lowercaseName, sourceConfig.Name, lowercaseName, serviceName, lowercaseName, lowercaseName)

	// PUT
	fileContent += fmt.Sprintf(putTemplate, lowercaseName, sourceConfig.Name, lowercaseName, serviceName, lowercaseName, lowercaseName)

	// DELETE
	fileContent += fmt.Sprintf(deleteTemplate, serviceName)

	// FACTORY
	fileContent += fmt.Sprintf(factoryTemplate, serviceName, serviceName, serviceName)

	log.Println("delivery generated!")
	utility.WriteFile(fmt.Sprintf("%s/%s", utility.ToSnakeCase(sourceConfig.Name), fileName), []byte(fileContent))
}

func generateDeliveryHeader(fileContent *string, sourceConfig models.GeneratorSource) {
	*fileContent += utility.GeneratePackage(utility.ToSnakeCase(sourceConfig.Name))
	*fileContent += utility.GenerateImports([]string{fmt.Sprintf("%s/models", sourceConfig.Service)},
		[]string{fmt.Sprintf("%s/utility", sourceConfig.Service)}, []string{sourceConfig.EchoVersion}, []string{"net/http"})
}
