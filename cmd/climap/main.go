package main

import (
	"climap/internal/bench"
	"climap/internal/creds"
	"climap/internal/sentences"
	"github.com/alexflint/go-arg"
	"github.com/fpawel/errorx"
	"log/slog"
	"sync/atomic"
)

type Args struct {
	ImapAddr          string `arg:"--addr,env:IMAP_ADDR" help:"IMAP server address to dial."`
	ImapCreds         string `arg:"-a,env:IMAP_CREDS" help:"IMAP server credentials in the form of login and password separated by a comma."`
	Connections       int    `arg:"--connections,env:CONNECTIONS" help:"number of connections"`
	SentencesFilePath string `arg:"--sentences,env:SENTENCES_FILE_PATH" default:"sentences.txt" help:"path to sentences file for random test mails"`
}

func main() {
	var args Args
	arg.MustParse(&args)

	mails, err := sentences.NewFileSentences(args.SentencesFilePath)
	check(err, "read sentences file")

	var connections atomic.Int64
	credentials, err := creds.Parse(args.ImapCreds)
	check(err, "parse credentials")

	for N := 0; N < args.Connections; N++ {
		b := bench.NewBuilder(args.ImapAddr, N, credentials, mails, &connections)

		go func() {
			for {
				if err := b.Do(); err != nil {
					slog.Error("failed", errorx.Attr(err), "N", b.N, "connections", connections.Load())
				}
			}
		}()
	}
	<-make(chan bool)
}
