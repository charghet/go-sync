package main

import (
	"github.com/charghet/go-sync/internal/run"
	"github.com/charghet/go-sync/internal/web"
)

func main() {
	run.GetRunner().Run()
	web.StartUp()
}
