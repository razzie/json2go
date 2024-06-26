package main

import (
	"fmt"
	"mime"
	"net"
	"net/http"
	"net/url"
	"strings"
	"unicode"

	"github.com/goccy/go-json"
)

func getPrivateIPBlocks() (blocks []*net.IPNet) {
	// https://stackoverflow.com/a/50825191
	for _, cidr := range []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"169.254.0.0/16", // RFC3927 link-local
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
		"fc00::/7",       // IPv6 unique local addr
	} {
		_, block, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(fmt.Errorf("parse error on %q: %v", cidr, err))
		}
		blocks = append(blocks, block)
	}
	return
}

var privateIPBlocks = getPrivateIPBlocks()

func IsPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

func GetHost(u *url.URL) string {
	host := u.Host
	if strings.HasPrefix(host, "[") && strings.Contains(host, "]") {
		host = host[1:strings.Index(host, "]")]
	} else if strings.Contains(host, ":") {
		parts := strings.Split(host, ":")
		host = parts[0]
	}
	return host
}

func ToCamelCase(name string) string {
	if len(name) == 0 {
		return "X"
	}
	if unicode.IsDigit([]rune(name)[0]) {
		name = "X" + name
	}
	name = strings.ReplaceAll(name, "-", "_")
	parts := strings.Split(name, "_")
	for i, part := range parts {
		runes := make([]rune, 0, len(part))
		for _, r := range part {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
				continue
			}
			runes = append(runes, r)
		}
		if len(runes) > 0 {
			runes[0] = unicode.ToUpper(runes[0])
		}
		parts[i] = string(runes)
	}
	return strings.Join(parts, "")
}

func DownloadJSON(u *url.URL) (map[string]interface{}, error) {
	ips, err := net.LookupIP(GetHost(u))
	if err != nil {
		return nil, err
	}
	for _, ip := range ips {
		if IsPrivateIP(ip) {
			return nil, fmt.Errorf("access to private networks is restricted")
		}
	}

	response, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 400 {
		return nil, fmt.Errorf("remote replied with %s", http.StatusText(response.StatusCode))
	}

	contentType, _, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("cannot parse content-type: %v", err)
	}
	if contentType != "application/json" {
		return nil, fmt.Errorf("content-type is not application/json")
	}

	var data map[string]interface{}
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}
