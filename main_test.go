package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestParseField(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		min     int
		max     int
		want    []int
		wantErr bool
		errMsg  string
	}{
		// Happy Paths
		{
			name:    "single number",
			input:   "5",
			min:     0,
			max:     59,
			want:    []int{5},
			wantErr: false,
		},
		{
			name:    "asterisk",
			input:   "*",
			min:     0,
			max:     23,
			want:    []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
			wantErr: false,
		},
		{
			name:    "step with asterisk",
			input:   "*/15",
			min:     0,
			max:     59,
			want:    []int{0, 15, 30, 45},
			wantErr: false,
		},
		{
			name:    "range",
			input:   "1-5",
			min:     0,
			max:     7,
			want:    []int{1, 2, 3, 4, 5},
			wantErr: false,
		},
		{
			name:    "list",
			input:   "1,15",
			min:     1,
			max:     31,
			want:    []int{1, 15},
			wantErr: false,
		},
		{
			name:    "range with step",
			input:   "0-10/2",
			min:     0,
			max:     59,
			want:    []int{0, 2, 4, 6, 8, 10},
			wantErr: false,
		},
		// Sad Paths
		{
			name:    "invalid number",
			input:   "abc",
			min:     0,
			max:     59,
			want:    nil,
			wantErr: true,
			errMsg:  "invalid number: abc",
		},
		{
			name:    "invalid range format",
			input:   "1-",
			min:     0,
			max:     23,
			want:    nil,
			wantErr: true,
			errMsg:  "invalid range format: 1-",
		},
		{
			name:    "min greater than max",
			input:   "5-2",
			min:     0,
			max:     59,
			want:    nil,
			wantErr: true,
			errMsg:  "min > max in range: 5-2",
		},
		{
			name:    "invalid step format",
			input:   "*/-1",
			min:     0,
			max:     59,
			want:    nil,
			wantErr: true,
			errMsg:  "invalid step value: -1",
		},
		{
			name:    "invalid step syntax",
			input:   "*/",
			min:     0,
			max:     59,
			want:    nil,
			wantErr: true,
			errMsg:  "invalid step format: */",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseField(tt.input, tt.min, tt.max)
			if tt.wantErr {
				if err == nil {
					t.Errorf("parseField(%q) expected error, got nil", tt.input)
					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("parseField(%q) error = %v, want error containing %q", tt.input, err, tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("parseField(%q) unexpected error: %v", tt.input, err)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("parseField(%q) = %v, want %v", tt.input, got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("parseField(%q) = %v, want %v", tt.input, got, tt.want)
					return
				}
			}
		})
	}
}

func TestMainOutput(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantOut string
		wantErr bool
	}{
		{
			name: "valid cron string",
			args: []string{"cron_parser", "*/15 0 1,15 * 1-5 /usr/bin/find"},
			wantOut: `minute        0 15 30 45
hour          0
day of month  1 15
month         1 2 3 4 5 6 7 8 9 10 11 12
day of week   1 2 3 4 5
command       /usr/bin/find
`,
			wantErr: false,
		},
		{
			name:    "no arguments",
			args:    []string{"cron_parser"},
			wantOut: "Usage: cron_parser \"cron-string\"\nExample: cron_parser \"*/15 0 1,15 * 1-5 /usr/bin/find\"\n",
			wantErr: true,
		},
		{
			name:    "invalid cron string",
			args:    []string{"cron_parser", "*/15 0 1,15 *"}, // Kurang field
			wantOut: "Invalid input: cron string must have 5 time fields and a command\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args

			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			defer func() { os.Stdout = old }()

			var exitCode int
			func() {
				defer func() {
					if r := recover(); r != nil {
						if code, ok := r.(int); ok {
							exitCode = code
						}
					}
				}()
				main()
			}()

			w.Close()
			var buf bytes.Buffer
			_, _ = buf.ReadFrom(r)

			if got := buf.String(); got != tt.wantOut {
				t.Errorf("main() output = %q, want %q", got, tt.wantOut)
			}

			if tt.wantErr && exitCode == 0 {
				t.Errorf("main() expected to exit with error, but exited with code 0")
			}
			if !tt.wantErr && exitCode != 0 {
				t.Errorf("main() unexpected exit code %d", exitCode)
			}
		})
	}
}
