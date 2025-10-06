package option

import (
	"testing"
)

func TestNewOptions(t *testing.T) {
	opts := NewOptions()
	if opts == nil {
		t.Error("NewOptions() returned nil")
	}
	if len(opts) != 0 {
		t.Errorf("NewOptions() length = %d, want 0", len(opts))
	}
}

func TestNewOption(t *testing.T) {
	tests := []struct {
		name        string
		nm          string
		al          []string
		sn          rune
		sa          []rune
		desc        string
		longdesc    string
		value       interface{}
		expectError bool
	}{
		{
			name:     "validStringOption",
			nm:       "test",
			al:       []string{"alias"},
			sn:       'x',
			sa:       []rune{'y'},
			desc:     "test option",
			longdesc: "long test option",
			value:    "default",
		},
		{
			name:     "validWithNils",
			nm:       "test",
			al:       nil,
			sn:       'x',
			sa:       nil,
			desc:     "test option",
			longdesc: "long test option",
			value:    "default",
		},
		{
			name:        "blankName",
			nm:          "    ",
			al:          []string{"alias"},
			sn:          'x',
			sa:          []rune{'y'},
			desc:        "test option",
			longdesc:    "long test option",
			value:       "default",
			expectError: true,
		},
		{
			name:        "singleCharAlias",
			nm:          "abc",
			al:          []string{"b"},
			sn:          'x',
			sa:          []rune{'y'},
			desc:        "",
			longdesc:    "",
			value:       "default",
			expectError: true,
		},
		{
			name:        "blankAlias",
			nm:          "abc",
			al:          []string{"\t\n "},
			sn:          'x',
			sa:          []rune{'y'},
			desc:        "",
			longdesc:    "",
			value:       "default",
			expectError: true,
		},
		{
			name:        "singleCharName",
			nm:          "a",
			al:          []string{"alias"},
			sn:          'x',
			sa:          []rune{'y'},
			desc:        "test option",
			longdesc:    "long test option",
			value:       "default",
			expectError: true,
		}, {
			name:        "singleRuneName",
			nm:          "✓", // check mark U+2713
			al:          []string{"alias"},
			sn:          'x',
			sa:          []rune{'y'},
			desc:        "test option",
			longdesc:    "long test option",
			value:       "default",
			expectError: true,
		},
		{
			name:        "duplicateAlias",
			nm:          "test",
			al:          []string{"test"},
			sn:          'x',
			sa:          []rune{'y'},
			desc:        "test option",
			longdesc:    "long test option",
			value:       "default",
			expectError: true,
		},
		{
			name:        "duplicateAlias2",
			nm:          "test",
			al:          []string{"test2", "test2"},
			sn:          0,
			sa:          []rune{},
			desc:        "test option",
			longdesc:    "long test option",
			value:       "default",
			expectError: true,
		},
		{
			name:        "ShortNameWhitespace",
			nm:          "test",
			al:          []string{"alias"},
			sn:          ' ',
			sa:          nil,
			desc:        "test option",
			longdesc:    "long test option",
			value:       "default",
			expectError: true,
		},
		{
			name:        "shortNameWithoutShortName",
			nm:          "test",
			al:          []string{"alias"},
			sn:          0,
			sa:          []rune{'y'},
			desc:        "test option",
			longdesc:    "long test option",
			value:       "default",
			expectError: true,
		},
		{
			name:        "ShortAliasZero",
			nm:          "test",
			al:          []string{"alias"},
			sn:          'x',
			sa:          []rune{'y', 0},
			desc:        "test option",
			longdesc:    "long test option",
			value:       "default",
			expectError: true,
		},
		{
			name:        "ShortAliasWhitespace",
			nm:          "test",
			al:          []string{"alias"},
			sn:          'x',
			sa:          []rune{'y', '\t'},
			desc:        "test option",
			longdesc:    "long test option",
			value:       "default",
			expectError: true,
		},
		{
			name:        "ShortAliasDuplicate",
			nm:          "test",
			al:          []string{"alias"},
			sn:          'x',
			sa:          []rune{'y', 'z', 'y'},
			desc:        "test option",
			longdesc:    "long test option",
			value:       "default",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt, err := NewOption(tt.nm, tt.al, tt.sn, tt.sa, tt.desc, tt.longdesc, tt.value.(string), nil)
			if tt.expectError {
				if err == nil {
					t.Errorf("NewOption() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("NewOption() unexpected error: %v", err)
				}
				if opt == nil {
					t.Error("NewOption() returned nil option")
				} else {
					if opt.Name() != "test" {
						t.Errorf("NewOption() name = %s, want test", opt.Name())
					}
					if opt.ShortName() != tt.sn {
						t.Errorf("NewOption() shortName = %c, want %c", opt.ShortName(), tt.sn)
					}
					if val, _ := GetValue[string](opt); val != tt.value.(string) {
						t.Errorf("NewOption() value = %v, want %v", val, tt.value)
					}
				}
			}
		})
	}
}

func TestGetValue(t *testing.T) {
	opt, _ := NewOption("test", nil, 't', nil, "test", "test", "value", nil)
	v, err := GetValue[string](opt)
	if err != nil {
		t.Errorf("GetValue() error: %v", err)
	}
	if v != "value" {
		t.Errorf("GetValue[string]() = %s, want value", v)
	}
}

