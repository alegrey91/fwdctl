package systemd_template

import (
	_ "embed"
	"path/filepath"
)

//go:embed fwdctl.service.tpl
var systemdTemplate string
var systemdTemplateName = "systemd"
var systemdFileName = "fwdctl.service"

type SystemdService struct {
	InstallationPath string
	RulesFile        string
}

func NewSystemdService(installationPath, rulesFile string) *SystemdService {
	if !filepath.IsAbs(installationPath) {
		panic("installation path is not absolute")
	}

	if !filepath.IsAbs(rulesFile) {
		panic("rules file is not absolute")
	}

	return &SystemdService{
		InstallationPath: installationPath,
		RulesFile:        rulesFile,
	}
}

func (s *SystemdService) GetTemplateStruct() interface{} {
	return s
}

func (s *SystemdService) GetFileContent() string {
	return systemdTemplate
}

func (s *SystemdService) GetTemplateName() string {
	return systemdTemplateName
}

func (s *SystemdService) GetFileName() string {
	return systemdFileName
}
