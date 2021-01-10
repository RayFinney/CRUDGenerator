package utility

import (
	"fmt"
)

func GeneratePackage(name string) string {
	return fmt.Sprintf("package %s \n\n", ReverseTitle(name))
}

func GenerateImports(imports ...[]string) string {
	importContent := "import (\n"
	for _, singleImport := range imports {
		if len(singleImport) == 1 {
			importContent += fmt.Sprintf("\t\"%s\"\n", singleImport[0])
		} else if len(singleImport) >= 2 {
			importContent += fmt.Sprintf("\t%s \"%s\"\n", singleImport[0], singleImport[1])
		}
	}
	importContent += ")\n\n"
	return importContent
}
