package main

import (
	"fmt"
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

	if response.StatusCode < 200 || response.StatusCode >= 400 {
		return "", fmt.Errorf("remote replied with %s", http.StatusText(response.StatusCode))
	}

	if response.Header.Get("Content-Type") != "application/json" {
		return "", fmt.Errorf("content-type is not application/json")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
