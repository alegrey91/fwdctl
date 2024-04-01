package systemd_template

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
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
	// checks for systemd service type
	if !serviceTypeAllowed(serviceType) {
		return nil, fmt.Errorf("service type is not allowed: %s", serviceType)
	}
	// checks for installation path
	if !filepath.IsAbs(installationPath) {
		return nil, fmt.Errorf("installation path is not absolute: %s", installationPath)
	}
	if _, err := os.Stat(installationPath); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("installation path does not exist: %s", installationPath)
	}
	// checks for rules file
	if !filepath.IsAbs(rulesFile) {
		return nil, fmt.Errorf("rules file path is not absolute: %s", rulesFile)
	}
	if _, err := os.Stat(installationPath); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("rules file path does not exist: %s", rulesFile)
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
