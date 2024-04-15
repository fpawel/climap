package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"slices"
	"strings"
)

type testMailProvider []string

func (x testMailProvider) NewTestMail() string {
	sentences := slices.Clone(x)
	rand.Shuffle(len(sentences), func(i, j int) {
		sentences[i], sentences[j] = sentences[j], sentences[i]
	})
	return fmt.Sprintf("From: <root@nsa.gov>\r\nSubject: %s\r\n\r\n%s", sentences[0], sentences[1])
}

func mustNewtestMailProviderFromSentencesFile(filePath string) (sentences testMailProvider) {
	f, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	check(err, "read sentences from file")
	defer func() {
		check(f.Close(), "close sentences file")
	}()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		s := strings.TrimSpace(sc.Text())
		if s == "" {
			continue
		}
		sentences = append(sentences, s)
	}
	check(sc.Err(), "scan sentences file")
	return
}
