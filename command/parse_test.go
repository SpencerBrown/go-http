package command

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	cmds := &Commands{
		args: []string{},
		root: &Command{
			name:  "root",
			alias: []string{"r"},
		},
	}

	tests := []struct {
		name     string
		commands *Commands
		args     []string
		want     *Commands
		wantErr  bool
	}{
		{"NilCommands", nil, nil, nil, true},
		{"NilRootCommand", &Commands{}, nil, nil, true},
		{"NilArgs", cmds, nil, nil, true},
		{"EmptyArgs", cmds, []string{}, nil, true},
		{"BlankArg", cmds, []string{"\n"}, nil, true},
		{"UnrecognizedCommand", cmds, []string{"unknown"}, nil, true},
		{"RecognizedCommand", cmds, []string{"root"},
			&Commands{
				args: []string{},
				root: &Command{
					name:  "root",
					alias: []string{"r"},
				},
			}, false,
		},
		{"RecognizedAlias", cmds, []string{"r"},
			&Commands{
				args: []string{},
				root: &Command{
					name:  "root",
					alias: []string{"r"},
				},
			}, false,
		},
		{"EndOfFlags", cmds, []string{"root", "--", "arg1", "arg2"},
			&Commands{
				args: []string{"arg1", "arg2"},
				root: &Command{
					name:  "root",
					alias: []string{"r"},
				},
			}, false,
		},
		{"StartOfArgs", cmds, []string{"root", "arg1", "arg2"},
			&Commands{
				args: []string{"arg1", "arg2"},
				root: &Command{
					name:  "root",
					alias: []string{"r"},
				},
			}, false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Parse(tt.commands, tt.args)
			if err == nil {
				if tt.wantErr      {
					t.Errorf("Parse() error = nil, wantErr %v", tt.wantErr)
				} else {
					if !reflect.DeepEqual(tt.commands, tt.want) {
						t.Errorf("Parse() got = %v, want %v", tt.commands, tt.want)
					}
				}
			} else {
				if !tt.wantErr {
					t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}
