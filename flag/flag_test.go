package flag

import (
	"reflect"
	"strconv"
	"testing"
)

func TestNewFlags(t *testing.T) {
	t.Run("NewFlags", func(t *testing.T) {
		want := make(Flags)
		if got := NewFlags(); !reflect.DeepEqual(got, want) {
			t.Errorf("NewFlags() = %v, want %v", got, want)
		}
	})
}

// TestNewFlag also tests Flag.Name(), Flag.Alias(), Flag.ShortName(), Flag.Description()
func TestNewFlag(t *testing.T) {
	type args struct {
		name        string
		alias       []string
		shortName   string
		description string
		value       any
	}
	tests := []struct {
		name        string
		args        args
		want        Flag
		shouldPanic bool
	}{
		{
			name: "newIntFlag",
			args: args{
				name:        "IntFlag",
				alias:       []string{"IF", "IntFlg"},
				shortName:   "I",
				description: "Int Flag",
				value:       42,
			},
			want: Flag{
				name:        "intflag",
				alias:       []string{"if", "intflg"},
				shortName:   "I",
				description: "Int Flag",
				value:       42,
			},
		},
		{
			name: "blankName",
			args: args{
				name:      "",
				alias:     []string{},
				shortName: "",
				value:     42,
			},
			want:        Flag{},
			shouldPanic: true,
		},
		{
			name: "singleCharacterName",
			args: args{
				name:      "x",
				alias:     []string{},
				shortName: "",
				value:     42,
			},
			want:        Flag{},
			shouldPanic: true,
		},
		{
			name: "singleRuneName",
			args: args{
				name:      "✓", // check mark U+2713
				alias:     []string{},
				shortName: "",
				value:     42,
			},
			want:        Flag{},
			shouldPanic: true,
		},
		{
			name: "blankAlias",
			args: args{
				name:      "foobar",
				alias:     []string{"barfoo", "\n"},
				shortName: "",
				value:     42,
			},
			want:        Flag{},
			shouldPanic: true,
		},
		{
			name: "singleCharacterAlias",
			args: args{
				name:      "foobar",
				alias:     []string{"barfoo", "x"},
				shortName: "",
				value:     42,
			},
			want:        Flag{},
			shouldPanic: true,
		},
		{
			name: "singleRuneAlias",
			args: args{
				name:      "foobar",
				alias:     []string{"barfoo", "✓"}, // check mark U+2713
				shortName: "",
				value:     42,
			},
			want:        Flag{},
			shouldPanic: true,
		},
		{
			name: "blankShortName",
			args: args{
				name:      "Foobar",
				alias:     []string{"barFoo"},
				shortName: " ",
				value:     42,
			},
			want:        Flag{name: "foobar", alias: []string{"barfoo"}, shortName: "", value: 42},
			shouldPanic: false,
		},
		{
			name: "longShortName",
			args: args{
				name:      "foobar",
				alias:     []string{"barfoo", "\n"},
				shortName: "to",
				value:     42,
			},
			want:        Flag{},
			shouldPanic: true,
		},
		{
			name: "dupNameAndAlias",
			args: args{
				name:      "foobar",
				alias:     []string{"foobar", "barfoo"},
				shortName: "t",
				value:     42,
			},
			want:        Flag{},
			shouldPanic: true,
		},
		{
			name: "dupAliasAndAlias",
			args: args{
				name:      "foobar",
				alias:     []string{"barfoo", "barfoo"},
				shortName: "t",
				value:     42,
			},
			want:        Flag{},
			shouldPanic: true,
		},
		{
			name: "dupShortName",
			args: args{
				name:      "foobar",
				alias:     []string{"barfoo", "foobarfoo", "t"},
				shortName: "t",
				value:     42,
			},
			want:        Flag{},
			shouldPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				defer func() {
					rec := recover()
					if rec != nil {
						t.Logf("Ignoring panic: %v\n", rec)
					}
				}()
				v := tt.args.value.(int)                                                            // panics if value is not an int
				_ = NewFlag(tt.args.name, tt.args.alias, tt.args.shortName, tt.args.description, v) // should panic
				t.Errorf("NewFlag test %s should have panicked, but didn't", tt.name)
			} else {
				flg := NewFlag(tt.args.name, tt.args.alias, tt.args.shortName, tt.args.description, tt.args.value.(int))
				if flg.Name() != tt.want.name {
					t.Errorf("Name() got %s, want %s", flg.Name(), tt.want.name)
				}
				if !reflect.DeepEqual(flg.Alias(), tt.want.alias) {
					t.Errorf("Alias() got %v, want %v", flg.Alias(), tt.want.alias)
				}
				if flg.ShortName() != tt.want.shortName {
					t.Errorf("ShortName() got %s, want %s", flg.ShortName(), tt.want.shortName)
				}
				if flg.Description() != tt.want.description {
					t.Errorf("Description() got %s, want %s", flg.Description(), tt.want.description)
				}
			}
		})
	}
}

