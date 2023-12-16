// All handler functions should be located in this file
package web

import (
	"html/template"
	"net/http"
)

func Root(w http.ResponseWriter, resp *http.Request) {
    tmpl := template.Must(template.ParseFiles("./static/index.html"))
    tmpl.Execute(w, nil)
}
