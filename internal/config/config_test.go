package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func createTestConfig(p string) {
	s := `repos:
 - path: /go-sync
   url: https://github.com/charghet/go-sync.git
	 username: user
   password: 123456
   email: user@example.com`
	_, err := os.Stat(filepath.Dir(p))
	if !os.IsExist(err) {
		err = os.MkdirAll(filepath.Dir(p), 0755)
		if err != nil {
			panic("Failed to create directory for test config file: " + err.Error())
		}
	}
	err = os.WriteFile(p, []byte(s), 0644)
	if err != nil {
		panic("Failed to create test config file: " + err.Error())
	}
}

func TestGetConfig(t *testing.T) {
	p := "../../test/config/config.yaml"
	createTestConfig(p)
	SetPath(p)
	config := GetConfig()
	if config == nil {
		t.Error("Expected config to be loaded, but got nil")
		return
	}
	b, err := json.Marshal(config)
	if err != nil {
		t.Errorf("Failed to marshal config to JSON: %v", err)
	}
	t.Logf("Config: %s", b)
	if len(config.Repos) == 0 {
		t.Error("Expected at least one repo in config, but got none")
	}
}
