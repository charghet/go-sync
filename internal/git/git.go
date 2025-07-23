package git

import (
	"time"

	"github.com/charghet/go-sync/internal/config"
	"github.com/charghet/go-sync/pkg/logger"
	"github.com/go-git/go-git/v6"
	gitConfig "github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
)

type GitRepo struct {
	RepoConfig config.RepoConfig
	repo       *git.Repository
	worktree   *git.Worktree
	Auth       *http.BasicAuth
}

func NewGitRepo(repoConfig config.RepoConfig) *GitRepo {
	return &GitRepo{
		RepoConfig: repoConfig,
		Auth: &http.BasicAuth{
			Username: repoConfig.Username,
			Password: repoConfig.Password,
		},
	}
}

func (r *GitRepo) Open() error {
	var err error
	r.repo, err = git.PlainOpen(r.RepoConfig.Path)

	if err != nil {
		if err == git.ErrRepositoryNotExists {
			logger.Info("Repository does not exist, initing:", r.RepoConfig.Url)
			r.repo, err = git.PlainInit(r.RepoConfig.Path, false)
			if err != nil {
				logger.Fatal("Failed to init git repository:", r.RepoConfig.Path, "Error:", err)
				return err
			}
			_, err = r.repo.CreateRemote(&gitConfig.RemoteConfig{
				Name: "origin",
				URLs: []string{r.RepoConfig.Url},
			})
			if err != nil {
				logger.Fatal("Failed to create remote repository:", err)
				return err
			}
			logger.Info("Created remote repository 'origin' for:", r.RepoConfig.Path)

			err = r.repo.CreateBranch(&gitConfig.Branch{
				Name:   r.RepoConfig.Branch,
				Remote: "origin",
				Merge:  plumbing.NewBranchReferenceName(r.RepoConfig.Branch),
			})
			if err != nil {
				logger.Fatal("Failed to create branch:", r.RepoConfig.Branch, "Error:", err)
				return err
			}
		}
	}
	if err != nil {
		logger.Fatal("Failed to open git repository:", err)
		return err
	}

	r.worktree, err = r.repo.Worktree()
	if err != nil {
		logger.Fatal("Failed to get worktree:", err)
		return err
	}

	err = r.Pull()
	if err != nil {
		logger.Fatal("Failed to pull changes after init:", err)
	}
	c, err := r.Commit("auto commit by init in " + time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		logger.Fatal("Failed to commit after init:", err)
		return err
	}

	if c {
		err = r.Push()
		if err != nil {
			logger.Fatal("Failed to push after init:", err)
			return err
		}
	}

	return nil
}

func (r *GitRepo) Clone() error {
	var err error
	r.repo, err = git.PlainClone(r.RepoConfig.Path, &git.CloneOptions{
		URL:  r.RepoConfig.Url,
		Auth: r.Auth,
	})
	if err != nil {
		logger.Danger("Failed to clone repository:", r.RepoConfig.Url)
		return err
	}
	return nil
}

func (r *GitRepo) Commit(message string) (commit bool, err error) {
	_, err = r.worktree.Add(".")
	if err != nil {
		logger.Danger("Failed to add changes to worktree:", err)
	}
	status, err := r.worktree.Status()
	if err != nil {
		logger.Danger("Failed to get worktree status:", err)
		return false, err
	}
	if status.IsClean() {
		logger.Info("No changes to commit, worktree is clean.")
		return false, nil
	}
	commit = true
	h, err := r.worktree.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  r.RepoConfig.Username,
			Email: r.RepoConfig.Email,
			When:  time.Now(),
		},
	})

	if err != nil {
		logger.Danger("Failed to commit changes:", err)
		return false, err
	}
	logger.Info("Committed changes:", h.String(), message)
	return commit, nil
}

func (r *GitRepo) Push() error {
	err := r.repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       r.Auth,
	})
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			logger.Info("No changes to push, repository is up to date.")
			return nil
		}
		logger.Danger("Failed to push changes:", err)
		return err
	}
	logger.Info("Pushed changes to remote repository: ", r.RepoConfig.Url)
	return nil
}

func (r *GitRepo) Pull() error {
	err := r.worktree.Pull(&git.PullOptions{
		RemoteName:    "origin",
		ReferenceName: plumbing.NewBranchReferenceName(r.RepoConfig.Branch),
		Auth:          r.Auth,
	})
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			logger.Info("No changes to pull, repository is up to date.")
			return nil
		}
		logger.Danger("Failed to pull changes:", err)
		return err
	}
	return nil
}
