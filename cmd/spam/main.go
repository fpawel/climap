package main

import (
	"climap/internal/creds"
	"climap/internal/sentences"
	"climap/internal/session"
	"crypto/tls"
	"github.com/alexflint/go-arg"
	"github.com/elliotchance/pie/v2"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/fpawel/errorx"
	"github.com/sourcegraph/conc/pool"
	"log/slog"
)

type Args struct {
	ImapAddr          string `arg:"--addr,env:IMAP_ADDR" help:"IMAP server address to dial."`
	ImapCreds         string `arg:"-a,env:IMAP_CREDS" help:"IMAP server credentials in the form of login and password separated by a comma."`
	From              string `arg:"-f,--from,env:MAIL_FROM"`
	Mailbox           string `arg:"-m,--mailbox,env:MAILBOX"`
	Count             int    `arg:"-c,--count,env:MAIL_COUNT"`
	SentencesFilePath string `arg:"--sentences,env:SENTENCES_FILE_PATH" default:"sentences.txt" help:"path to sentences file for random test mails"`
}

func main() {
	var args Args
	arg.MustParse(&args)

	mailBuilder, err := sentences.NewFileSentences(args.SentencesFilePath)
	check(err, "read sentences file")

	credentials, err := creds.Parse(args.ImapCreds)
	check(err, "parse credentials")

	imapClient, ses, err := session.NewSessionClient(args.ImapAddr, &tls.Config{InsecureSkipVerify: true})
	check(err, "failed to dial IMAP server")

	defer func() {
		if err := ses.Close(); err != nil {
			slog.Error("close session", errorx.Attr(err))
		}
	}()

	check(imapClient.Login(credentials.Login, credentials.Password).Wait(), "login")

	mailboxes, err := imapClient.List("", "%", nil).Collect()
	check(err, "failed to list mailboxes")

	mailboxesNames := pie.Map(mailboxes, func(m *imap.ListData) string { return m.Mailbox })

	mailbox, err := imapClient.Select(args.Mailbox, nil).Wait()
	check(err, "select", "mailbox", args.Mailbox, "mailboxes", mailboxesNames)

	slog.Info("mailbox", "mailbox", mailbox)

	p := pool.New().WithErrors().WithFirstError().WithMaxGoroutines(16)
	for i := 0; i < args.Count; i++ {
		mail := mailBuilder.NewMail(args.From)
		p.Go(func() error {
			d, err := appendMail(imapClient, args.Mailbox, []byte(mail))
			if err != nil {
				return err
			}
			slog.Info("sent", "mail", mail, "data", d)
			return nil
		})
	}
	check(p.Wait())
}

func appendMail(imapClient *imapclient.Client, mailbox string, mail []byte) (*imap.AppendData, error) {
	size := int64(len(mail))

	rc := imapClient.Append(mailbox, size, nil)
	_, err := rc.Write(mail)
	if err != nil {
		return nil, errorx.Prepend("failed to write message").Wrap(err)
	}

	if err := rc.Close(); err != nil {
		return nil, errorx.Prepend("failed to close message").Wrap(err)
	}

	d, err := rc.Wait()
	if err != nil {
		return nil, errorx.Prepend("APPEND command failed").Wrap(err)
	}
	return d, nil
}
