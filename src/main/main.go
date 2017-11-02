package main

/*
IMPORTE
*/
import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

/*
Typen
*/

type Startseite struct {
	Menu   []Menuitem
	Seiten []Seite
}

type Menuitem struct {
	Ziel string
	Text string
}

type Kommentar struct {
	Autor  string
	Datum  time.Time
	Inhalt string
}

type Seite struct {
	Menu       []Menuitem
	Dateiname  string
	Titel      string
	Inhalt     string
	Datum      time.Time
	Autor      string
	Kommentare []Kommentar
}

type Profil struct {
	Menu []Menuitem
	Name string
}

type Sitzung struct {
	Keks  http.Cookie
	Datum time.Time
	Name  string
}

/*
Globale Variablen
*/

var sitzungen []Sitzung
var timeout int

/*
Funktionen
*/

func gebeSitzung(r *http.Request) (bool, string) {
	keks, err := r.Cookie("id")
	if err != nil {
		return false, ""
	}
	for _, s := range sitzungen {
		if s.Keks.Value == keks.Value {
			return true, s.Name
		}
	}
	return false, ""
}

func loescheSitzung(r *http.Request) {
	keks, err := r.Cookie("id")
	if err != nil {
		return
	}
	for i, s := range sitzungen {
		if s.Keks.Value == keks.Value {
			sitzungen = append(sitzungen[:i], sitzungen[i+1:]...)
		}
	}
}

func kekse(w http.ResponseWriter) http.Cookie {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println(err)
	}
	s := make([]string, len(b))
	for i := range b {
		s[i] = strconv.Itoa(int(b[i]))
	}
	ablauf := time.Now().Add(time.Minute * time.Duration(timeout))
	c := http.Cookie{Name: "id", Value: strings.Join(s, ""), Expires: ablauf}
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
	start.Menu = machMenu(start.Menu, r)
	t.Execute(w, start)
}

func pruefeLogin(name string, pass string) bool {
	return (name == "test" || name == "admin") && pass != ""
}

func machLogin(w http.ResponseWriter, r *http.Request) (bool, bool) {
	name := r.FormValue("name")
	pass := r.FormValue("pass")
	if name == "" && pass == "" {
		return false, false
	} else {
		if !pruefeLogin(name, pass) {
			return true, false
		}
	}
	s := Sitzung{Name: name, Keks: kekse(w), Datum: time.Now()}
	sitzungen = append(sitzungen, s)
	return true, true
}

func login(w http.ResponseWriter, r *http.Request) {
	login, erfolg := machLogin(w, r)
	if !login {
		t, _ := template.ParseFiles("login.html")
		t.Execute(w, nil)
	} else {
		if erfolg {
			http.Redirect(w, r, "/", 302)
		} else {
			t, _ := template.ParseFiles("login.html")
			t.Execute(w, "Login fehlgeschlagen")
		}
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	loescheSitzung(r)
	http.Redirect(w, r, "/", 302)
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
	s.Menu = machMenu(s.Menu, r)
	t.Execute(w, s)
}

func machMenu(m []Menuitem, r *http.Request) []Menuitem {
	m = append(m, Menuitem{Ziel: "/", Text: "Startseite"})
	login, name := gebeSitzung(r)
	if login {
		m = append(m, Menuitem{Ziel: "/profil", Text: name})
		m = append(m, Menuitem{Ziel: "/logout", Text: "Logout"})
	} else {
		m = append(m, Menuitem{Ziel: "/login", Text: "Login"})
	}
	return m
}

func profil(w http.ResponseWriter, r *http.Request) {
	login, name := gebeSitzung(r)
	if !login {
		http.Redirect(w, r, "/", 302)
		return
	}
	p := Profil{Name: name}
	t, _ := template.ParseFiles("profil.html")
	p.Menu = machMenu(p.Menu, r)
	t.Execute(w, p)
}

func main() {
	port := flag.Int("port", 8000, "Port für den Webserver")
	flag.IntVar(&timeout, "timeout", 15, "Timeout von Sitzungen in Minuten")
	flag.Parse()
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("css"))))
	http.Handle("/images/", http.StripPrefix("/images", http.FileServer(http.Dir("images"))))
	http.HandleFunc("/", startseite)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/profil", profil)
	http.HandleFunc("/seite/", seite)
	http.ListenAndServe(":"+strconv.Itoa(*port), nil)
}
