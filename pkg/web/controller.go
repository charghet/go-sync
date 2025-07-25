package web

import (
	"github.com/charghet/go-sync/pkg/util"
	"github.com/gin-gonic/gin"
)

type BaseController struct {
}

type Pager struct {
	Index int `json:"index"`
	Size int `json:"size"`
}

func (*BaseController) ResponseJson(ctx *gin.Context, code int, msg string, data interface{}) {
	ctx.JSON(code, &Result{
		Code: code,
		Msg:  msg,
		Data: data,
	})
	ctx.Abort()
}

func (c *BaseController) ResponseOkJson(ctx *gin.Context, data interface{}) {
	c.ResponseJson(ctx, 200, "ok", data)
}

func (c *BaseController) BindJSON(ctx *gin.Context, req interface{}) {
	err := ctx.ShouldBindJSON(req)
	CheckErr(err, ServiceErr{Code: 300, Msg: "request json parse error"})
}

func (c *BaseController) BindParam(ctx *gin.Context, req interface{}) {
	err := ctx.ShouldBindQuery(req)
	CheckErr(err, ServiceErr{Code: 300, Msg: "request param parse error"})
}

func (c *BaseController) SetLogin(ctx *gin.Context, name string) string {
	token, _ := util.GenerateJWT(name)
	ctx.SetCookie("token", token, 3600, "/", "", false, true)
	return token
}

func (c *BaseController) GetLogin(ctx *gin.Context) string {
	token, _ := ctx.Cookie("token")
	name, _ := util.ValidateJWT(token)
	return name
}

