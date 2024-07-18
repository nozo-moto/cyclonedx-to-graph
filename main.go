package main

import (
	"crypto/md5"
	"fmt"
	"os"

	"github.com/CycloneDX/cyclonedx-go"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <cyclonedx_file>")
		return
	}
	filePath := os.Args[1]

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("ファイルの読み込み中にエラーが発生しました: %s\n", err)
		return
	}
	defer file.Close()

	// CycloneDX構造体を初期化
	var bom cyclonedx.BOM
	err = cyclonedx.NewBOMDecoder(file, cyclonedx.BOMFileFormatJSON).Decode(&bom)
	if err != nil {
		fmt.Printf("JSONのパース中にエラーが発生しました: %s\n", err)
		return
	}

	fmt.Println(exportGraphiz(bom))
	os.Exit(0)
}

func exportMermaid(bom cyclonedx.BOM) string {
	result := "graph TD;\n"
	for _, component := range *bom.Components {
		refId := genId(component.BOMRef)
		for _, dep := range getDependencies(bom, component.BOMRef) {
			depId := genId(dep)
			result += fmt.Sprintf(
				"\t%s(%s: %s) --> %s(%s: %s);\n",
				refId,
				component.Name,
				component.Version,
				depId,
				getComponent(bom, dep).Name,
				getComponent(bom, dep).Version,
			)
		}
	}

	// return fmt.Sprintf("```mermaid\n%s```", result)
	return result
}

func exportGraphiz(bom cyclonedx.BOM) string {
	result := "digraph {\n"
	for _, component := range *bom.Components {
		refId := genId(component.BOMRef)
		result += fmt.Sprintf("%s [\n\tlabel = \"%s:%s\"\n];\n", refId, component.Name, component.Version)
		for _, dep := range getDependencies(bom, component.BOMRef) {
			result += fmt.Sprintf(
				"\t%s -> %s;\n",
				refId,
				genId(dep),
			)
		}
	}
	result += "}\n"
	return result
}

func getComponent(bom cyclonedx.BOM, ref string) *cyclonedx.Component {
	for _, component := range *bom.Components {
		if component.BOMRef == ref {
			return &component
		}
	}
	return nil
}

func getDependencies(bom cyclonedx.BOM, root string) []string {
	dependencies := []string{}
	for _, ref := range *bom.Dependencies {
		if ref.Ref == root {
			dependencies = append(dependencies, *ref.Dependencies...)
		}
	}
	return dependencies
}

func genId(ref string) string {
	return fmt.Sprintf("md5_%x", md5.Sum([]byte(ref)))
}
