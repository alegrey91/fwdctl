package template

import "testing"

type Test struct {
}

func NewTest() *Test {
	return &Test{}
}

func (t *Test) GetTemplateStruct() interface{} {
	return t
}

func (t *Test) GetFileContent() string {
	return "this_is_the_content"
}

func (t *Test) GetTemplateName() string {
	return "this_is_the_name"
}

func (t *Test) GetFileName() string {
	return "this_is_the_file_name"
}

func TestGenerateTemplate(t *testing.T) {
	type args struct {
		g          Generator
		outputPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "this_should_pass",
			args: args{
				g:          NewTest(),
				outputPath: "/tmp/",
			},
			wantErr: false,
		},
		{
			name: "this_should_fail_due_to_wrong_output_path",
			args: args{
				g:          NewTest(),
				outputPath: "tmp/",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GenerateTemplate(tt.args.g, tt.args.outputPath); (err != nil) != tt.wantErr {
				t.Errorf("GenerateTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
