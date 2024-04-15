package main

import (
	"climap/pkg/slog/attr"
	"github.com/alexflint/go-arg"
	"log/slog"
	"strings"
	"sync/atomic"
)

type (
	Args struct {
		ImapAddr          string `arg:"--addr,env:IMAP_ADDR" help:"IMAP server address to dial."`
		ImapCreds         string `arg:"-a,env:IMAP_CREDS" help:"IMAP server credentials in the form of login and password separated by a comma."`
		Connections       int    `arg:"--connections,env:CONNECTIONS" help:"number of connections"`
		SentencesFilePath string `arg:"--sentences,env:SENTENCES_FILE_PATH" default:"sentences.txt" help:"path to sentences file for random test mails"`
	}
)

var connections int64

func main() {
	var args Args
	arg.MustParse(&args)
	xs := strings.Split(args.ImapCreds, ",")

	sentences = mustReadVSentencesFromFile(args.SentencesFilePath)

	for N := 0; N < args.Connections; N++ {
		b := imapBenchmarkBuilder{
			Addr: args.ImapAddr,
			N:    N,
			Creds: creds{
				login:    xs[0],
				password: xs[1],
			},
		}
		go func() {
			for {
				log := slog.Default().With("N", b.N)
				if err := b.doBenchmark(); err != nil {
					log.Error("failed", attr.Err(err), "connections", atomic.LoadInt64(&connections))
				}
			}
		}()
	}
	<-make(chan bool)
}
