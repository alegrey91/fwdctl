package template

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

type Generator interface {
	GetTemplateStruct() interface{}
	GetFileContent() string
	GetTemplateName() string
	GetFileName() string
}

func GenerateTemplate(g Generator, destinationPath string) error {
	tpl, err := template.New(g.GetTemplateName()).Parse(g.GetFileContent())
	if err != nil {
		return err
	}

	if destinationPath == "" {
		return nil
	}
	if !filepath.IsAbs(destinationPath){
		return fmt.Errorf("destinationPath path is not absolute")
	}

	outFile, err := os.Create(destinationPath + "/" + g.GetFileName())
	if err != nil {
		return err
	}

	err = tpl.Execute(outFile, g.GetTemplateStruct())
	if err != nil {
		return err
	}
	return nil
}
