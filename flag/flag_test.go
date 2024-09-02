package flag

import (
	"fmt"
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

// TestNewFlag also tests GetFlagOK, GetValueOK, and GetValue
// as well as f.Name(), f.Alias(), f.ShortName(), f.Description()
func TestNewFlag(t *testing.T) {
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
		want        Flag
		shouldPanic bool
	}{
		{
			name: "newIntFlag",
			args: args{
				f:           NewFlags(),
				name:        "IntFlag",
				alias:       []string{"IF", "IntFlg"},
				shortName:   "I",
				description: "Int Flag",
				value:       42,
			},
			want: Flag{
				name:        "IntFlag",
				alias:       []string{"IF", "IntFlg"},
				shortName:   "I",
				description: "Int Flag",
				value:       42,
			},
		},
		{
			name: "blankName",
			args: args{
				f:         NewFlags(),
				name:      "",
				alias:     []string{},
				shortName: "",
				value:     "something",
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
				value:     "something",
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
				value:     "something",
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
				value:     "something",
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
				value:     "something",
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
				value:     "something",
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
				value:     "something",
			},
			want:        Flag{},
			shouldPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				vv := tt.args.value.(string)
				defer func() { rec := recover(); fmt.Printf("Ignoring panic: %v\n", rec); _ = rec }() // ignore panic
				_ = NewFlag(tt.args.f, tt.args.name, tt.args.alias, tt.args.shortName, tt.args.description, vv)
				t.Error("NewFlag should have panicked, but didn't")
			} else {
				switch v := tt.args.value.(type) {
				case int:
					newFlg := NewFlag(tt.args.f, tt.args.name, tt.args.alias, tt.args.shortName, tt.args.description, v)
					flg, ok := GetFlagOK(tt.args.f, tt.args.name)
					if ok {
						if !reflect.DeepEqual(*flg, tt.want) {
							t.Errorf("GetFlagOK for %s got Flag %v, want %v", tt.args.name, flg, tt.want)
						}
						if !reflect.DeepEqual(*flg, *newFlg) {
							t.Errorf("GetFlagOK for %s returned wrong flag", tt.args.name)
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

func TestParseFlags(t *testing.T) {
	type args struct {
		fs Flags
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ParseFlags(tt.args.fs); (err != nil) != tt.wantErr {
				t.Errorf("GetFlags() error = %v, wantErr %v", err, tt.wantErr)
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
