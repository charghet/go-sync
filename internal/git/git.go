package git

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/charghet/go-sync/internal/config"
	"github.com/charghet/go-sync/pkg/logger"
	"github.com/charghet/go-sync/pkg/util"
	"github.com/go-git/go-git/v6"
	gitConfig "github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
	"github.com/go-git/go-git/v6/utils/merkletrie"
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

func (r *GitRepo) Open(pull bool) error {
	var err error
	r.repo, err = git.PlainOpen(r.RepoConfig.Path)

	if err != nil {
		if err == git.ErrRepositoryNotExists {
			logger.Info(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Repository does not exist, initing:", r.RepoConfig.Url)
			r.repo, err = git.PlainInit(r.RepoConfig.Path, false)
			if err != nil {
				logger.Fatal(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Failed to init git repository:", r.RepoConfig.Path, "Error:", err)
				return err
			}
			_, err = r.repo.CreateRemote(&gitConfig.RemoteConfig{
				Name: "origin",
				URLs: []string{r.RepoConfig.Url},
			})
			if err != nil {
				logger.Fatal(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Failed to create remote repository:", err)
				return err
			}
			logger.Info(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Created remote repository 'origin' for:", r.RepoConfig.Path)

			err = r.repo.CreateBranch(&gitConfig.Branch{
				Name:   r.RepoConfig.Branch,
				Remote: "origin",
				Merge:  plumbing.NewBranchReferenceName(r.RepoConfig.Branch),
			})
			if err != nil {
				logger.Fatal(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Failed to create branch:", r.RepoConfig.Branch, "Error:", err)
				return err
			}
		}
	}
	if err != nil {
		logger.Fatal(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Failed to open git repository:", err)
		return err
	}

	r.worktree, err = r.repo.Worktree()
	if err != nil {
		logger.Fatal(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Failed to get worktree:", err)
		return err
	}

	if pull {
		err = r.Pull()
		if err != nil {
			logger.Fatal(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Failed to pull changes after init:", err)
		}
		c, err := r.Commit("auto commit by init in " + time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			logger.Fatal(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Failed to commit after init:", err)
			return err
		}

		if c {
			err = r.Push()
			if err != nil {
				logger.Fatal(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Failed to push after init:", err)
				return err
			}
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
		logger.Danger(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Failed to clone repository:", r.RepoConfig.Url)
		return err
	}
	return nil
}

func (r *GitRepo) Commit(message string) (commit bool, err error) {
	_, err = r.worktree.Add(".")
	if err != nil {
		logger.Danger(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Failed to add changes to worktree:", err)
	}
	status, err := r.worktree.Status()
	if err != nil {
		logger.Danger(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Failed to get worktree status:", err)
		return false, err
	}
	if status.IsClean() {
		logger.Info(fmt.Sprintf("[%v]", r.RepoConfig.Name), "No changes to commit, worktree is clean.")
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
		logger.Danger(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Failed to commit changes:", err)
		return false, err
	}
	logger.Info(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Committed changes:", h.String(), message)
	return commit, nil
}

func (r *GitRepo) Push() error {
	err := r.repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       r.Auth,
	})
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			logger.Info(fmt.Sprintf("[%v]", r.RepoConfig.Name), "No changes to push, repository is up to date.")
			return nil
		}
		logger.Danger(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Failed to push changes:", err)
		return err
	}
	logger.Info(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Pushed changes to remote repository: ", r.RepoConfig.Url)
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
			logger.Info(fmt.Sprintf("[%v]", r.RepoConfig.Name), "No changes to pull, repository is up to date.")
			return nil
		}
		logger.Danger(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Failed to pull changes:", err)
		return err
	}
	return nil
}

func (r *GitRepo) Checkout(hash string, files []string) error {
	err := r.worktree.Checkout(&git.CheckoutOptions{
		Hash:                      plumbing.NewHash(hash),
		SparseCheckoutDirectories: files,
	})
	if err != nil {
		logger.Danger(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Failed to checkout:", hash, "Error:", err)
		return err
	}
	logger.Info(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Checked out:", hash, "with files:", files)
	return nil
}

func (r *GitRepo) Restore(files []string) error {
	err := r.worktree.Restore(&git.RestoreOptions{
		Files: files,
	})

	if err != nil {
		logger.Danger(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Failed to restore files:", files, "Error:", err)
		return err
	}
	logger.Info(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Restored files:", files)
	return nil
}

func (r *GitRepo) Reset(hash string, files []string) error {
	err := r.worktree.Reset(&git.ResetOptions{
		Commit: plumbing.NewHash(hash),
		Files:  files,
	})
	if err != nil {
		logger.Danger(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Failed to reset to hash:", hash, "Error:", err)
		return err
	}
	logger.Info(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Reset worktree to hash:", hash)
	return nil
}

func (r *GitRepo) RevertFile(hash string, files []string) error {
	until := time.Now()
	cIter, err := r.repo.Log(&git.LogOptions{Until: &until})
	if err != nil {
		logger.Danger(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Failed to get commit iterator:", err)
		return err
	}

	rct := plumbing.NewHash(hash)

	foundHash := false
	err = cIter.ForEach(func(c *object.Commit) error {
		if c.Hash == rct {
			foundHash = true
			fi, err := c.Files()
			if err != nil {
				logger.Danger(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Failed to get files from commit:", c.Hash, "Error:", err)
				return err
			}
			fileSet := util.SliceToSet(files)
			_, all := fileSet["."]
			err = fi.ForEach(func(cf *object.File) error {
				if _, e := fileSet[cf.Name]; all || e {
					delete(fileSet, cf.Name)
					fr, err := cf.Blob.Reader()
					if err != nil {
						logger.Danger(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Failed to get file reader for:", cf.Name, "Error:", err)
						return err
					}
					fw, err := os.Create(filepath.Join(r.RepoConfig.Path, cf.Name))
					if err != nil {
						logger.Danger(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Failed to create file:", cf.Name, "Error:", err)
						return err
					}
					io.Copy(fw, fr)
					logger.Info(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Reverted file:", cf.Name, "to commit:", c.Hash)
				}
				return nil
			})
			if err != nil {
				return err
			}
			if !all && len(fileSet) > 0 {
				s := fmt.Sprintf("Some files were not found in commit:%v Files not found:%v", c.Hash, fileSet)
				logger.Warn(fmt.Sprintf("[%v]", r.RepoConfig.Name), s)
				return errors.New(s)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	if !foundHash {
		s := fmt.Sprintf("Commit hash not found:%v", hash)
		logger.Warn(fmt.Sprintf("[%v]", r.RepoConfig.Name), s)
		return errors.New(s)
	}
	return nil
}

type Commit struct {
	Hash    string `json:"hash"`
	Message string `json:"message"`
	Author  string `json:"author"`
	Date    string `json:"date"`
	Email   string `json:"email"`
}

func (r *GitRepo) GetCommit(pageIndex, pageSize int) (commits []Commit, total int, err error) {
	until := time.Now()
	cIter, err := r.repo.Log(&git.LogOptions{Until: &until})
	if err != nil {
		logger.Danger(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Failed to get commit iterator:", err)
		return nil, 0, err
	}

	start := (pageIndex - 1) * pageSize
	end := pageIndex * pageSize
	cIter.ForEach(func(c *object.Commit) error {
		total += 1
		if pageIndex == 0 || (total > start && total <= end) {
			commits = append(commits, Commit{
				Hash:    c.Hash.String(),
				Message: c.Message,
				Author:  c.Author.Name,
				Date:    c.Author.When.Format("2006-01-02 15:04:05"),
				Email:   c.Author.Email,
			})
		}
		return nil
	})
	return commits, total, nil
}

type Change struct {
	Action string `json:"action"`
	Name   string `json:"name"`
}

func (r *GitRepo) GetChange(hash string) ([]Change, error) {
	commitHash := plumbing.NewHash(hash)
	commit, err := r.repo.CommitObject(commitHash)
	if err != nil {
		logger.Danger(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Failed to get commit:", err)
		return nil, err
	}

	parentCommit, err := commit.Parents().Next()
	if err != nil {
		if err != io.EOF {
			logger.Danger(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Failed to get parent commit:", err)
			return nil, err
		}
	}

	var commitTree, parentTree *object.Tree
	if parentCommit == nil {
		parentTree = &object.Tree{}
	} else {
		parentTree, err = parentCommit.Tree()
		if err != nil {
			logger.Danger(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Failed to get parent commit tree:", err)
			return nil, err
		}
	}
	commitTree, err = commit.Tree()
	if err != nil {
		logger.Danger(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Failed to get commit tree:", err)
		return nil, err
	}

	changes, err := parentTree.Diff(commitTree)
	if err != nil {
		logger.Danger(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Failed to get changes:", err)
		return nil, err
	}
	res := make([]Change, changes.Len())
	for i, change := range changes {
		action, err := change.Action()
		if err != nil {
			logger.Danger(fmt.Sprintf("[%v]", r.RepoConfig.Name), "Failed to get action:", err)
			return nil, err
		}
		var c Change
		switch action {
		case merkletrie.Insert:
			c.Action = "A"
			c.Name = change.To.Name
		case merkletrie.Modify:
			c.Action = "M"
			c.Name = change.To.Name
		case merkletrie.Delete:
			c.Action = "D"
			c.Name = change.From.Name
		}
		res[i] = c
	}
	return res, nil
}
