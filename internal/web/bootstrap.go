package web

import (
	"fmt"

	"github.com/charghet/go-sync/internal/config"
	"github.com/charghet/go-sync/pkg/logger"
)

func StartUp(dist EmbedFS) {
	r := SetupRoute(dist)
	con := config.GetConfig()
	err := r.Run(fmt.Sprintf("%s:%d", con.Server.Host, con.Server.Port))
	if err != nil {
		logger.Fatal(err)
	}
}
