package run

import (
	"os"
	"time"

	"github.com/charghet/go-sync/internal/config"
	"github.com/charghet/go-sync/internal/git"
	"github.com/charghet/go-sync/internal/notify"
	"github.com/charghet/go-sync/pkg/logger"
)

func Run() error {
	logger.SetLogFile("run.log")
	logger.Info("Starting go-sync...")

	repos := config.GetConfig().Repos
	if len(repos) == 0 {
		logger.Warn("No repositories configured, exiting.")
		os.Exit(0)
	}
	for _, repoConfig := range repos {
		repo := git.NewGitRepo(repoConfig)
		err := repo.Open()
		if err != nil {
			continue
		}

		n, err := notify.NewNotify()
		if err != nil {
			continue
		}
		defer n.Close()

		err = n.Add(repoConfig.Path)
		if err != nil {
			continue
		}

		go func() {
			for {
				select {
				case event, ok := <-n.Events:
					if !ok {
						return
					}
					logger.Info("Received event:", event, "for path:", event.Name)
					repo.Commit("auto commit in " + time.Now().Format("2006-01-02 15:04:05"))
					repo.Push()
				case err, ok := <-n.Errors:
					if !ok {
						return
					}
					logger.Danger("Error:", err)
				}
			}
		}()
	}
	<-make(chan struct{})
	return nil
}
