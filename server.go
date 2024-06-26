package main

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/alecthomas/chroma/quick"
	"github.com/goccy/go-json"
)

var form = `
<form action="/" method="post">
  <label for="json">Insert JSON or URL</label><br />
  <textarea id="json" name="json" rows="50" cols="80"></textarea><br />
  <br />
  <input type="submit" value="Submit" />
</form>`

func NewServer() http.Handler {
	return http.HandlerFunc(handleRequest)
}

func getJSONData(jsonForm string) (map[string]interface{}, error) {
	var data map[string]interface{}
	if u, err := url.Parse(jsonForm); err == nil {
		return DownloadJSON(u)
	} else if err := json.Unmarshal([]byte(jsonForm), &data); err != nil {
		return nil, err
	}
	return data, nil
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		jsonForm := r.FormValue("json")
		data, err := getJSONData(jsonForm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		structDef, err := GenerateStruct(data, "Generated")
		w.Header().Set("Content-Type", "text/html")
		for ; err != nil; err = errors.Unwrap(err) {
			w.Write([]byte("<p>" + err.Error() + "</p>"))
		}
		quick.Highlight(w, structDef, "go", "html", "vs")
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(form))
}
