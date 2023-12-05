package main

import (
	"flag"
	"net/http"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":8080", "HTTP listen address")
	flag.Parse()
	http.ListenAndServe(addr, NewServer())
}
