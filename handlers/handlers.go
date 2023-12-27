package handlers

import (
	"html/template"
	"net/http"
)

func StaticFileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))).ServeHTTP(w, r)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/fragments/layout.html", "web/templates/routes/index.html"))
	err := tmpl.ExecuteTemplate(w, "layout", nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func TicTacToeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/fragments/layout.html", "web/templates/routes/tic-tac-toe.html"))
	err := tmpl.ExecuteTemplate(w, "layout", nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
