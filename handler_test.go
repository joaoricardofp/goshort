package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func postShorten(h http.Handler, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("POST", "/shorten", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestShortenCriadoEReutilizado(t *testing.T) {
	h := NewHandler(NewStore()).Routes()

	a := postShorten(h, `{"url":"https://example.com/a"}`)
	b := postShorten(h, `{"url":"HTTPS://Example.COM:443/a#x"}`)

	if a.Code != http.StatusCreated {
		t.Errorf("primeira chamada: status %d, esperado 201", a.Code)
	}
	if b.Code != http.StatusOK {
		t.Errorf("segunda chamada: status %d, esperado 200", b.Code)
	}

	var ra, rb shortenResponse
	json.Unmarshal(a.Body.Bytes(), &ra)
	json.Unmarshal(b.Body.Bytes(), &rb)

	if ra.Code != rb.Code {
		t.Errorf("códigos diferentes: %q e %q", ra.Code, rb.Code)
	}
	if want := "http://example.com/" + ra.Code; ra.ShortURL != want {
		t.Errorf("short_url = %q, esperado %q", ra.ShortURL, want)
	}
}

func TestShortenEntradasInvalidas(t *testing.T) {
	h := NewHandler(NewStore()).Routes()

	casos := map[string]string{
		"json quebrado": `{"url":`,
		"sem campo url": `{}`,
		"scheme ftp":    `{"url":"ftp://example.com"}`,
		"sem host":      `{"url":"http://"}`,
	}

	for nome, body := range casos {
		t.Run(nome, func(t *testing.T) {
			if rec := postShorten(h, body); rec.Code != http.StatusBadRequest {
				t.Errorf("status %d, esperado 400", rec.Code)
			}
		})
	}
}

func TestRedirect(t *testing.T) {
	h := NewHandler(NewStore()).Routes()
	postShorten(h, `{"url":"https://example.com/a"}`) // gera "000001"

	req := httptest.NewRequest("GET", "/000001", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("status %d, esperado 302", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "https://example.com/a" {
		t.Errorf("Location = %q", loc)
	}
}

func TestRedirectNaoEncontrado(t *testing.T) {
	h := NewHandler(NewStore()).Routes()

	req := httptest.NewRequest("GET", "/zzzzzz", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status %d, esperado 404", rec.Code)
	}
}
