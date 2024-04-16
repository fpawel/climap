package messages

import (
	"bufio"
	"fmt"
	"github.com/fpawel/errorx"
	"log/slog"
	"math/rand/v2"
	"os"
	"slices"
	"strings"
)

type FileSentences []string

func (x FileSentences) NewMail() string {
	sentences := slices.Clone(x)
	rand.Shuffle(len(sentences), func(i, j int) {
		sentences[i], sentences[j] = sentences[j], sentences[i]
	})
	return fmt.Sprintf("From: <root@nsa.gov>\r\nSubject: %s\r\n\r\n%s", sentences[0], sentences[1])
}

func NewFileSentences(filePath string) (sentences FileSentences, _ error) {
	f, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return sentences, errorx.Prepend("read sentences from file").Wrap(err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("close sentences file", errorx.Attr(err))
		}
	}()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		s := strings.TrimSpace(sc.Text())
		if s == "" {
			continue
		}
		sentences = append(sentences, s)
	}
	if sc.Err() != nil {
		return sentences, errorx.Prepend("scan sentences file").Wrap(sc.Err())
	}
	return
}
