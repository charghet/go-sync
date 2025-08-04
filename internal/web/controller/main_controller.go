package controller

import (
	"github.com/charghet/go-sync/internal/config"
	"github.com/charghet/go-sync/internal/git"
	"github.com/charghet/go-sync/internal/run"
	"github.com/charghet/go-sync/pkg/web"
	"github.com/gin-gonic/gin"
)

type MainController struct {
	web.BaseController
}

func NewMainController() *MainController {
	return &MainController{}
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var conf = config.GetConfig()

func (c *MainController) Login(ctx *gin.Context) {
	var req LoginReq
	c.BindJSON(ctx, &req)
	if conf.User.Username == req.Username && conf.User.Password == req.Password {
		token := c.SetLogin(ctx, req.Username)
		c.ResponseOkJson(ctx, token)
		return
	}
	panic(web.ServiceErr{Code: 400, Msg: "username or password is incorrect"})
}

type RepoIdReq struct {
	Id int `json:"id"`
}

func (c *MainController) Repos(ctx *gin.Context) {
	c.ResponseOkJson(ctx, config.RepoInfo())
}

type CommitsReq struct {
	RepoIdReq
	Pager web.Pager `json:"pager"`
}

func (c *MainController) Commits(ctx *gin.Context) {
	var req CommitsReq
	c.BindJSON(ctx, &req)
	r := getRepo(req.Id)
	commits, total, err := r.GetCommit(req.Pager.Index, req.Pager.Size)
	web.CheckInnerErr(err, "can not get commits")
	c.ResponseOkJson(ctx, struct {
		Total int          `json:"total"`
		List  []git.Commit `json:"list"`
	}{
		Total: total,
		List:  commits,
	})
}

type RevertReq struct {
	RepoIdReq
	Hash string   `json:"hash"`
	File []string `json:"file"`
}

func (c *MainController) Revert(ctx *gin.Context) {
	var req RevertReq
	c.BindJSON(ctx, &req)
	r := getRepo(req.Id)
	if len(req.File) == 0 {
		req.File = []string{"."}
	}
	err := r.RevertFile(req.Hash, req.File)
	web.CheckServiceErr(err, "")
	run.GetRunner().Ignore(req.Id)
	c.ResponseOkJson(ctx, "ok")
}

type ChangesReq struct {
	RepoIdReq
	Hash string `json:"hash"`
}

func (c *MainController) Changes(ctx *gin.Context) {
	var req ChangesReq
	c.BindJSON(ctx, &req)
	r := getRepo(req.Id)
	changes, err := r.GetChange(req.Hash)
	web.CheckServiceErr(err, "")
	c.ResponseOkJson(ctx, changes)
}

func getRepo(id int) *git.GitRepo {
	if id <= 0 || id > len(run.GetRunner().Repos) {
		panic(web.ServiceErr{Code: 300, Msg: "id not found"})
	}
	return run.GetRunner().Repos[id-1]
}
