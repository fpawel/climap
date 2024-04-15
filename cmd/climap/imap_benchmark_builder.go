package main

import (
	"climap/pkg/errorx"
	"climap/pkg/slog/attr"
	"crypto/tls"
	"github.com/emersion/go-imap/v2/imapclient"
	"log/slog"
	"time"
)

type (
	imapBenchmarkBuilder struct {
		Addr  string
		N     int
		Creds CredentialsGetter
	}

	CredentialsGetter interface {
		GetCredentials(N int) (string, string, error)
	}
)

func (x imapBenchmarkBuilder) doBenchmark() error {
	c, err := x.newBenchmark()
	if err != nil {
		return errorx.Wrap(err)
	}
	defer func() {
		if err := c.Close(); err != nil {
			c.log.Error("close", "session-id", c.s.sessionID, attr.Err(err))
		}
	}()
	if err := c.do(); err != nil {
		return errorx.Args("N", x.N, "session-id", c.s.sessionID).Wrap(err)
	}
	return nil
}

func (x imapBenchmarkBuilder) newBenchmark() (r imapBenchmark, _ error) {
	var sesID sessionSniffer
	sesID.wg.Add(1)
	tm := time.Now()
	imapClient, err := imapclient.DialTLS(x.Addr, &imapclient.Options{
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DebugWriter: &sesID,
	})
	if err != nil {
		return r, errorx.Prepend("failed to dial IMAP server").Args("since", time.Since(tm).String()).Wrap(err)
	}
	sesID.wg.Wait()
	return imapBenchmark{
		s:   &sesID,
		b:   x,
		c:   imapClient,
		log: slog.Default().With("session-id", sesID.sessionID, "N", x.N),
	}, nil
}