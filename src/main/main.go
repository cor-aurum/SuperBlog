package main

import (
	"html/template"
	"io"
	"net/http"
)

func kekse(w http.ResponseWriter) {
	c := http.Cookie{Name: "name", Value: "value"}
	http.SetCookie(w, &c)
}

func startseite(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "seite/index.html", 300)
}

func login(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "TODO Formular für Login und hinterstehende Logik")
}

type Seite struct {
	Titel  string
	Inhalt string
	Datum string
	Autor string
}

func seite(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("template.html")
	s := Seite{Titel: "Testseite", Inhalt: "Automatisch eingesetzter <br>Inhalt", Datum: "31.10.2017", Autor: "Max Mustermann"}
	t.Execute(w, s)
}

func main() {
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("css"))))
	http.Handle("/images/", http.StripPrefix("/images", http.FileServer(http.Dir("images"))))
	http.HandleFunc("/", startseite)
	http.HandleFunc("/login", login)
	http.HandleFunc("/seite/", seite)
	http.ListenAndServe(":8000", nil)
}
