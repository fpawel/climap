package main

import (
	"climap/internal/bench"
	"github.com/alexflint/go-arg"
	"github.com/fpawel/errorx"
	"log/slog"
	"strings"
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
	xs := strings.Split(args.ImapCreds, ",")

	sentences := bench.mustNewtestMailProviderFromSentencesFile(args.SentencesFilePath)

	var connections atomic.Int64

	for N := 0; N < args.Connections; N++ {
		b := imapBenchmarkBuilder{
			Addr: args.ImapAddr,
			N:    N,
			Creds: creds{
				login:    xs[0],
				password: xs[1],
			},
			Cons:             bench.connections{Int64: &connections},
			TestMailProvider: sentences,
		}
		go func() {
			for {
				if err := b.doBenchmark(); err != nil {
					slog.Error("failed", errorx.Attr(err), "N", b.N, "connections", connections.Load())
				}
			}
		}()
	}
	<-make(chan bool)
}
