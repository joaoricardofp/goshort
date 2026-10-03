package main

import (
	"sync"
	"testing"
)

func TestGetOrCreateConcorrente(t *testing.T) {
	s := NewStore()
	const n = 100

	codes := make([]string, n)
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			code, err := s.GetOrCreate("https://Example.com:443/a")

			if err != nil {
				t.Errorf("erro: %v", err)
				return
			}

			codes[i] = code
		}(i)
	}

	wg.Wait()

	for i, c := range codes {
		if c != codes[0] {
			t.Fatalf("goroutine %d code %s, expected %s", i, c, codes[0])
		}
	}

	if len(s.urlToCode) != 1 || len(s.codeToURL) != 1 {
		t.Errorf("urlToCode: %d, codeToURL: %d", len(s.urlToCode), len(s.codeToURL))
	}
}

func TestEncodeCode(t *testing.T) {
	tests := []struct {
		in   int
		want string
	}{
		{1, "000001"},
		{10, "00000a"},
		{36, "000010"},
		{2176782335, "zzzzzz"},  // último com 6 caracteres
		{2176782336, "1000000"}, // primeiro com 7
	}

	for _, tt := range tests {
		if got := encodeCode(tt.in); got != tt.want {
			t.Errorf("encodeCode(%d) = %q, esperado %q", tt.in, got, tt.want)
		}
	}
}
