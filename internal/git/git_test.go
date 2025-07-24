package git

import (
	"testing"

	"github.com/charghet/go-sync/internal/config"
)

func TestInit(t *testing.T) {
	config.SetPath("../../config.yaml")
	r := NewGitRepo(config.GetConfig().Repos[0])
	err := r.Open(true)
	if err != nil {
		t.Fatalf("Failed to open git repository: %v", err)
		return
	}
	c, err := r.Commit("test")
	if err != nil {
		t.Fatalf("Failed to commit changes: %v", err)
		return
	}
	if c {
		err = r.Push()
		if err != nil {
			t.Fatalf("Failed to push changes: %v", err)
			return
		}
	}
}

func TestRevertFile(t *testing.T) {
	config.SetPath("../../config.yaml")
	r := NewGitRepo(config.GetConfig().Repos[0])
	err := r.Open(false)
	if err != nil {
		t.Fatalf("Failed to open git repository: %v", err)
		return
	}
	h := "50f2c8891ad7d9cc0af6690ae0539aab160b99be"
	files := []string{"."}
	err = r.RevertFile(h, files)
	if err != nil {
		t.Fatalf("Failed to revert file: %v", err)
		return
	}
}
