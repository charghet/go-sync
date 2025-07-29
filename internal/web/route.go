package web

import (
	"embed"
	"io/fs"
	"net/http"
	"sync"

	"github.com/charghet/go-sync/internal/web/controller"
	"github.com/charghet/go-sync/pkg/logger"
	"github.com/charghet/go-sync/pkg/web"
	"github.com/gin-gonic/gin"
)

var once sync.Once

type EmbedFS struct {
	FS     embed.FS
	Prefix string
}

func SetupRoute(dist EmbedFS) (router *gin.Engine) {
	once.Do(func() {
		router = gin.Default()
		router.Use(web.ErrorHandler)
		router.Use(web.CookieHandler)
		RegisterWebRoutes(router, dist)
	})
	return router
}

func RegisterWebRoutes(router *gin.Engine, dist EmbedFS) {
	prefix := "/api"
	c := controller.NewMainController()
	router.POST(prefix+"/login", c.Login)
	router.POST(prefix+"/repos", c.Repos)
	router.POST(prefix+"/commits", c.Commits)
	router.POST(prefix+"/revert", c.Revert)

	sfs, err := fs.Sub(dist.FS, dist.Prefix)
	if err != nil {
		logger.Fatal("Failed to get static file:", err)
	}
	router.StaticFS("/", http.FS(sfs))
	router.NoRoute(func(ctx *gin.Context) {
		index, err := dist.FS.ReadFile(dist.Prefix + "/index.html")
		if err != nil {
			ctx.String(404, "index.html not found")
			return
		}
		ctx.Data(200, "text/html; charset=utf-8", index)
	})
}
