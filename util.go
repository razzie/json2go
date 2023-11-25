package main

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

func ToCamelCase(name string) string {
	parts := strings.Split(name, "_")
	for i, part := range parts {
		parts[i] = strings.Title(part)
	}
	return strings.Join(parts, "")
}

func DownloadJSON(u *url.URL) (string, error) {
	response, err := http.Get(u.String())
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
