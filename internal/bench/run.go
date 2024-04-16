package bench

import (
	"github.com/fpawel/errorx"
	"log/slog"
	"sync/atomic"
)

func Run(imapAddr string, connectionsCount int, creds CredentialsGetter, mailsProvider MailProvider) {
	var cons atomic.Int64
	for N := 0; N < connectionsCount; N++ {
		b := builder{
			addr:         imapAddr,
			N:            N,
			creds:        creds,
			cons:         connections{Int64: &cons},
			MailProvider: mailsProvider,
		}
		go func() {
			for {
				if err := b.Do(); err != nil {
					slog.Error("failed", errorx.Attr(err), "N", b.N, "connections", cons.Load())
				}
			}
		}()
	}
}
