package main

import (
	"climap/pkg/slog/attr"
	"log/slog"
	"os"
)

func check(err error, a ...interface{}) {
	if err == nil {
		return
	}

	switch len(a) {
	case 0:
		slog.Error(err.Error(), attr.Err(err))
	case 1:
		slog.Error(a[0].(string), attr.Err(err))
	default:
		msg := a[0].(string)
		a = append([]interface{}{attr.Err(err)}, a[1:]...)
		slog.Error(msg, a...)
	}
	os.Exit(1)
}
