package main

import (
	"fmt"
	"os"

	"github.com/charghet/go-sync/internal/notify"
	"github.com/charghet/go-sync/pkg/logger"
)

func Main() {
	n, err := notify.NewNotify()
	if err != nil {
		return
	}

	testPath := "."
	err = os.MkdirAll(testPath, 0755)
	if err != nil {
		return
	}
	err = n.Add(testPath)
	if err != nil {
		return
	}

	go func() {
		for {
			select {
			case event, ok := <-n.Events:
				if !ok {
					return
				}
				logger.Info(fmt.Sprintf("Event: %s for path: %s", event.Op, event.Name))
			case err, ok := <-n.Errors:
				if !ok {
					return
				}
				logger.Danger(fmt.Sprintf("Error: %v", err))
			}
		}
	}()
	<-make(chan struct{})
}
