package web

import (
	"fmt"
	"runtime/debug"

	"github.com/charghet/go-sync/pkg/logger"
	"github.com/charghet/go-sync/pkg/util"
	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				code := 500
				msg := "internal server error"

				if serr, ok := err.(ServiceErr); ok {
					msg = serr.Msg
					code = 400
					if serr.Code != 0 {
						code = serr.Code
					}
					logger.Warn(fmt.Sprintf("service error: %s", serr.Msg))
				} else {
					code = 501
					logger.Danger(fmt.Sprintf("inner error: %s", err.Error()))
					if ierr, ok := err.(InnerErr); ok {
						if ierr.Code != 0 {
							code = ierr.Code
						}
					}
					debug.PrintStack()
				}
				c.JSON(200, &Result{
					Code: code,
					Msg:  msg,
				})
				c.Abort()
			} else {
				logger.Danger(fmt.Sprintf("panic: %v", r))
				debug.PrintStack()
				c.JSON(200, &Result{
					Code: 500,
					Msg:  "internal server error",
				})
				c.Abort()
			}
		}
	}()

	c.Next()
}

func CookieHandler(c *gin.Context) {
	if c.Request.URL.Path == "/api/login" {
		c.Next()
		return
	}
	token, err := c.Cookie("token")
	if err == nil {
		_, err = util.ValidateJWT(token)
	}
	if err != nil {
		c.JSON(401, Result{
			Code: 401,
			Msg:  "auth failed",
			Data: c.Request.Host,
		})
		c.Abort()
		return
	}
	c.Next()
}