func TestGetValueAny(t *testing.T) {
	tests := []struct {
		name string
		flag *Flag
		want any
	}{
		{
			name: "intValue",
			flag: &Flag{
				name:        "intFlag",
				alias:       []string{"iF"},
				shortName:   "i",
				description: "An integer flag",
				value:       42,
			},
			want: 42,
		},
		{
			name: "stringValue",
			flag: &Flag{
				name:        "stringFlag",
				alias:       []string{"sF"},
				shortName:   "s",
				description: "A string flag",
				value:       "hello",
			},
			want: "hello",
		},
		{
			name: "boolValue",
			flag: &Flag{
				name:        "boolFlag",
				alias:       []string{"bF"},
				shortName:   "b",
				description: "A boolean flag",
				value:       true,
			},
			want: true,
		},
		{
			name: "int64Value",
			flag: &Flag{
				name:        "int64Flag",
				alias:       []string{"i64F"},
				shortName:   "i64",
				description: "An int64 flag",
				value:       int64(64),
			},
			want: int64(64),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.flag.GetValueAny(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetValueAny() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetValueOK(t *testing.T) {
	tests := []struct {
		name   string
		flag   *Flag
		wantOK bool
	}{
		{
			name: "intValue",
			flag: &Flag{
				name:        "intFlag",
				alias:       []string{"iF"},
				shortName:   "i",
				description: "An integer flag",
				value:       42,
			},
			wantOK: true,
		},
		{
			name: "stringValue",
			flag: &Flag{
				name:        "stringFlag",
				alias:       []string{"sF"},
				shortName:   "s",
				description: "A string flag",
				value:       "hello",
			},
			wantOK: false,
		},
		{
			name: "boolValue",
			flag: &Flag{
				name:        "boolFlag",
				alias:       []string{"bF"},
				shortName:   "b",
				description: "A boolean flag",
				value:       true,
			},
			wantOK: false,
		},
		{
			name: "int64Value",
			flag: &Flag{
				name:        "int64Flag",
				alias:       []string{"i64F"},
				shortName:   "i64",
				description: "An int64 flag",
				value:       int64(64),
			},
			wantOK: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := GetValueOK[int](tt.flag)
			if ok != tt.wantOK {
				t.Errorf("GetValueOK() ok = %v, want %v", ok, tt.wantOK)
			}
			if ok && got != tt.flag.value {
				t.Errorf("GetValueOK() value = %v, want %v", got, tt.flag.value)
			}
		})
	}
}

func TestGetValue(t *testing.T) {
	tests := []struct {
		name        string
		flag        *Flag
		shouldPanic bool
	}{
		{
			name: "intValue",
			flag: &Flag{
				name:        "intFlag",
				alias:       []string{"iF"},
				shortName:   "i",
				description: "An integer flag",
				value:       42,
			},
			shouldPanic: false,
		},
		{
			name: "stringValue",
			flag: &Flag{
				name:        "stringFlag",
				alias:       []string{"sF"},
				shortName:   "s",
				description: "A string flag",
				value:       "hello",
			},
			shouldPanic: true,
		},
		{
			name: "boolValue",
			flag: &Flag{
				name:        "boolFlag",
				alias:       []string{"bF"},
				shortName:   "b",
				description: "A boolean flag",
				value:       true,
			},
			shouldPanic: true,
		},
		{
			name: "int64Value",
			flag: &Flag{
				name:        "int64Flag",
				alias:       []string{"i64F"},
				shortName:   "i64",
				description: "An int64 flag",
				value:       int64(64),
			},
			shouldPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				defer func() {
					rec := recover()
					if rec != nil {
						t.Logf("Ignoring panic: %v\n", rec)
					}
				}()
				GetValue[int](tt.flag) // should panic
				t.Errorf("GetValue test %s should have panicked, but didn't", tt.name)
			} else {
				if got := GetValue[int](tt.flag); got != tt.flag.value {
					t.Errorf("GetValue() = %v, want %v", got, tt.flag.value)
				}
			}
		})
	}
}

func TestAddFlag(t *testing.T) {
	type args struct {
		f           Flags
		name        string
		alias       []string
		shortName   string
		description string
		value       any
	}
	tests := []struct {
		name        string
		args        args
		want        Flags
		shouldPanic bool
	}{
		{
			name: "addIntFlag",
			args: args{
				f:           NewFlags(),
				name:        "IntFlag",
				alias:       []string{"IF", "IntFlg"},
				shortName:   "I",
				description: "Int Flag",
				value:       42,
			},
			want:        NewFlags().AddFlag(NewFlag("intflag", []string{"if", "intflg"}, "I", "Int Flag", 42)),
			shouldPanic: false,
		},
		{
			name: "nilFlags",
			args: args{
				f:         nil,
				name:      "",
				alias:     []string{},
				shortName: "",
				value:     42,
			},
			want:        nil,
			shouldPanic: true,
		},
		{
			name: "dupFlagName",
			args: args{
				f:         NewFlags().AddFlag(NewFlag("foobar", []string{}, "", "some description", 42)),
				name:      "foobar",
				alias:     []string{},
				shortName: "",
				value:     42,
			},
			want:        nil,
			shouldPanic: true,
		},
		{
			name: "dupFlagName2",
			args: args{
				f:           NewFlags().AddFlag(NewFlag("foobar✓", []string{}, "", "some description", 42)),
				name:        "foobar✓",
				alias:       []string{},
				shortName:   "f",
				description: "some description",
				value:       42,
			},
			want:        nil,
			shouldPanic: true,
		},
		{
			name: "dupFlagShortname",
			args: args{
				f:           NewFlags().AddFlag(NewFlag("foobar", []string{}, "f", "some description", 42)),
				name:        "foobar2",
				alias:       []string{},
				shortName:   "f",
				description: "some description",
				value:       42,
			},
			want:        nil,
			shouldPanic: true,
		},
		{
			name: "dupFlagShortnameCase",
			args: args{
				f:           NewFlags().AddFlag(NewFlag("foobar", []string{}, "F", "some description", 42)),
				name:        "foobar2",
				alias:       []string{},
				shortName:   "f",
				description: "some description",
				value:       42,
			},
			want:        NewFlags().AddFlag(NewFlag("foobar", []string{}, "F", "some description", 42)).AddFlag(NewFlag("foobar2", []string{}, "f", "some description", 42)),
			shouldPanic: false,
		},
		{
			name: "dupFlagAlias",
			args: args{
				f:         NewFlags().AddFlag(NewFlag("foobar", []string{"able", "baker", "charlie"}, "", "some description", 42)),
				name:      "foobar2",
				alias:     []string{"alpha", "bravo", "charlie"},
				shortName: "",
				value:     42,
			},
			want:        nil,
			shouldPanic: true,
		},
		{
			name: "dupFlagAlias2",
			args: args{
				f:         NewFlags().AddFlag(NewFlag("foobar", []string{"able", "baker", "charlie"}, "", "some description", 42)),
				name:      "foobar2",
				alias:     []string{"alpha", "bravo", "foobar"},
				shortName: "",
				value:     42,
			},
			want:        nil,
			shouldPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flgs := tt.args.f
			if tt.shouldPanic {
				defer func() {
					rec := recover()
					if rec != nil {
						t.Logf("Ignoring panic: %v\n", rec)
					}
				}()
				//TODO test other value types
				v := tt.args.value.(int) // panics if value is not an int
				flgs.AddFlag(NewFlag(tt.args.name, tt.args.alias, tt.args.shortName, tt.args.description, v))
				t.Errorf("AddFlag test %s should have panicked, but didn't", tt.name)
			} else {
				switch v := tt.args.value.(type) {
				case int:
					newFlags := flgs.AddFlag(NewFlag(tt.args.name, tt.args.alias, tt.args.shortName, tt.args.description, v))
					if !reflect.DeepEqual(newFlags, tt.want) {
						t.Errorf("AddFlag() = %v, want %v", newFlags, tt.want)
					}
				default:
					t.Errorf("Unknown value type %T val %v", tt.args.value, tt.args.value)
				}
			}
		})
	}
}

