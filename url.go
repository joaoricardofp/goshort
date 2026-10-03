package main

import (
	"net"
	"net/url"
	"strconv"
	"strings"
)

func validateURL(u string) bool {
	parsed, err := url.Parse(u)

	if err != nil {
		return false
	}

	scheme := strings.ToLower(parsed.Scheme)

	if scheme != "http" && scheme != "https" {
		return false
	}

	if parsed.Host == "" {
		return false
	}

	if port := parsed.Port(); port != "" {
		if _, err := strconv.Atoi(port); err != nil {
			return false
		}
	}

	return true
}

func normalizeURL(u string) (string, error) {
	parsed, err := url.Parse(u)

	if err != nil {
		return "", err
	}

	parsed.Scheme = strings.ToLower(parsed.Scheme)

	hostname := strings.ToLower(parsed.Hostname())
	port := parsed.Port()

	switch {
	case parsed.Scheme == "http" && port == "80":
		port = ""
	case parsed.Scheme == "https" && port == "443":
		port = ""
	}

	if port != "" {
		parsed.Host = net.JoinHostPort(hostname, port)
	} else if strings.Contains(hostname, ":") {
		parsed.Host = "[" + hostname + "]"
	} else {
		parsed.Host = hostname
	}

	if parsed.Path == "" {
		parsed.Path = "/"
	}

	parsed.Fragment = ""

	return parsed.String(), nil
}
