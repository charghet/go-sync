package run

import (
	"fmt"
	"os"
	"time"

	"github.com/charghet/go-sync/internal/config"
	"github.com/charghet/go-sync/internal/git"
	"github.com/charghet/go-sync/internal/notify"
	"github.com/charghet/go-sync/pkg/logger"
)

type Runner struct {
	Repos        []*git.GitRepo
	ignoreTimers []*time.Timer
}

var runner *Runner

func GetRunner() *Runner {
	if runner == nil {
		runner = &Runner{
			Repos:        make([]*git.GitRepo, len(config.GetConfig().Repos)),
			ignoreTimers: make([]*time.Timer, len(config.GetConfig().Repos)),
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
		err := repo.Open(*repoConfig.Pull)
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

		r.ignoreTimers[i] = time.NewTimer(100 * time.Millisecond)
		go func() {
			timer := time.NewTimer(50 * time.Millisecond)
			defer timer.Stop()
			<-timer.C
			ignoreTimer := r.ignoreTimers[i]
			for {
				select {
				case event, ok := <-n.Events:
					if !ok {
						return
					}

					timer.Stop()
					timer.Reset(time.Duration(*repo.RepoConfig.Debounce) * time.Second)
					logger.Info(fmt.Sprintf("[%v]", repo.RepoConfig.Name), "Received event:", event, "for path:", event.Name)
				case <-timer.C:
					select {
					case <-ignoreTimer.C:
						logger.Info(fmt.Sprintf("[%v]", repo.RepoConfig.Name), "Timer expired, committing changes.")
						c, err := repo.Commit("auto commit in " + time.Now().Format("2006-01-02 15:04:05"))
						if err != nil {
							logger.Warn(fmt.Sprintf("[%v]", repo.RepoConfig.Name), "Failed to commit changes:", err)
						}
						if c {
							repo.Push()
						}
						ignoreTimer.Reset(100 * time.Millisecond)
					default:
						logger.Debug(fmt.Sprintf("[%v]", repo.RepoConfig.Name), "ignoreTimer not stop, skip..")
					}

				case err, ok := <-n.Errors:
					if !ok {
						return
					}
					logger.Danger(fmt.Sprintf("[%v]", repo.RepoConfig.Name), "Error:", err)
				}
			}
		}()
	}
}

func (r *Runner) Ignore(id int) {
	i := id - 1
	r.ignoreTimers[i].Stop()
	r.ignoreTimers[i].Reset(time.Duration(*config.GetConfig().Repos[i].Ignore) * time.Second)
}