func TestGetFlagOK(t *testing.T) {
	tests := []struct {
		name   string
		f      Flags
		n      string
		wantOK bool
		want   *Flag
	}{
		{
			name: "getFlagOK valid",
			f: Flags{
				"flag1": &Flag{
					name:        "flag1",
					alias:       []string{"f1", "flg1"},
					shortName:   "f",
					description: "Flag 1",
					value:       1,
				},
			},
			n:      "flag1",
			wantOK: true,
			want: &Flag{
				name:        "flag1",
				alias:       []string{"f1", "flg1"},
				shortName:   "f",
				description: "Flag 1",
				value:       1,
			},
		},
		{
			name: "getFlagOK invalid",
			f: Flags{
				"flag2": &Flag{
					name:        "flag2",
					alias:       []string{"f2", "flg2"},
					shortName:   "g",
					description: "Flag 2",
					value:       2,
				},
			},
			n:      "flag1",
			wantOK: false,
			want:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := GetFlagOK(tt.f, "flag1")
			if ok != tt.wantOK {
				t.Errorf("GetFlagOK() ok = %v, want %v", ok, tt.wantOK)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFlagOK() value = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFlag(t *testing.T) {
	tests := []struct {
		name        string
		f           Flags
		n           string
		shouldPanic bool
		want        *Flag
	}{
		{
			n: "flag1",
			f: Flags{
				"flag1": &Flag{
					name:        "flag1",
					alias:       []string{"f1", "flg1"},
					shortName:   "f",
					description: "Flag 1",
					value:       1,
				},
			},
			want: &Flag{
				name:        "flag1",
				alias:       []string{"f1", "flg1"},
				shortName:   "f",
				description: "Flag 1",
				value:       1,
			},
			shouldPanic: false,
		},
		{
			name: "getFlag invalid",
			f: Flags{
				"flag2": &Flag{
					name:        "flag2",
					alias:       []string{"f2", "flg2"},
					shortName:   "g",
					description: "Flag 2",
					value:       2,
				},
			},
			n:           "flag1",
			want:        nil,
			shouldPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				defer func() {
					rec := recover()
					if rec != nil {
						t.Logf("Ignoring panic: %v\n", rec)
					}
				}()
				GetFlag(tt.f, tt.n)
				t.Errorf("GetFlag test %s should have panicked, but didn't", tt.name)
			} else {
				if got := GetFlag(tt.f, tt.n); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("GetFlag() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestFlags_String(t *testing.T) {
	tests := []struct {
		name string
		fs   Flags
		want string
	}{
		{
			name: "stringTest",
			fs:   Flags{},
			want: "Flags:\nName Short Aliases Default Type Description\n",
		},
		{
			name: "stringTest2",
			fs: Flags{
				"flg1": &Flag{
					name:        "flg1",
					alias:       []string{"alias1"},
					shortName:   "",
					description: "",
					value:       "val",
				},
			},
			want: "Flags:\nName Short Aliases Default Type   Description\nflg1       alias1  val     string \n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fs.String(); got != tt.want {
				t.Errorf("Flags.String() = %s, want %s", strconv.Quote(got), strconv.Quote(tt.want))
			}
		})
	}
}

func TestFindFlag(t *testing.T) {
	tests := []struct {
		name     string
		flags    Flags
		search   string
		expected *Flag
	}{
		{
			name: "findByName",
			flags: Flags{
				"flag1\u2713": &Flag{
					name:        "flag1\u2713",
					alias:       []string{"f1", "flg1"},
					shortName:   "f",
					description: "Flag 1",
					value:       1,
				},
			},
			search: " Flag1\u2713",
			expected: &Flag{
				name:        "flag1\u2713",
				alias:       []string{"f1", "flg1"},
				shortName:   "f",
				description: "Flag 1",
				value:       1,
			},
		},
		{
			name: "findByAlias",
			flags: Flags{
				"flag1": &Flag{
					name:        "flag1",
					alias:       []string{"f1", "flg1"},
					shortName:   "f",
					description: "Flag 1",
					value:       1,
				},
			},
			search: "f1",
			expected: &Flag{
				name:        "flag1",
				alias:       []string{"f1", "flg1"},
				shortName:   "f",
				description: "Flag 1",
				value:       1,
			},
		},
		{
			name: "notFound",
			flags: Flags{
				"flag1": &Flag{
					name:        "flag1",
					alias:       []string{"f1", "flg1"},
					shortName:   "f",
					description: "Flag 1",
					value:       1,
				},
			},
			search:   "nonexistent",
			expected: nil,
		},
		{
			name: "findByShortName",
			flags: Flags{
				"flag1": &Flag{
					name:        "flag1",
					alias:       []string{"f1", "flg1"},
					shortName:   "f",
					description: "Flag 1",
					value:       1,
				},
			},
			search: "f",
			expected: &Flag{
				name:        "flag1",
				alias:       []string{"f1", "flg1"},
				shortName:   "f",
				description: "Flag 1",
				value:       1,
			},
		},
		{
			name: "findByAliasCaseInsensitive",
			flags: Flags{
				"flag1": &Flag{
					name:        "flag1",
					alias:       []string{"f1", "flg1"},
					shortName:   "f",
					description: "Flag 1",
					value:       1,
				},
			},
			search: "F1",
			expected: &Flag{
				name:        "flag1",
				alias:       []string{"f1", "flg1"},
				shortName:   "f",
				description: "Flag 1",
				value:       1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.flags.FindFlag(tt.search)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("FindFlag() = %v, want %v", result, tt.expected)
			}
		})
	}
}
