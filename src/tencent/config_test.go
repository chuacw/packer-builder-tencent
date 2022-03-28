package tencent

import (
	"reflect"
	"testing"
)

// returns the minimal configuration that will supposedly work successfully
func createWorkingConfig() (map[string]interface{}, string) {
	filename := TempFileName()
	return map[string]interface{}{
		"ImageID":      "img-3wnd9xpl",
		"Placement":    map[string]interface{}{"Zone": "ap-guangzhou-1"},
		"KeyName":      filename,
		"SecretID":     "secret",
		"SecretKey":    "key",
		"ssh_username": "ubuntu",
	}, filename
}

// returns an empty configuration that will cause errors
func createEmptyConfig() map[string]interface{} {
	return map[string]interface{}{}
}

// This tests the NewConfig function in config.go works as expected.
func TestNewConfig(t *testing.T) {
	// Test an empty configuration
	raw1 := createEmptyConfig()
	_, warns, errs := NewConfig(raw1)
	CheckConfigHasErrors(t, warns, errs)

	// test a default working configuration
	raw2, _ := createWorkingConfig()
	_, warns, errs = NewConfig(raw2)
	CheckConfigIsOk(t, warns, errs)
}

func TestConfig_Keys(t *testing.T) {
	config := new(Config)
	result := make(map[string]string)
	tests := []struct {
		name string
		c    *Config
		want map[string]string
	}{
		{"Keys() test case 1", config, result},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Keys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}
