package main

import (
	"climap/pkg/slog/pretty"
	"climap/pkg/slog/slogctx"
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
