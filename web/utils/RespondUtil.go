package utils

import (
	"html/template"
	"io"
	"log"
	"net/http"
)

func Respond(w http.ResponseWriter, status int, html string) {
	w.WriteHeader(status)
	io.WriteString(w, html)
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Length", string(len(html)))
}

func RespondNotFound(w http.ResponseWriter) {
	Respond(w, http.StatusNotFound, "<h1>Page Not Found</h1>")
}

func RespondServerError(w http.ResponseWriter) {
	Respond(w, http.StatusInternalServerError, "<h1>服务器内部错误</h1>")
}

func RespondTemplate(w http.ResponseWriter, status int, templateFile string, data interface{}) {
	w.WriteHeader(status)
	t, err := template.ParseFiles(templateFile)
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, data)
}
