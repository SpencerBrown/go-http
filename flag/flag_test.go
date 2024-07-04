package flag

import (
	"reflect"
	"testing"
)

func TestNewFlags(t *testing.T) {
	t.Run("NewFlags", func(t *testing.T) {
		if got := NewFlags(); !reflect.DeepEqual(got, map[string]Flag{}) {
			t.Errorf("NewFlags() = %v, want %v", got, map[string]Flag{})
		}
	})
}

// TestNewFlag also tests GetValue and GetValueMust
func TestNewFlag(t *testing.T) {
	type args struct {
		f           Flags
		name        string
		shortName   string
		description string
		value       any
	}
	tests := []struct {
		name string
		args args
		want Flag
	}{
		{
			name: "newIntFlag",
			args: args{
				f:           NewFlags(),
				name:        "IntFlag",
				shortName:   "I",
				description: "Int Flag",
				value:       42,
			},
			want: Flag{
				ShortName:   "I",
				Description: "Int Flag",
				Value:       42,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.args.value.(type) {
			case int:
				NewFlag(tt.args.f, tt.args.name, tt.args.shortName, tt.args.description, v)
				flg, val, ok := GetValue[int](tt.args.f, tt.args.name)
				if ok {
					if !reflect.DeepEqual(flg, tt.want) {
						t.Errorf("GetValue[int] got Flag %v, want %v", flg, tt.want)
					}
					if val != tt.want.Value {
						t.Errorf("GetValue[int] got %d, want %d", val, tt.want.Value)
					}
				} else {
					t.Errorf("GetValue for %s returned not ok", tt.args.name)
				}
				flg, val = GetValueMust[int](tt.args.f, tt.args.name) // should not panic
				if !reflect.DeepEqual(flg, tt.want) {
					t.Errorf("GetValueMust[int] got Flag %v, want %v", flg, tt.want)
				}
				if val != tt.want.Value {
					t.Errorf("GetValueMust[int] got %d, want %d", val, tt.want.Value)
				}
				defer func() { _ = recover() }()                    // ignore panic
				_, _ = GetValueMust[int64](tt.args.f, tt.args.name) // should panic
				t.Errorf("GetValueMust[int64] should have paniced, but didn't")
			default:
				t.Errorf("Unknown value type %T val %v", tt.args.value, tt.args.value)
			}
		})
	}
}

func TestGetFlags(t *testing.T) {
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
			if err := GetFlags(tt.args.fs); (err != nil) != tt.wantErr {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fs.String(); got != tt.want {
				t.Errorf("Flags.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
