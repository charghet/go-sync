package web

import (
	"github.com/charghet/go-sync/pkg/logger"
)

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func CheckErr(err error, perr error) {
	if err != nil {
		panic(perr)
	}
}

func CheckServiceErr(err error, msg string) {
	if err != nil {
		logger.Warn(err)
		if msg == "" {
			msg = err.Error()
		}
		CheckErr(err, ServiceErr{Msg: msg})
	}
}

func CheckInnerErr(err error, msg string) {
	if err != nil {
		logger.Danger(err)
		if msg == "" {
			msg = err.Error()
		}
		CheckErr(err, InnerErr{Msg: msg})
	}
}
