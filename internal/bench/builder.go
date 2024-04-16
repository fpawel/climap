package bench

import (
	"climap/internal/session"
	"crypto/tls"
	"github.com/fpawel/errorx"
	"log/slog"
	"sync/atomic"
	"time"
)

type (
	Builder struct {
		addr  string
		N     int
		creds CredentialsGetter
		cons  connections
		MailProvider
	}

	MailProvider interface {
		NewMail(from string) string
	}

	CredentialsGetter interface {
		GetCredentials(N int) (string, string, error)
	}
)

func NewBuilder(addr string, N int, creds CredentialsGetter, m MailProvider, cons *atomic.Int64) Builder {
	return Builder{
		addr:         addr,
		N:            N,
		creds:        creds,
		MailProvider: m,
		cons:         connections{cons},
	}
}

func (x Builder) Do() error {
	c, err := x.new()
	if err != nil {
		return errorx.Wrap(err)
	}
	defer func() {
		if err := c.Close(); err != nil {
			c.log.Error("close", "session-id", c.s.SessionID(), errorx.Attr(err))
		}
	}()
	if err := c.do(); err != nil {
		return errorx.Args("N", x.N, "session-id", c.s.SessionID()).Wrap(err)
	}
	return nil
}

func (x Builder) new() (r benchmark, _ error) {
	tm := time.Now()
	imapClient, ses, err := session.NewSessionClient(x.addr, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return r, errorx.Prepend("failed to dial IMAP server").Args("since", time.Since(tm).String()).Wrap(err)
	}
	return benchmark{
		s:   ses,
		b:   x,
		c:   imapClient,
		log: slog.Default().With("session-id", ses.SessionID(), "N", x.N),
	}, nil
}
