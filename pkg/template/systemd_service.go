package template

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed tpl/fwdctl.service.tpl
var systemdServiceTpl string

type SystemdService struct {
	InstallationPath string
	RulesFile        string
}

func (systemdSvc *SystemdService) GenerateTemplate(toOutput string) error {
	tpl, err := template.New("systemd").Parse(systemdServiceTpl)
	if err != nil {
		return err
	}

	if toOutput == "" {
		return nil
	}
	if !filepath.IsAbs(toOutput) {
		return fmt.Errorf("output path is not absolute")
	}

	outFile, err := os.Create(toOutput)
	if err != nil {
		return err
	}

	err = tpl.Execute(outFile, systemdSvc)
	if err != nil {
		return err
	}
	return nil
}
