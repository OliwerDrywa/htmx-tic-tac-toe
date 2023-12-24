package handlers

import (
	"html/template"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseGlob("web/templates/*.html"))
	err := tmpl.ExecuteTemplate(w, "layout", nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
