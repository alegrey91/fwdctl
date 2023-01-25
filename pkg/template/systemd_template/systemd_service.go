package systemd_template

import (
	_ "embed"
	"fmt"
	"path/filepath"
)

//go:embed fwdctl.service.tpl
var systemdTemplate string
var systemdTemplateName = "systemd"
var systemdFileName = "fwdctl.service"
var allowedServiceTypes = [2]string{"oneshot", "fork"}

type SystemdService struct {
	ServiceType      string
	InstallationPath string
	RulesFile        string
}

func serviceTypeAllowed(st string) bool {
	for _, ast := range allowedServiceTypes {
		if ast == st {
			return true
		}
	}
	return false
}

func NewSystemdService(serviceType, installationPath, rulesFile string) (*SystemdService, error) {
	if !serviceTypeAllowed(serviceType) {
		return nil, fmt.Errorf("service type provided not allowed")
	}
	if !filepath.IsAbs(installationPath) {
		// TODO: substitute with fmt.Errorf()
		return nil, fmt.Errorf("installation path is not absolute")
	}
	if !filepath.IsAbs(rulesFile) {
		return nil, fmt.Errorf("rules file is not absolute")
	}

	return &SystemdService{
		ServiceType:      serviceType,
		InstallationPath: installationPath,
		RulesFile:        rulesFile,
	}, nil
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
