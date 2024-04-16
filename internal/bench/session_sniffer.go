package bench

import (
	"bytes"
	"github.com/fpawel/errorx"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

type sessionSniffer struct {
	sessionID string
	buff      bytes.Buffer
	fd        *os.File
	wg        sync.WaitGroup
}

var reUUID = regexp.MustCompile(`\* OK New Cloud Technologies IMAP welcomes you -- ([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})`)

func (x *sessionSniffer) write(p []byte) (int, error) {
	if _, err := x.fd.Write(p); err != nil {
		return 0, errorx.Prepend("write").Wrap(err)
	}
	if err := x.fd.Sync(); err != nil {
		return 0, errorx.Prepend("flush").Wrap(err)
	}
	return len(p), nil
}

func (x *sessionSniffer) Write(p []byte) (int, error) {
	if x.fd != nil {
		return x.write(p)
	}
	xs := reUUID.FindStringSubmatch(string(p))
	if len(xs) != 2 {
		return x.buff.Write(p)
	}
	x.sessionID = xs[1]

	defer x.wg.Done()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return 0, errorx.Prepend("UserHomeDir").Wrap(err)
	}
	dir := filepath.Join(homeDir, ".climap")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return 0, errorx.Prepend("MkdirAll").Wrap(err)
	}
	filePath := filepath.Join(homeDir, ".climap", x.sessionID)
	x.fd, err = os.Create(filePath)
	if err != nil {
		return 0, errorx.Prepend("Create").Wrap(err)
	}
	if _, err := x.fd.Write(x.buff.Bytes()); err != nil {
		return 0, errorx.Prepend("write buffer").Wrap(err)
	}
	x.buff.Reset()
	return x.write(p)
}
