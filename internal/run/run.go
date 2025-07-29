package run

import (
	"os"
	"time"

	"github.com/charghet/go-sync/internal/config"
	"github.com/charghet/go-sync/internal/git"
	"github.com/charghet/go-sync/internal/notify"
	"github.com/charghet/go-sync/pkg/logger"
)

type Runner struct {
	Repos       []*git.GitRepo
	ignoreTimer *time.Timer
}

var runner *Runner

func GetRunner() *Runner {
	if runner == nil {
		runner = &Runner{
			Repos:       make([]*git.GitRepo, len(config.GetConfig().Repos)),
			ignoreTimer: time.NewTimer(100 * time.Millisecond),
		}
	}
	return runner
}

func (r *Runner) Run() {
	logger.SetLogFile("run.log")
	logger.Info("Starting go-sync...")

	repos := config.GetConfig().Repos
	if len(repos) == 0 {
		logger.Warn("No repositories configured, exiting.")
		os.Exit(0)
	}
	for i, repoConfig := range repos {
		repo := git.NewGitRepo(repoConfig)
		r.Repos[i] = repo
		err := repo.Open(true)
		if err != nil {
			continue
		}

		n, err := notify.NewNotify()
		if err != nil {
			continue
		}

		err = n.Add(repoConfig.Path)
		if err != nil {
			continue
		}

		go func() {
			timer := time.NewTimer(50 * time.Millisecond)
			defer timer.Stop()
			<-timer.C
			for {
				select {
				case event, ok := <-n.Events:
					if !ok {
						return
					}

					timer.Stop()
					timer.Reset(time.Duration(repo.RepoConfig.Debounce) * time.Second)
					logger.Info("Received event:", event, "for path:", event.Name)
				case <-timer.C:
					select {
					case <-r.ignoreTimer.C:
						logger.Info("Timer expired, committing changes.")
						c, err := repo.Commit("auto commit in " + time.Now().Format("2006-01-02 15:04:05"))
						if err != nil {
							logger.Warn("Failed to commit changes:", err)
						}
						if c {
							repo.Push()
						}
						r.ignoreTimer.Reset(100 * time.Millisecond)
					default:
						logger.Debug("ignoreTimer not stop, skip..")
					}

				case err, ok := <-n.Errors:
					if !ok {
						return
					}
					logger.Danger("Error:", err)
				}
			}
		}()
	}
}

func (r *Runner) Ignore() {
	r.ignoreTimer.Stop()
	r.ignoreTimer.Reset(time.Duration(config.GetConfig().Ignore) * time.Second)
}
