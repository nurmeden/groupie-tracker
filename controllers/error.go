package controllers

import (
	"html/template"
	"net/http"
)

type ErrorStatus struct {
	Code    int
	Message string
}

func HandlerErrors(w http.ResponseWriter, status int) {
	tmpl, err := template.ParseFiles("error.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	Message := ErrorStatus{status, http.StatusText(status)}
	err = tmpl.Execute(w, Message)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
