package main

import "testing"

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want bool
	}{
		// Válidas
		{"http simples", "http://example.com", true},
		{"https com caminho", "https://example.com/a/b?x=1", true},
		{"com porta", "https://example.com:8443", true},
		{"IPv4", "http://192.168.1.1:8080", true},
		{"IPv6 sem porta", "http://[::1]", true},
		{"IPv6 com porta", "http://[::1]:8080", true},
		{"scheme em maiúsculas", "HTTP://example.com", true},

		// Inválidas: scheme
		{"scheme ftp", "ftp://example.com", false},
		{"scheme mailto", "mailto:a@b.com", false},
		{"sem scheme", "example.com", false},
		{"scheme relativo", "//example.com", false},
		{"string vazia", "", false},

		// Inválidas: host
		{"sem host", "http://", false},
		{"sem host, com caminho", "http:///path", false},
		{"espaço no host", "http://exa mple.com", false},

		// Inválidas: porta
		{"porta não numérica", "http://example.com:abc", false},
		{"porta negativa", "http://example.com:-1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validateURL(tt.in)

			if got != tt.want {
				t.Errorf("validateURL(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestNormalizeURL(t *testing.T) {
	tests := []struct{ in, want string }{
		{"HTTP://Example.COM:80", "http://example.com/"},
		{"https://example.com:8443/a#x", "https://example.com:8443/a"},
		{"http://[::1]:80/", "http://[::1]/"},
		{"http://[::1]:8080/", "http://[::1]:8080/"},
		{"https://[2001:DB8::1]/p", "https://[2001:db8::1]/p"},
	}

	for _, tt := range tests {
		got, err := normalizeURL(tt.in)

		if err != nil {
			t.Fatalf("%s: unexpected error: %v", tt.in, err)
		}

		if got != tt.want {
			t.Errorf("normalizeURL(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
