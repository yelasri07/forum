package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, page string, data any, status int) error {
	temp, err := template.ParseFiles("./frontend/templates/" + page)
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(status)
	err = temp.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return err
	}

	return nil
}

func RenderError(w http.ResponseWriter, status int) {
	e := struct {
		Type   string
		Status int
	}{
		Type:   http.StatusText(status),
		Status: status,
	}

	err := RenderTemplate(w, "error.html", e, status)
	if err != nil {
		log.Printf("Error rendering error page: %v", err)
	}
}
