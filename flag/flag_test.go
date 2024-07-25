package flag

import (
	"reflect"
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
func TestNewFlag(t *testing.T) {
	type args struct {
		f           Flags
		name        string
		alias       []string
		shortName   rune
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
				alias:       []string{"IF", "IntFlg"},
				shortName:   'I',
				description: "Int Flag",
				value:       42,
			},
			want: Flag{
				name:        "IntFlag",
				alias:       []string{"IF", "IntFlg"},
				shortName:   'I',
				description: "Int Flag",
				value:       42,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
				val := GetValue[int](flg) // should not panic
				if val != tt.want.value {
					t.Errorf("GetValue[int] got %d, want %d", val, tt.want.value)
				}
				defer func() { _ = recover() }() // ignore panic
				GetValue[int64](flg)             // should panic
				t.Errorf("GetValueint64] should have panicked, but didn't")
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
