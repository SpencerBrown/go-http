package command

import (
	"testing"

	"github.com/SpencerBrown/go-http/flag"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		cmds      []string
		flgs      []string
		args      []string
		wantNames []string
		wantArgs  []string
		wantFlags flag.Flags
		wantErr   bool
	}{
		{
			name:    "no cmds",
			cmds:    []string{},
			wantErr: true,
		},
		{
			name:      "no flags or args",
			cmds:      []string{"root"},
			wantNames: []string{"root"},
			wantArgs:  []string{},
			wantErr:   false,
		},
		{
			name:      "single command, flags and args with equal sign",
			cmds:      []string{"root"},
			flgs:      []string{"flag1", "value1", "flag2", "value2"},
			args:      []string{"--flag1=value1", "--flag2=value2", "arg1", "arg2"},
			wantNames: []string{"root"},
			wantArgs:  []string{"arg1", "arg2"},
			wantFlags: flag.NewFlags().AddFlag(flag.NewFlag("flag1", nil, "", "", "value1")).AddFlag(flag.NewFlag("flag2", nil, "", "", "value2")),
			wantErr:   false,
		},
		{
			name:      "subcommand, flags and args",
			cmds:      []string{"root", "sub"},
			args:      []string{"sub", "--flag2=value2", "arg1", "arg2"},
			wantNames: []string{"root", "sub"},
			wantArgs:  []string{"arg1", "arg2"},
			wantFlags: flag.NewFlags().AddFlag(flag.NewFlag("flag2", nil, "", "", "value2")),
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.cmds)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if !equalSlices(got.names, tt.wantNames) {
					t.Errorf("Parse() names = %v, want %v", got.names, tt.wantNames)
				}
				if !equalSlices(got.args, tt.wantArgs) {
					t.Errorf("Parse() args = %v, want %v", got.args, tt.wantArgs)
				}
			}
		})
	}
}

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
