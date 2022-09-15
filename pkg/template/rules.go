package template

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed tpl/rules.yml.tpl
var rulesTpl string

type Rule struct {
}

func (rule *Rule) GenerateTemplate(output string) error {
	tpl, err := template.New("rules").Parse(rulesTpl)
	if err != nil {
		return err
	}

	if output == "" {
		return nil
	}
	if !filepath.IsAbs(output) {
		return fmt.Errorf("output path is not absolute")
	}

	outFile, err := os.Create(output)
	if err != nil {
		return err
	}

	err = tpl.Execute(outFile, rule)
	if err != nil {
		return err
	}
	return nil
}
