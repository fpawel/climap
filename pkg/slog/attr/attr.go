package attr

import (
	"climap/pkg/errorx"
	"log/slog"
	"time"
)

func Err(err error) slog.Attr {
	if e := errorx.Get(err); len(e.Frames) != 0 {
		return slog.Any(err.Error(), e.Details())
	}
	return slog.String("error", err.Error())
}

func Since(tm time.Time) slog.Attr {
	return slog.String("since", time.Since(tm).String())
}
