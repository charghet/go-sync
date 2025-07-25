package web

import (
	"sync"

	"github.com/charghet/go-sync/internal/web/controller"
	"github.com/charghet/go-sync/pkg/web"
	"github.com/gin-gonic/gin"
)

var once sync.Once

func SetupRoute() (router *gin.Engine) {
	once.Do(func() {
		router = gin.Default()
		router.Use(web.ErrorHandler)
		router.Use(web.CookieHandler)
		RegisterWebRoutes(router)
	})
	return router
}

func RegisterWebRoutes(router *gin.Engine) {
	prefix := "/api"
	c := controller.NewMainController()
	router.POST(prefix+"/login", c.Login)
	router.POST(prefix+"/repos", c.Repos)
	router.POST(prefix+"/commits", c.Commits)
	router.POST(prefix+"/revert", c.Revert)
}
