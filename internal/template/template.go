package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Generator interface {
	GetTemplateStruct() interface{}
	GetFileContent() string
	GetTemplateName() string
	GetFileName() string
}

func GenerateTemplate(g Generator, outputPath string) error {
	tpl, err := template.New(g.GetTemplateName()).Parse(g.GetFileContent())
	if err != nil {
		return fmt.Errorf("error getting template instance: %v", err)
	}

	if !filepath.IsAbs(outputPath) {
		return fmt.Errorf("output path is not absolute: %s", outputPath)
	}
	// if last char of outputPath is "/" we want to remove,
	// so the final output will be cleaned.
	// this way: /root/template.file instead of /root//template.file
	if outputPath != "/" && outputPath[len(outputPath)-1:] == "/" {
		outputPath = strings.TrimSuffix(outputPath, "/")
	}

	outFile, err := os.Create(filepath.Join(outputPath, g.GetFileName()))
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}

	err = tpl.Execute(outFile, g.GetTemplateStruct())
	if err != nil {
		return fmt.Errorf("error writing content into file: %v", err)
	}
	return nil
}
