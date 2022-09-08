package template

import (
	_ "embed"
	"os"
	"text/template"
)

//go:embed tpl/rules.yml.tpl
var rulesTpl string

func GetRuleTemplate(toOutput string) error {
	tpl, err := template.New("rules").Parse(rulesTpl)
	if err != nil {
		return err
	}

	if toOutput == "" {
		return nil
	}

	outFile, err := os.Create(toOutput)
	if err != nil {
		return err
	}

	err = tpl.Execute(outFile, nil)
	if err != nil {
		return err
	}
	return nil
}
