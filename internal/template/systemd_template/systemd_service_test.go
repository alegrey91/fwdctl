package systemd_template

import (
	_ "embed"
	"os"
	"reflect"
	"testing"
)

func Test_serviceTypeAllowed(t *testing.T) {
	type args struct {
		st string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "service_type_not_allowed",
			args: args{
				st: "none",
			},
			want: false,
		},
		{
			name: "service_type_allowed",
			args: args{
				st: "fork",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := serviceTypeAllowed(tt.args.st); got != tt.want {
				t.Errorf("serviceTypeAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSystemdService(t *testing.T) {
	type args struct {
		serviceType      string
		installationPath string
		rulesFile        string
	}
	tests := []struct {
		name    string
		args    args
		want    *SystemdService
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "should_return_valid_struct",
			args: args{
				serviceType:      "fork",
				installationPath: "/opt/",
				rulesFile:        "/tmp/rules.yml",
			},
			want: &SystemdService{
				ServiceType:      "fork",
				InstallationPath: "/opt/",
				RulesFile:        "/tmp/rules.yml",
			},
			wantErr: false,
		},
		{
			name: "should_return_nil_due_to_wrong_service_type",
			args: args{
				serviceType:      "aaa",
				installationPath: "/opt/",
				rulesFile:        "/home/user/rules.yml",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should_return_nil_due_to_wrong_installation_path",
			args: args{
				serviceType:      "oneshot",
				installationPath: "../../tmp/",
				rulesFile:        "/home/user/rules.yml",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should_return_nil_due_to_wrong_installation_path_2",
			args: args{
				serviceType:      "oneshot",
				installationPath: "/this/path/doesnt/exist",
				rulesFile:        "/home/user/rules.yml",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should_return_nil_due_to_wrong_rules_file",
			args: args{
				serviceType:      "oneshot",
				installationPath: "/opt/",
				rulesFile:        "rules.yml",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should_return_nil_due_to_wrong_rules_file_2",
			args: args{
				serviceType:      "oneshot",
				installationPath: "/opt/",
				rulesFile:        "/this/path/doesnt/exist/rules.yml",
			},
			want:    nil,
			wantErr: true,
		},
	}
	fileName := "/tmp/rules.yml"
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		file, err := os.Create(fileName)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSystemdService(tt.args.serviceType, tt.args.installationPath, tt.args.rulesFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSystemdService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSystemdService() = %v, want %v", got, tt.want)
			}
		})
	}
}
