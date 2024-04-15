package main

import (
	"github.com/fpawel/errorx"
	"log/slog"
	"os"
)

func check(err error, a ...interface{}) {
	if err == nil {
		return
	}

	switch len(a) {
	case 0:
		slog.Error(err.Error(), errorx.Attr(err))
	case 1:
		slog.Error(a[0].(string), errorx.Attr(err))
	default:
		msg := a[0].(string)
		a = append([]interface{}{errorx.Attr(err)}, a[1:]...)
		slog.Error(msg, a...)
	}
	os.Exit(1)
}
