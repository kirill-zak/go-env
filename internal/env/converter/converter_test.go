package converter

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_strToInt(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Valid integer",
			args: args{
				str: "123",
			},
			want:    123,
			wantErr: false,
		},
		{
			name: "Negative number",
			args: args{
				str: "-456",
			},
			want:    -456,
			wantErr: false,
		},
		{
			name: "Zero value",
			args: args{
				str: "0",
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "Empty string",
			args: args{
				str: "",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Non-integer input",
			args: args{
				str: "abc",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Leading whitespace",
			args: args{
				str: "     789",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Trailing whitespace",
			args: args{
				str: "1011     ",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Float as input",
			args: args{
				str: "123.45",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Maximum positive integer",
			args: args{
				str: "2147483647",
			},
			want:    2147483647,
			wantErr: false,
		},
		{
			name: "Minimum negative integer",
			args: args{
				str: "-2147483648",
			},
			want:    -2147483648,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := strToInt(tt.args.str)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.want, got)
		})
	}
}

func Test_strToDuration(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Duration
		wantErr bool
	}{
		{
			name: "valid duration",
			args: args{
				str: "1h",
			},
			want:    time.Hour,
			wantErr: false,
		},
		{
			name: "invalid format",
			args: args{
				str: "abc",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "negative duration",
			args: args{
				str: "-1m",
			},
			want:    -time.Minute,
			wantErr: false,
		},
		{
			name: "zero value",
			args: args{
				str: "0s",
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "combined units",
			args: args{
				str: "1h30m",
			},
			want:    time.Hour + 30*time.Minute,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := strToDuration(tt.args.str)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.want, got)
		})
	}
}

func Test_strToFloat(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "valid float",
			args: args{
				str: "123.45",
			},
			want:    123.45,
			wantErr: false,
		},
		{
			name: "invalid input",
			args: args{
				str: "not a number",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "scientific notation",
			args: args{
				str: "1.23e+2",
			},
			want:    123,
			wantErr: false,
		},
		{
			name: "integer as string",
			args: args{
				str: "42",
			},
			want:    42,
			wantErr: false,
		},
		{
			name: "empty string",
			args: args{
				str: "",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := strToFloat(tt.args.str)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.want, got)
		})
	}
}

func Test_strToSlice(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "simple comma separated list",
			args: args{
				str: "apple,banana,cherry",
			},
			want: []string{"apple", "banana", "cherry"},
		},
		{
			name: "list with spaces",
			args: args{
				str: " apple , banana , cherry ",
			},
			want: []string{"apple", "banana", "cherry"},
		},
		{
			name: "single element",
			args: args{
				str: "orange",
			},
			want: []string{"orange"},
		},
		{
			name: "empty string",
			args: args{
				str: "",
			},
			want: []string{""},
		},
		{
			name: "extra commas",
			args: args{
				str: ",,,grapefruit",
			},
			want: []string{"", "", "", "grapefruit"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := strToSlice(tt.args.str)

			require.ElementsMatch(t, tt.want, got)
		})
	}
}

func Test_strToBool(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "true case",
			args: args{
				str: "true",
			},
			want: true,
		},
		{
			name: "false case",
			args: args{
				str: "false",
			},
			want: false,
		},
		{
			name: "numeric zero",
			args: args{
				str: "0",
			},
			want: false,
		},
		{
			name: "non-empty string",
			args: args{
				str: "random text",
			},
			want: true,
		},
		{
			name: "empty string",
			args: args{
				str: "",
			},
			want: false,
		},
		{
			name: "mixed case",
			args: args{
				str: "FaLsE",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := strToBool(tt.args.str)

			require.Equal(t, tt.want, got)
		})
	}
}

func Test_strToURLSlice(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name     string
		args     args
		wantUrls []*url.URL
		wantErr  bool
	}{
		{
			name: "valid URLs",
			args: args{
				str: "https://example.com, https://domain.com",
			},
			wantUrls: []*url.URL{
				{Host: "example.com", Scheme: "https", Path: ""},
				{Host: "domain.com", Scheme: "https", Path: ""},
			},
			wantErr: false,
		},
		{
			name: "one valid URL",
			args: args{
				str: "http://localhost",
			},
			wantUrls: []*url.URL{{Host: "localhost", Scheme: "http"}},
			wantErr:  false,
		},
		{
			name: "empty string",
			args: args{
				str: "",
			},
			wantUrls: []*url.URL{{}},
			wantErr:  false,
		},
		{
			name: "malformed URL",
			args: args{
				str: "htts://bad .url",
			},
			wantUrls: nil,
			wantErr:  true,
		},
		{
			name: "mix of valid and invalid URLs",
			args: args{
				str: "htts://bad .url, https://valid.example.com",
			},
			wantUrls: nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUrls, err := strToURLSlice(tt.args.str)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.wantUrls, gotUrls)
		})
	}
}
