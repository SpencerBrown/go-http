package command

import (
	"testing"

	"github.com/SpencerBrown/go-http/flag"
)

func TestNewCommands(t *testing.T) {
	cmds := NewCommands()
	if cmds == nil {
		t.Fatal("NewCommands() returned nil")
	}
	if cmds.root != nil {
		t.Errorf("NewCommands() root = %v, want nil", cmds.root)
	}
	if len(cmds.args) != 0 {
		t.Errorf("NewCommands() args length = %d, want 0", len(cmds.args))
	}
}

func TestNewCommand(t *testing.T) {
	tests := []struct {
		name        string
		nm          string
		al          []string
		desc        string
		flgs        flag.Flags
		shouldPanic bool
	}{
		{
			name: "validCommand",
			nm:   "test",
			al:   []string{"t"},
			desc: "test command",
			flgs: flag.NewFlags(),
		},
		{
			name:        "blankName",
			nm:          " ",
			al:          []string{"t"},
			desc:        "test command",
			flgs:        flag.NewFlags(),
			shouldPanic: true,
		},
		{
			name:        "blankAlias",
			nm:          "test",
			al:          []string{" "},
			desc:        "test command",
			flgs:        flag.NewFlags(),
			shouldPanic: true,
		},
		{
			name:        "duplicateNameAlias",
			nm:          "test",
			al:          []string{"test"},
			desc:        "test command",
			flgs:        flag.NewFlags(),
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("NewCommand() did not panic for test %s", tt.name)
					} else {
						t.Logf("Recovered from panic for test %s: %s", tt.name, r)
					}
				}()
			}
			// The following line should panic if shouldPanic is true
			// in which case the deferred function will run and recover the panic
			// and log the recovery (only printedin verbose mode)
			// then the test will continue to the next iteration of the loop
			// because the t.Run() function exits after the deferred function runs.
			cmd := NewCommand(tt.nm, tt.al, tt.desc, tt.flgs)
			if tt.shouldPanic {
				t.Errorf("NewCommand() did not panic for test %s", tt.name)
			} else {
				if cmd.Name() != tt.nm {
					t.Errorf("NewCommand() name = %s, want %s", cmd.Name(), tt.nm)
				}
				if !equalStringSlices(cmd.Alias(), tt.al) {
					t.Errorf("NewCommand() alias = %v, want %v", cmd.Alias(), tt.al)
				}
				if cmd.Description() != tt.desc {
					t.Errorf("NewCommand() description = %s, want %s", cmd.Description(), tt.desc)
				}
				if !equalFlags(cmd.Flags(), tt.flgs) {
					t.Errorf("NewCommand() flags = %v, want %v", cmd.Flags(), tt.flgs)
				}
			}
		})
	}
}

func TestSetRoot(t *testing.T) {
	tests := []struct {
		name        string
		commands    *Commands
		command     *Command
		want        *Commands
		shouldPanic bool
	}{
		{
			name:        "nilCmds",
			commands:    nil,
			command:     NewCommand("root", []string{"r"}, "root command", flag.NewFlags()),
			shouldPanic: true,
		},
		{
			name:        "nilCmd",
			commands:    NewCommands(),
			command:     nil,
			shouldPanic: true,
		},
		{
			name: "doubleRoot",
			commands: &Commands{
				args: make([]string, 0),
				root: NewCommand("root", nil, "root command", flag.NewFlags()),
			},
			command:     NewCommand("root", nil, "root command", flag.NewFlags()),
			shouldPanic: true,
		},
		{
			name:     "validRoot",
			commands: NewCommands(),
			command:  NewCommand("root", []string{"r"}, "root command", flag.NewFlags()),
			want: &Commands{
				args: []string{"r"},
				root: NewCommand("root", []string{"r"}, "root command", flag.NewFlags()),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("SetRoot() did not panic in test %s when it should have", tt.name)
					}
				}()
			}
			tt.commands.SetRoot(tt.command)
			if tt.shouldPanic {
				t.Errorf("SetRoot() did not panic for test %s when it should have", tt.name)
			}
			if !equalCommand(tt.commands.root, tt.want.root) {
				t.Errorf("SetRoot() args = %v, want %v", tt.commands.args, tt.want.args)
			}
		})
	}
}
func TestSetSub(t *testing.T) {
	tests := []struct {
		name        string
		parent      *Command
		command     *Command
		want        *Command
		shouldPanic bool
	}{
		{
			name:        "nilParent",
			parent:      nil,
			command:     NewCommand("sub", nil, "subcommand", flag.NewFlags()),
			shouldPanic: true,
		},
		{
			name:        "nilCmd",
			parent:      NewCommand("root", nil, "root command", flag.NewFlags()),
			command:     nil,
			shouldPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("SetSub() did not panic in test %s when it should have", tt.name)
					}
				}()
			}
			tt.parent.SetSub(tt.command)
			if tt.shouldPanic {
				t.Errorf("SetSub() did not panic for test %s when it should have", tt.name)
			}
		})
	}
	SetSubTester1(t)
	SetSubTester2(t)
}

func equalCommand(a, b *Command) bool {
	return a.Name() == b.Name() && equalStringSlices(a.Alias(), b.Alias()) && a.Description() == b.Description() && equalFlags(a.Flags(), b.Flags())
}

func SetSubTester1(t *testing.T) {
	rootCmd := NewCommand("root", []string{"r"}, "root command", flag.NewFlags())
	subCmd := NewCommand("sub", []string{"s"}, "sub command", flag.NewFlags())

	rootCmd.SetSub(subCmd)
	if len(rootCmd.sub) != 1 {
		t.Errorf("SetSub() sub length = %d, want 1", len(rootCmd.sub))
	}
	if rootCmd.sub[0] != subCmd {
		t.Errorf("SetSub() sub[0] = %v, want %v", rootCmd.sub[0], subCmd)
	}
	if subCmd.parent != rootCmd {
		t.Errorf("SetSub() parent = %v, want %v", subCmd.parent, rootCmd)
	}

	duplicateCmd := NewCommand("sub", []string{"s"}, "duplicate sub command", flag.NewFlags())
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("SetSub() did not panic on duplicate subcommand")
		}
	}()
	rootCmd.SetSub(duplicateCmd)
}
func SetSubTester2(t *testing.T) {
	rootCmd := NewCommand("root", []string{"r"}, "root command", flag.NewFlags())
	subCmd := NewCommand("sub", []string{"s"}, "sub command", flag.NewFlags())

	rootCmd.SetSub(subCmd)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("SetSub() did not panic on duplicate subcommand")
		}
	}()
	rootCmd.SetSub(subCmd)
}

func equalStringSlices(a, b []string) bool {
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

func equalFlags(a, b flag.Flags) bool {
	if len(a) != len(b) {
		return false
	}
	for key, valA := range a {
		valB, ok := b[key]
		if !ok || !equalFlag(valA, valB) {
			return false
		}
	}
	return true
}

func equalFlag(a, b *flag.Flag) bool {
	return a.Name() == b.Name() && a.ShortName() == b.ShortName() && a.Description() == b.Description() && equalStringSlices(a.Alias(), b.Alias()) && a.GetValueAny() == b.GetValueAny()
}
