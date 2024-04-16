package main

import (
	"github.com/fpawel/slogx/pretty"
	"github.com/fpawel/slogx/slogctx"
	"log/slog"
)

func init() {
	slog.SetDefault(
		slog.New(
			slogctx.NewHandler(
				pretty.NewHandler().
					WithAddSource(false).
					WithLevel(slog.LevelInfo))))
}
