package command

import (
	"testing"

	"github.com/SpencerBrown/go-http/flag"
)

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
func TestName(t *testing.T) {
	tests := []struct {
		name    string
		command *Command
		want    string
	}{
		{
			name:    "validName",
			command: NewCommand("test", nil, "test command", flag.NewFlags()),
			want:    "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.command.Name() != tt.want {
				t.Errorf("Name() = %s, want %s", tt.command.Name(), tt.want)
			}
		})
	}
}

func TestAlias(t *testing.T) {
	tests := []struct {
		name    string
		command *Command
		want    []string
	}{
		{
			name:    "validAlias",
			command: NewCommand("test", []string{"t"}, "test command", flag.NewFlags()),
			want:    []string{"t"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !equalStringSlices(tt.command.Alias(), tt.want) {
				t.Errorf("Alias() = %v, want %v", tt.command.Alias(), tt.want)
			}
		})
	}
}

func TestDescription(t *testing.T) {
	tests := []struct {
		name    string
		command *Command
		want    string
	}{
		{
			name:    "validDescription",
			command: NewCommand("test", nil, "test command", flag.NewFlags()),
			want:    "test command",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.command.Description() != tt.want {
				t.Errorf("Description() = %s, want %s", tt.command.Description(), tt.want)
			}
		})
	}
}

func TestFlags(t *testing.T) {
	tests := []struct {
		name    string
		command *Command
		want    flag.Flags
	}{
		{
			name:    "validFlags",
			command: NewCommand("test", nil, "test command", flag.NewFlags()),
			want:    flag.NewFlags(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !equalFlags(tt.command.Flags(), tt.want) {
				t.Errorf("Flags() = %v, want %v", tt.command.Flags(), tt.want)
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		command  *Command
		want1    string
	}{
		{
			name:    "validCommand",
			command: NewCommand("test", []string{"t"}, "test command", flag.NewFlags()),
			want1: `Command: test
  Aliases: t
  Description: test command
  Flags:
  Name Short Aliases Default Type Description
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.command.String() != tt.want1 {
				t.Errorf("command.String() = \n%s, want1 \n%s", tt.command.String(), tt.want1)
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
	SetSubTester3(t)
}

// SetSubTester1 tests the SetSub method for a valid subcommand
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

// SetSubTester2 tests the SetSub method for a duplicate subcommand
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

// SetSubTester3 tests the SetSub method for duplicate flags between subcommands
func SetSubTester3(t *testing.T) {
	rootCmd := NewCommand("root", []string{"r"}, "root command", flag.NewFlags())
	subCmd1 := NewCommand("sub1", []string{"s1"}, "sub1 command", flag.NewFlags())
	subCmd2 := NewCommand("sub2", []string{"s2"}, "sub2 command", flag.NewFlags())

	rootCmd.flags.AddFlag(flag.NewFlag("aflag", nil, "f", "A Flag", 42))
	subCmd1.flags.AddFlag(flag.NewFlag("aflag", nil, "g", "Another Flag", 43))

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("SetSub() did not panic on duplicate flags between subcommands")
		}
	}()
	rootCmd.SetSub(subCmd1)
	rootCmd.SetSub(subCmd2)
	subCmd2.SetSub(subCmd1)
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
