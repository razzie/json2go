package main

import (
	"flag"
	"net/http"
	"net/url"

	"github.com/alecthomas/chroma/quick"
)

var form = `
<form action="/" method="post">
  <label for="json">Insert JSON or URL</label><br />
  <textarea id="json" name="json" rows="50" cols="80"></textarea><br />
  <br />
  <input type="submit" value="Submit" />
</form>`

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":8080", "HTTP listen address")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			r.ParseForm()
			json := r.FormValue("json")

			if u, err := url.Parse(json); err == nil {
				json, err = DownloadJSON(u)
				if err != nil {
					http.Error(w, "Failed to download JSON: "+err.Error(), http.StatusInternalServerError)
					return
				}
			}

			structDef, err := GenerateStruct(json, "Generated")
			if err != nil {
				http.Error(w, "Failed to generate Go struct: "+err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "text/html")
			quick.Highlight(w, structDef, "go", "html", "vs")
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(form))
	})

	http.ListenAndServe(addr, nil)
}
