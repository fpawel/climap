package bench

import (
	"crypto/tls"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/fpawel/errorx"
	"log/slog"
	"time"
)

type (
	builder struct {
		addr  string
		N     int
		creds CredentialsGetter
		cons  connections
		MailProvider
	}

	MailProvider interface {
		NewMail() string
	}

	CredentialsGetter interface {
		GetCredentials(N int) (string, string, error)
	}
)

func (x builder) Do() error {
	c, err := x.new()
	if err != nil {
		return errorx.Wrap(err)
	}
	defer func() {
		if err := c.Close(); err != nil {
			c.log.Error("close", "session-id", c.s.sessionID, errorx.Attr(err))
		}
	}()
	if err := c.do(); err != nil {
		return errorx.Args("N", x.N, "session-id", c.s.sessionID).Wrap(err)
	}
	return nil
}

func (x builder) new() (r benchmark, _ error) {
	var sesID sessionSniffer
	sesID.wg.Add(1)
	tm := time.Now()
	imapClient, err := imapclient.DialTLS(x.addr, &imapclient.Options{
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DebugWriter: &sesID,
	})
	if err != nil {
		return r, errorx.Prepend("failed to dial IMAP server").Args("since", time.Since(tm).String()).Wrap(err)
	}
	sesID.wg.Wait()
	return benchmark{
		s:   &sesID,
		b:   x,
		c:   imapClient,
		log: slog.Default().With("session-id", sesID.sessionID, "N", x.N),
	}, nil
}
