package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"slices"
	"strings"
)

var sentences []string

func newRandomMail() string {
	verbs := slices.Clone(sentences)
	rand.Shuffle(len(verbs), func(i, j int) {
		verbs[i], verbs[j] = verbs[j], verbs[i]
	})
	return fmt.Sprintf("From: <root@nsa.gov>\r\nSubject: %s\r\n\r\n%s", verbs[0], verbs[1])
}

func mustReadVSentencesFromFile(filePath string) (verbs []string) {
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
		verbs = append(verbs, s)
	}
	check(sc.Err(), "scan sentences file")
	if err := sc.Err(); err != nil {
		log.Fatalf("scan file error: %v", err)
		return
	}
	return
}