func TestAddOption(t *testing.T) {
	opts := NewOptions()
	opt1, _ := NewOption("test1", nil, 'a', nil, "test1", "test1", "value1", nil)
	opt2, _ := NewOption("test2", nil, 'b', nil, "test2", "test2", "value2", nil)
	err := opts.AddOption(opt1)
	if err != nil {
		t.Errorf("AddOption() error: %v", err)
	}
	err = opts.AddOption(opt2)
	if err != nil {
		t.Errorf("AddOption() error: %v", err)
	}
	if len(opts) != 2 {
		t.Errorf("AddOption() length = %d, want 2", len(opts))
	}
}

func TestFindOption(t *testing.T) {
	opts := NewOptions()
	opt, _ := NewOption("test", []string{"alias"}, 'x', nil, "test", "test", "value", nil)
	err := opts.AddOption(opt)
	if err != nil {
		t.Errorf("AddOption() error: %v", err)
	}
	found := opts.FindOption("test")
	if found == nil {
		t.Error("FindOption() returned nil for existing option")
	}
	foundAlias := opts.FindOption("alias")
	if foundAlias == nil {
		t.Error("FindOption() returned nil for alias")
	}
	notFound := opts.FindOption("nonexistent")
	if notFound != nil {
		t.Error("FindOption() returned non-nil for nonexistent option")
	}
}

func TestParseValue(t *testing.T) {
	optInt, _ := NewOption("intopt", nil, 0, nil, "int", "int", 0, nil)
	err := optInt.ParseValue("42")
	if err != nil {
		t.Errorf("ParseValue() error = %v", err)
	}
	if val, _ := GetValue[int](optInt); val != 42 {
		t.Errorf("ParseValue() int value = %d, want 42", val)
	}

	optBool, _ := NewOption("boolopt", nil, 0, nil, "bool", "bool", false, nil)
	err = optBool.ParseValue("true")
	if err != nil {
		t.Errorf("ParseValue() error = %v", err)
	}
	if val, _ := GetValue[bool](optBool); !val {
		t.Error("ParseValue() bool value = false, want true")
	}
}
func TestGetOptionOK(t *testing.T) {
	opts := NewOptions()
	opt, _ := NewOption("test", []string{"alias"}, 'x', nil, "desc", "longdesc", "value", nil)
	err := opts.AddOption(opt)
	if err != nil {
		t.Fatalf("AddOption() error: %v", err)
	}

	gotOpt, ok := GetOptionOK(opts, "test")
	if !ok {
		t.Error("GetOptionOK() did not find existing option by name")
	}
	if gotOpt != opt {
		t.Error("GetOptionOK() returned wrong option for name")
	}

	gotOptAlias, okAlias := GetOptionOK(opts, "alias")
	if okAlias {
		t.Error("GetOptionOK() should not find option by alias")
	}
	if gotOptAlias != nil {
		t.Error("GetOptionOK() should return nil for alias")
	}

	gotOptMissing, okMissing := GetOptionOK(opts, "missing")
	if okMissing {
		t.Error("GetOptionOK() should not find missing option")
	}
	if gotOptMissing != nil {
		t.Error("GetOptionOK() should return nil for missing option")
	}
}

func TestGetOption(t *testing.T) {
	opts := NewOptions()
	opt, _ := NewOption("test", nil, 'x', nil, "desc", "longdesc", "value", nil)
	_ = opts.AddOption(opt)

	gotOpt, err := GetOption(opts, "test")
	if err != nil {
		t.Errorf("GetOption() error: %v", err)
	}
	if gotOpt != opt {
		t.Error("GetOption() returned wrong option")
	}

	_, err = GetOption(opts, "missing")
	if err == nil {
		t.Error("GetOption() should return error for missing option")
	}
}

func TestGetValueOK(t *testing.T) {
	optStr, _ := NewOption("str", nil, 0, nil, "", "", "abc", nil)
	valStr, okStr := GetValueOK[string](optStr)
	if !okStr || valStr != "abc" {
		t.Errorf("GetValueOK[string] = %v, %v; want abc, true", valStr, okStr)
	}

	optInt, _ := NewOption("int", nil, 0, nil, "", "", 42, nil)
	valInt, okInt := GetValueOK[int](optInt)
	if !okInt || valInt != 42 {
		t.Errorf("GetValueOK[int] = %v, %v; want 42, true", valInt, okInt)
	}

	valWrong, okWrong := GetValueOK[bool](optStr)
	if okWrong {
		t.Error("GetValueOK[bool] should be false for string value")
	}
	_ = valWrong // just to avoid unused variable warning
}

func TestGetValueAny(t *testing.T) {
	optStr, _ := NewOption("str", nil, 0, nil, "", "", "abc", nil)
	val := optStr.GetValueAny()
	if val != "abc" {
		t.Errorf("GetValueAny() = %v, want abc", val)
	}

	optInt, _ := NewOption("int", nil, 0, nil, "", "", 42, nil)
	valInt := optInt.GetValueAny()
	if valInt != 42 {
		t.Errorf("GetValueAny() = %v, want 42", valInt)
	}
}