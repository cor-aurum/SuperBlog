package main

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"io/ioutil"
	"fmt"
)

func kekse(w http.ResponseWriter) {
	c := http.Cookie{Name: "name", Value: "value"}
	http.SetCookie(w, &c)
}

func startseite(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "seite/index", 301)
}

func login(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "TODO Formular für Login und hinterstehende Logik")
}

type Seite struct {
	Titel  string
	Inhalt string
	Datum  string
	Autor  string
}

func seite(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("template.html")
	var s Seite
	dat, err := ioutil.ReadFile(r.URL.Path[1:]+".json")
	if err != nil {
		http.Redirect(w, r, "index", 302)
		return
	}
	err = json.Unmarshal(dat, &s)
	if err != nil {
		fmt.Println(err)
	}
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
