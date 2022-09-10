package rules

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// NewRulesFile return the struct that contains informations about rules
func NewRulesFile(path string) (*RulesFile, error) {
	rulesFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	rules := RulesFile{}
	err = yaml.Unmarshal(rulesFile, &rules)
	if err != nil {
		return nil, err
	}
	return &rules, nil
}
