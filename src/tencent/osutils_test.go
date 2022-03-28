package tencent

import (
	"os"
	"testing"
)

func Test_pathExists(t *testing.T) {
	dir := os.TempDir()
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"path exists", args{dir}, true},
		{"path doesn't exist", args{dir + "1"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pathExists(tt.args.path); got != tt.want {
				t.Errorf("pathExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirectoryExists(t *testing.T) {
	dir := os.TempDir()
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty directory exists", args{""}, true},
		{"temp dir exists", args{dir}, true},
		{"dir doesn't exist", args{dir + "1"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DirectoryExists(tt.args.path); got != tt.want {
				t.Errorf("DirectoryExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	tempfilename := TempFileName()
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Temp file name doesn't exist", args{tempfilename}, false},
		{"Temp file name doesn't exist", args{tempfilename + "1"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileExists(tt.args.path); got != tt.want {
				t.Errorf("FileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
