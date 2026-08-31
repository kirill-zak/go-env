package converter

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAsGoType(t *testing.T) {
	type args struct {
		typ string
		val string
	}
	tests := []struct {
		name     string
		args     args
		wantResT string
		wantResV interface{}
		wantErr  bool
	}{
		{
			name: "valid duration",
			args: args{
				typ: "duration",
				val: "1h",
			},
			wantResT: "time.Duration",
			wantResV: time.Hour,
			wantErr:  false,
		},
		{
			name: "string slice",
			args: args{
				typ: "string_slice",
				val: "a,b,c",
			},
			wantResT: "[]string",
			wantResV: []string{"a", "b", "c"},
			wantErr:  false,
		},
		{
			name: "empty type",
			args: args{
				typ: "",
				val: "test",
			},
			wantResT: "string",
			wantResV: "test",
			wantErr:  false,
		},
		{
			name: "valid float",
			args: args{
				typ: "float",
				val: "123.45",
			},
			wantResT: "float64",
			wantResV: 123.45,
			wantErr:  false,
		},
		{
			name: "valid int",
			args: args{
				typ: "int",
				val: "42",
			},
			wantResT: "int",
			wantResV: 42,
			wantErr:  false,
		},
		{
			name: "valid bool",
			args: args{
				typ: "bool",
				val: "true",
			},
			wantResT: "bool",
			wantResV: true,
			wantErr:  false,
		},
		{
			name: "valid single URL",
			args: args{
				typ: "url",
				val: "https://example.com",
			},
			wantResT: "*url.URL",
			wantResV: &url.URL{Scheme: "https", Host: "example.com"},
			wantErr:  false,
		},
		{
			name: "valid URL slice",
			args: args{
				typ: "[]url",
				val: "https://example.com,https://domain.com",
			},
			wantResT: "[]*url.URL",
			wantResV: []*url.URL{
				{Scheme: "https", Host: "example.com"},
				{Scheme: "https", Host: "domain.com"},
			},
			wantErr: false,
		},
		{
			name: "invalid type",
			args: args{
				typ: "unknown_type",
				val: "value",
			},
			wantResT: "",
			wantResV: nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResT, gotResV, err := AsGoType(tt.args.typ, tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("AsGoType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Equal(t, tt.wantResT, gotResT)
			require.Equal(t, tt.wantResV, gotResV)
		})
	}
}
