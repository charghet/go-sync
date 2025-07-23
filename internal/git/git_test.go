package git

import (
	"testing"

	"github.com/charghet/go-sync/internal/config"
)

func TestInit(t *testing.T) {
	config.SetPath("../../config.yaml")
	r := NewGitRepo(config.GetConfig().Repos[0])
	err := r.Open()
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
