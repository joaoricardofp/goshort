package main

import (
	"strconv"
	"strings"
	"sync"
)

const minCodeLen = 6

type Store struct {
	mu        sync.Mutex
	urlToCode map[string]string
	codeToURL map[string]string
	counter   int
}

func NewStore() *Store {
	return &Store{
		urlToCode: make(map[string]string),
		codeToURL: make(map[string]string),
	}
}

func encodeCode(n int) string {
	s := strconv.FormatInt(int64(n), 36)

	if len(s) < minCodeLen {
		s = strings.Repeat("0", minCodeLen-len(s)) + s
	}

	return s
}

func (s *Store) GetOrCreate(rawURL string) (string, bool, error) {
	normalized, err := normalizeURL(rawURL)

	if err != nil {
		return "", false, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if code, ok := s.urlToCode[normalized]; ok {
		return code, false, nil
	}

	var code string

	for {
		s.counter++
		code = encodeCode(s.counter)

		if _, used := s.codeToURL[code]; !used {
			break
		}
	}

	s.urlToCode[normalized] = code
	s.codeToURL[code] = normalized

	return code, true, nil
}

func (s *Store) Lookup(code string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	url, ok := s.codeToURL[code]
	return url, ok
}
