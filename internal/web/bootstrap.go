package web

import (
	"fmt"

	"github.com/charghet/go-sync/pkg/logger"
)

func StartUp() {
	r := SetupRoute()
	err := r.Run(fmt.Sprintf("%s:%d", "127.0.0.1", 1212))
	if err != nil {
		logger.Fatal(err)
	}
}
