package template

import (
	_ "embed"
	"os"
	"text/template"
)

//go:embed tpl/fwdctl.service.tpl
var systemdServiceTpl string

type SystemdService struct {
    InstallationPath string
    RulesFile string
}

func GetSystemdTemplate(toOutput string, installationPath string, rulesFile string) error {
	systemdSvc := SystemdService{
		InstallationPath: installationPath,
		RulesFile: rulesFile,
	}

	tpl, err := template.New("systemd").Parse(systemdServiceTpl)
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

	err = tpl.Execute(outFile, systemdSvc)
	if err != nil {
		return err
	}
	return nil
}
