package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func kekse(w http.ResponseWriter) {
	c := http.Cookie{Name: "name", Value: "value"}
	http.SetCookie(w, &c)
}

func startseite(w http.ResponseWriter, r *http.Request) {
	seiten, err := ioutil.ReadDir("seite")
	if err != nil {
		fmt.Println(err)
	}
	start := Startseite{}
	for _, seite := range seiten {
		var s Seite
		dat, err := ioutil.ReadFile("seite/" + seite.Name())
		if err != nil {
			fmt.Println(err)
			return
		}
		err = json.Unmarshal(dat, &s)
		s.Dateiname="seite/" + seite.Name()[:len(seite.Name())-5]
		start.Seiten = append(start.Seiten, s)
	}
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, start)
}

func login(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "TODO Formular für Login und hinterstehende Logik")
}

type Startseite struct {
	Seiten []Seite
}

type Kommentar struct {
	Autor  string
	Datum  time.Time
	Inhalt string
}

type Seite struct {
	Dateiname  string
	Titel      string
	Inhalt     string
	Datum      string
	Autor      string
	Kommentare []Kommentar
}

func enthaelt(s []Kommentar, e Kommentar) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func erstelleKommentar(r *http.Request) {
	k := Kommentar{Autor: r.URL.Query().Get("autor"), Inhalt: r.URL.Query().Get("inhalt"), Datum: time.Now()}
	if len(k.Autor) == 0 || len(k.Inhalt) == 0 {
		return
	}
	var s Seite
	dat, err := ioutil.ReadFile(r.URL.Path[1:] + ".json")
	if err != nil {
		fmt.Println("Erstellen des Kommentars fehlgeschlagen", err)
		return
	}
	err = json.Unmarshal(dat, &s)
	if enthaelt(s.Kommentare, k) {
		return
	}
	s.Kommentare = append([]Kommentar{k}, s.Kommentare...)
	if err != nil {
		fmt.Println(err)
	}
	b, _ := json.Marshal(s)
	ioutil.WriteFile(r.URL.Path[1:]+".json", b, 0644)
}

func seite(w http.ResponseWriter, r *http.Request) {
	erstelleKommentar(r)
	t, _ := template.ParseFiles("template.html")
	var s Seite
	dat, err := ioutil.ReadFile(r.URL.Path[1:] + ".json")
	if err != nil {
		http.Redirect(w, r, "/", 302)
		//io.WriteString(w, "404")
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
