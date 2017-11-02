package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	//"io"
	"flag"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"
)

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
	Datum      time.Time
	Autor      string
	Kommentare []Kommentar
}

type Sitzung struct {
	Keks http.Cookie
	Datum time.Time
	Name string
}

var sitzungen []Sitzung

func kekse(w http.ResponseWriter) http.Cookie {
	c := http.Cookie{Name: "name", Value: "value"}
	http.SetCookie(w, &c)
	return c
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
		s.Dateiname = "seite/" + seite.Name()[:len(seite.Name())-5]
		start.Seiten = append(start.Seiten, s)
	}
	sort.Slice(start.Seiten, func(i, j int) bool { return start.Seiten[i].Datum.After(start.Seiten[j].Datum) })
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, start)
}

func pruefeLogin(name string, pass string) bool {
	return name!=""&&pass!=""
}

func machLogin(w http.ResponseWriter, r *http.Request) bool{
	name := r.FormValue("name")
	pass := r.FormValue("pass")
	if !pruefeLogin(name, pass) {
		return false
	}
	s := Sitzung{Name: name, Keks: kekse(w), Datum: time.Now()}
	sitzungen=append(sitzungen, s)
	fmt.Println("Login von", name)
	return true
}

func login(w http.ResponseWriter, r *http.Request) {
	if !machLogin(w,r) {
	t, _ := template.ParseFiles("login.html")
	t.Execute(w, nil)
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

func enthaelt(k []Kommentar, e Kommentar) bool {
	for _, a := range k {
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
		return
	}
	err = json.Unmarshal(dat, &s)
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(w, s)
}

func main() {
	port := flag.Int("port", 8000, "Port für den Webserver")
	flag.Parse()
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("css"))))
	http.Handle("/images/", http.StripPrefix("/images", http.FileServer(http.Dir("images"))))
	http.HandleFunc("/", startseite)
	http.HandleFunc("/login", login)
	http.HandleFunc("/seite/", seite)
	http.ListenAndServe(":"+strconv.Itoa(*port), nil)
}
