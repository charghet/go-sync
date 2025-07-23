package git

import (
	"time"

	"github.com/charghet/go-sync/internal/config"
	"github.com/charghet/go-sync/pkg/logger"
	"github.com/go-git/go-git/v6"
	gitConfig "github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
)

type GitRepo struct {
	RepoConfig   config.RepoConfig
	CloneOption  *git.CloneOptions
	PushOption   *git.PushOptions
	CommitOption *git.CommitOptions
	PullOption   *git.PullOptions
	FetchOption  *git.FetchOptions
	repo         *git.Repository
	worktree     *git.Worktree
}

func NewGitRepo(repoConfig config.RepoConfig) *GitRepo {
	return &GitRepo{
		RepoConfig: repoConfig,
		CloneOption: &git.CloneOptions{
			URL: repoConfig.Url,
			Auth: &http.BasicAuth{
				Username: repoConfig.Username,
				Password: repoConfig.Password,
			},
		},
		PushOption: &git.PushOptions{
			Auth: &http.BasicAuth{
				Username: repoConfig.Username,
				Password: repoConfig.Password,
			},
		},
		CommitOption: &git.CommitOptions{
			Author: &object.Signature{
				Name:  repoConfig.Username,
				Email: repoConfig.Email,
			},
		},
		PullOption: &git.PullOptions{
			RemoteName: "origin",
			// ReferenceName: plumbing.ReferenceName(repoConfig.Branch),
			Auth: &http.BasicAuth{
				Username: repoConfig.Username,
				Password: repoConfig.Password,
			},
		},
		FetchOption: &git.FetchOptions{
			RemoteName: "origin",
			Auth: &http.BasicAuth{
				Username: repoConfig.Username,
				Password: repoConfig.Password,
			},
		},
	}
}

func (r *GitRepo) Open() error {
	var err error
	init := false
	r.repo, err = git.PlainOpen(r.RepoConfig.Path)

	if err != nil {
		if err == git.ErrRepositoryNotExists {
			logger.Info("Repository does not exist, initing:", r.RepoConfig.Url)
			r.repo, err = git.PlainInit(r.RepoConfig.Path, false)
			if err != nil {
				logger.Fatal("Failed to init git repository:", r.RepoConfig.Path, "Error:", err)
				return err
			}
			init = true
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
	if init {
		_, err = r.repo.CreateRemote(&gitConfig.RemoteConfig{
			Name: "origin",
			URLs: []string{r.RepoConfig.Url},
		})
		if err != nil {
			logger.Fatal("Failed to create remote repository:", err)
			return err
		}
		logger.Info("Created remote repository 'origin' for:", r.RepoConfig.Path)
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
	r.repo, err = git.PlainClone(r.RepoConfig.Path, r.CloneOption)
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
	r.CommitOption.Author.When = time.Now()
	_, err = r.worktree.Commit(message, r.CommitOption)
	if err != nil {
		logger.Danger("Failed to commit changes:", err)
		return false, err
	}
	logger.Info("Committed changes with message:", message)
	return commit, nil
}

func (r *GitRepo) Push() error {
	err := r.repo.Push(r.PushOption)
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
	err := r.worktree.Pull(r.PullOption)
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
