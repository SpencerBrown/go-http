package flag

import (
	"reflect"
	"strconv"
	"testing"
)

func something() {
	flgs := NewFlags()
	flg := NewFlag("flag1", []string{"f1", "flg1"}, "f", "Flag 1", 1)
	_ = flgs.AddFlag(flg)
	z := y.AddFlag(NewFlag("flag1", nil, "", "", "value1"))
	x := NewFlags().AddFlag(NewFlag("flag1", nil, "", "", "value1"))
}

func TestNewFlags(t *testing.T) {
	t.Run("NewFlags", func(t *testing.T) {
		want := make(Flags)
		if got := NewFlags(); !reflect.DeepEqual(got, want) {
			t.Errorf("NewFlags() = %v, want %v", got, want)
		}
	})
}

// TestNewFlag also tests GetFlagOK, GetFlag, GetValueOK, and GetValue
// as well as f.Name(), f.Alias(), f.ShortName(), f.Description()
func TestNewFlag(t *testing.T) {
	type args struct {
		f           Flags
		setupFlags  func(f Flags) Flags
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
			name: "nilFlags",
			args: args{
				f:         nil,
				name:      "",
				alias:     []string{},
				shortName: "",
				value:     42,
			},
			want:        Flag{},
			shouldPanic: true,
		},
		{
			name: "blankName",
			args: args{
				f:         NewFlags(),
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
				f:         NewFlags(),
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
				f:         NewFlags(),
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
				f:         NewFlags(),
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
				f:         NewFlags(),
				name:      "foobar",
				alias:     []string{"barfoo", "x"},
				shortName: "",
				value:     42,
			},
			want:        Flag{},
			shouldPanic: true,
		},
		{
			name: "blankShortName",
			args: args{
				f:         NewFlags(),
				name:      "foobar",
				alias:     []string{"barfoo", "\n"},
				shortName: " ",
				value:     42,
			},
			want:        Flag{},
			shouldPanic: true,
		},
		{
			name: "longShortName",
			args: args{
				f:         NewFlags(),
				name:      "foobar",
				alias:     []string{"barfoo", "\n"},
				shortName: "to",
				value:     42,
			},
			want:        Flag{},
			shouldPanic: true,
		},
		{
			name: "dupAlias",
			args: args{
				f:         NewFlags(),
				name:      "foobar",
				alias:     []string{"foobar", "barfoo"},
				shortName: "t",
				value:     42,
			},
			want:        Flag{},
			shouldPanic: true,
		},
		{
			name: "dupAlias2",
			args: args{
				f:         NewFlags(),
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
				f:         NewFlags(),
				name:      "foobar",
				alias:     []string{"barfoo", "foobarfoo", "t"},
				shortName: "t",
				value:     42,
			},
			want:        Flag{},
			shouldPanic: true,
		},
		{
			name: "dupFlagName",
			args: args{
				f: NewFlags(),
				setupFlags: func(f Flags) Flags {
					NewFlag(f, "foobar", []string{}, "", "some description", 42)
					return f
				},
				name:      "foobar",
				alias:     []string{},
				shortName: "",
				value:     42,
			},
			want:        Flag{},
			shouldPanic: true,
		},
		{
			name: "dupFlagName2",
			args: args{
				f: NewFlags(),
				setupFlags: func(f Flags) Flags {
					NewFlag(f, "foobar\u2713", []string{}, "f", "some description", 42)
					return f
				},
				name:      "foobar✓",
				alias:     []string{},
				shortName: "f",
				value:     42,
			},
			want:        Flag{},
			shouldPanic: true,
		},
		{
			name: "dupFlagShortname",
			args: args{
				f: NewFlags(),
				setupFlags: func(f Flags) Flags {
					NewFlag(f, "foobar", []string{}, "f", "some description", 42)
					return f
				},
				name:      "foobar2",
				alias:     []string{},
				shortName: "f",
				value:     42,
			},
			want:        Flag{},
			shouldPanic: true,
		},
		{
			name: "dupFlagAlias",
			args: args{
				f: NewFlags(),
				setupFlags: func(f Flags) Flags {
					NewFlag(f, "foobar", []string{"able", "baker", "charlie"}, "", "some description", 42)
					return f
				},
				name:      "foobar2",
				alias:     []string{"alpha", "bravo", "charlie"},
				shortName: "",
				value:     42,
			},
			want:        Flag{},
			shouldPanic: true,
		},
		{
			name: "dupFlagAlias2",
			args: args{
				f: NewFlags(),
				setupFlags: func(f Flags) Flags {
					NewFlag(f, "foobar", []string{"able", "baker", "charlie"}, "", "some description", 42)
					return f
				},
				name:      "foobar2",
				alias:     []string{"alpha", "bravo", "foobar"},
				shortName: "",
				value:     42,
			},
			want:        Flag{},
			shouldPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var flgs Flags
			if tt.args.setupFlags == nil {
				flgs = tt.args.f
			} else {
				flgs = tt.args.setupFlags(tt.args.f)
			}
			if tt.shouldPanic {
				defer func() {
					rec := recover()
					if rec != nil {
						t.Logf("Ignoring panic: %v\n", rec)
					}
				}()
				// the following line should panic, either because the value conversion to string panics,
				// or because one of the args is invalid for some reason (which varies per test)
				v := tt.args.value.(int) // panics if value is not an int
				_ = NewFlag(flgs, tt.args.name, tt.args.alias, tt.args.shortName, tt.args.description, v)
				t.Errorf("NewFlag test %s should have panicked, but didn't", tt.name)
			} else {
				switch v := tt.args.value.(type) {
				case int:
					newFlg := NewFlag(flgs, tt.args.name, tt.args.alias, tt.args.shortName, tt.args.description, v)
					flg, ok := GetFlagOK(tt.args.f, tt.args.name)
					if ok {
						if !reflect.DeepEqual(*flg, tt.want) {
							t.Errorf("GetFlagOK for %s got Flag %v, want %v", tt.args.name, flg, tt.want)
						}
						if !reflect.DeepEqual(*flg, *newFlg) {
							t.Errorf("GetFlagOK for %s returned wrong flag", tt.args.name)
						}
						flg2 := GetFlag(tt.args.f, tt.args.name)
						if !reflect.DeepEqual(*flg2, tt.want) {
							t.Errorf("GetFlag for %s got Flag %v, want %v", tt.args.name, flg2, tt.want)
						}
						val, ok := GetValueOK[int](flg)
						if !ok {
							t.Errorf("GetValueOK[int] for %s returned not ok", tt.args.name)
						}
						if val != tt.want.value {
							t.Errorf("GetValue[int] got %d, want %d", val, tt.want.value)
						}
					} else {
						t.Errorf("GetFlagOK for %s returned not ok", tt.args.name)
					}
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
					val := GetValue[int](flg) // should not panic
					if val != tt.want.value {
						t.Errorf("GetValue[int] got %d, want %d", val, tt.want.value)
					}
					defer func() { _ = recover() }() // ignore panic
					GetValue[int64](flg)             // should panic
					t.Error("GetValue[int64] should have panicked, but didn't")
				default:
					t.Errorf("Unknown value type %T val %v", tt.args.value, tt.args.value)
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

func TestGetFlag(t *testing.T) {
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
		want        *Flag
		shouldPanic bool
	}{
		{
			name: "nonexistentFlag",
			args: args{
				f:    NewFlags(),
				name: "NonexistentFlag",
			},
			want:        nil,
			shouldPanic: true,
		},
		{
			name: "existingFlag",
			args: args{
				f:           NewFlags(),
				name:        "ExistingFlag",
				alias:       []string{"EF", "ExFlg"},
				shortName:   "E",
				description: "Existing Flag",
				value:       42,
			},
			want: &Flag{
				name:        "existingflag",
				alias:       []string{"ef", "exflg"},
				shortName:   "E",
				description: "Existing Flag",
				value:       42,
			},
			shouldPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				defer func() { rec := recover(); /*fmt.Printf("Ignoring panic: %v\n", rec);*/ _ = rec }() // ignore panic
				_ = GetFlag(tt.args.f, tt.args.name)
				t.Error("GetFlag should have panicked, but didn't")
			} else {
				switch v := tt.args.value.(type) {
				case int:
					NewFlag(tt.args.f, tt.args.name, tt.args.alias, tt.args.shortName, tt.args.description, v)
					if got := GetFlag(tt.args.f, tt.args.name); !reflect.DeepEqual(got, tt.want) {
						t.Errorf("GetFlag() = %v, want %v", got, tt.want)
					}
				default:
					t.Errorf("Unknown value type %T val %v", tt.args.value, tt.args.value)

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
			search:   "f",
			expected:  &Flag{
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