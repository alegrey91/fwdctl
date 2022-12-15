package rules_template

import (
	_ "embed"
)

//go:embed rules.yml.tpl
var rulesTemplate string
var rulesTemplateName = "rules"
var rulesFileName = "rules.yml"

type Rule struct {
}

func NewRules() *Rule {
	return &Rule{}
}

func (r *Rule) GetTemplateStruct() interface{} {
	return r
}

func (r *Rule) GetFileContent() string {
	return rulesTemplate
}

func (r *Rule) GetTemplateName() string {
	return rulesTemplateName
}

func (r *Rule) GetFileName() string {
	return rulesFileName
}
