package main

import (
	"climap/pkg/errorx"
	"climap/pkg/slog/attr"
	"errors"
	"github.com/elliotchance/pie/v2"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/sourcegraph/conc/pool"
	"log/slog"
	"time"
)

type imapBenchmark struct {
	s   *sessionSniffer
	b   imapBenchmarkBuilder
	c   *imapclient.Client
	log *slog.Logger
}

func (x imapBenchmark) Close() error {
	return errors.Join(x.s.fd.Close(), x.c.Close())
}

func (x imapBenchmark) do() error {
	login, pass, err := x.b.Creds.GetCredentials(x.b.N)
	if err != nil {
		return errorx.Prepend("get credentials").Wrap(err)
	}

	if err := x.Login(login, pass); err != nil {
		return errorx.Prepend("login").Args("login", login).Wrap(err)
	}

	x.log.Info("in", "connections", x.b.Cons.Inc())
	defer x.b.Cons.Dec()

	iteration := 0
	for {
		tm := time.Now()
		p := pool.New().WithErrors().WithFirstError()
		for i := 0; i < 16; i++ {
			p.Go(func() error {
				return x.doMailbox("Drafts")
			})
		}
		if err := p.Wait(); err != nil {
			return err
		}
		iteration++
		x.log.Info("done", attr.Since(tm), "iteration", iteration, "connections", x.b.Cons.Get())
	}
}

func (x imapBenchmark) Login(username, password string) error {
	if err := x.c.Login(username, password).Wait(); err != nil {
		return errorx.Prepend("failed to login").Wrap(err)
	}
	return nil
}

func (x imapBenchmark) Select(mailbox string) (*imap.SelectData, error) {
	r, err := x.c.Select(mailbox, nil).Wait()
	if err != nil {
		return nil, errorx.Prepend("failed to select").Args("mailbox", mailbox).Wrap(err)
	}
	return r, nil
}

func (x imapBenchmark) AppendRandomMails(mailbox string, n int) ([]imap.AppendData, error) {
	p := pool.NewWithResults[imap.AppendData]().WithErrors().WithFirstError().WithMaxGoroutines(16)
	for i := 0; i < n; i++ {
		p.Go(func() (imap.AppendData, error) {
			return x.AppendRandomMail(mailbox)
		})
	}
	return p.Wait()
}

func (x imapBenchmark) AppendRandomMail(mailbox string) (imap.AppendData, error) {
	e := errorx.Args("mailbox", mailbox)
	buf := []byte(x.b.NewTestMail())
	size := int64(len(buf))
	rc := x.c.Append(mailbox, size, nil)
	_, err := rc.Write(buf)
	if err != nil {
		return imap.AppendData{}, e.Prepend("failed to write message").Wrap(err)
	}

	if err := rc.Close(); err != nil {
		return imap.AppendData{}, e.Prepend("failed to close message").Wrap(err)
	}

	r, err := rc.Wait()
	if err != nil {
		return imap.AppendData{}, e.Prepend("APPEND command failed").Wrap(err)
	}
	return *r, nil
}

func (x imapBenchmark) doMailbox(mailbox string) error {
	mailboxes, err := x.c.List("", "%", nil).Collect()
	if err != nil {
		return errorx.Prependf("failed to list mailboxes").Wrap(err)
	}
	mailboxesNames := pie.Map(mailboxes, func(m *imap.ListData) string { return m.Mailbox })

	inbox, err := x.Select(mailbox)
	if err != nil {
		return errorx.Args("mailboxes", mailboxesNames, "mailbox", mailbox).Prepend("select").Wrap(err)
	}

	const minCount = 100
	if inbox.NumMessages < minCount {
		if _, err = x.AppendRandomMails(mailbox, minCount); err != nil {
			return errorx.Prepend("append random mails").Wrap(err)
		}
		inbox, err = x.Select(mailbox)
		if err != nil {
			return errorx.Args("mailboxes", mailboxesNames, "mailbox", mailbox).Prepend("select").Wrap(err)
		}
	}

	seqSet := imap.SeqSetNum()
	for i := uint32(1); i <= minCount; i++ {
		seqSet.AddNum(i)
	}
	if _, err = x.c.Fetch(seqSet, &imap.FetchOptions{Envelope: true, UID: true}).Collect(); err != nil {
		return errorx.Prepend("fetch").Wrap(err)
	}
	return nil
}
