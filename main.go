package main

import (
	"embed"

	"github.com/charghet/go-sync/internal/run"
	"github.com/charghet/go-sync/internal/web"
)

//go:embed frontend/dist
var fs embed.FS

func main() {
	dist := web.EmbedFS{
		FS:     fs,
		Prefix: "frontend/dist",
	}
	run.GetRunner().Run()
	web.StartUp(dist)
}
