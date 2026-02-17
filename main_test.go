package main

import "testing"

func TestBumpVersion(t *testing.T) {
	tests := []struct {
		name       string
		version    string
		major      uint64
		minor      uint64
		patch      uint64
		exact      string
		want       string
		wantErr    bool
	}{
		{"patch", "1.2.3", 0, 0, 1, "", "1.2.4", false},
		{"minor", "1.2.3", 0, 1, 0, "", "1.3.0", false},
		{"major", "1.2.3", 1, 0, 0, "", "2.0.0", false},
		{"set", "1.2.3", 0, 0, 0, "5.0.0", "5.0.0", false},
		{"set lower", "1.2.3", 0, 0, 0, "0.1.0", "", true},
		{"invalid version", "abc", 0, 0, 1, "", "", true},
		{"invalid exact", "1.2.3", 0, 0, 0, "bad", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bumpVersion(tt.version, tt.major, tt.minor, tt.patch, tt.exact)
			if (err != nil) != tt.wantErr {
				t.Errorf("bumpVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("bumpVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
