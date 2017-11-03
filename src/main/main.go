package main

/*
IMPORTE
*/
import (
	"bufio"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

/*
Konstanten
*/
const const_timeout int = 15
const const_port int = 8000

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

type Nutzerdaten struct {
	Profile []Profil
}

type Profil struct {
	Menu     []Menuitem
	Name     string
	Passwort string
	Meldung  string
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
var profile Nutzerdaten

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

func gebeProfil(name string) *Profil {
	for i, profil := range profile.Profile {
		if profil.Name == name {
			return &profile.Profile[i]
		}
	}
	return nil
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

func salzHash(name string, pass string) string {
	h := sha512.New()
	salz := name + pass
	return base64.URLEncoding.EncodeToString(h.Sum([]byte(salz)))
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
	pass = salzHash(name, pass)
	for _, profil := range profile.Profile {
		if profil.Name == name && profil.Passwort == pass {
			return true
		}
	}
	return false
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

func passwort(w http.ResponseWriter, r *http.Request) {
	login, name := gebeSitzung(r)
	if !login {
		http.Redirect(w, r, "/login", 302)
		return
	} else {
		passNeu := r.FormValue("neu_pass")
		passWdh := r.FormValue("pass_wdh")
		passAlt := r.FormValue("alt_pass")
		p := gebeProfil(name)
		p.Menu=nil
		p.Meldung=""
		if p.Passwort == salzHash(name, passAlt) {
			if passNeu == passWdh {
				p.Passwort = salzHash(name, passNeu)
				b, err := json.Marshal(profile)
				if err != nil {
					fmt.Println(err)
				}
				err = ioutil.WriteFile("user.json", b, 0644)
				if err != nil {
					fmt.Println(err)
				}
				p.Meldung = "Das Passwort wurde geändert"
			} else {
				p.Meldung = "Die Passwörter stimmen nicht überein"
			}
		} else {
			p.Meldung = "Das Passwort ist falsch"
		}
		t, _ := template.ParseFiles("profil.html")
		p.Menu = machMenu(p.Menu, r)
		t.Execute(w, p)
	}
}

func ladeProfile() {
	var p Nutzerdaten
	dat, err := ioutil.ReadFile("user.json")
	if err != nil {
		fmt.Println("Lesen der Nutzerdaten fehlgeschlagen oder noch keine Nutzer vorhanden", err)
		return
	}
	err = json.Unmarshal(dat, &p)
	profile = p
}

func erstelleNutzer() {
	input := bufio.NewReader(os.Stdin)
	fmt.Print("Namen des neuen Benutzers eingeben: ")
	name, _ := input.ReadString('\n')
	fmt.Print("Passwort des neuen Benutzers eingeben: ")
	pass, _ := input.ReadString('\n')
	profile.Profile = append(profile.Profile, Profil{Name: name[:len(name)-1], Passwort: salzHash(name[:len(name)-1], pass[:len(pass)-1])})
	b, err := json.Marshal(profile)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile("user.json", b, 0644)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print("Weitere Benutzer anlegen? (J/N): ")
	weiter, _ := input.ReadString('\n')
	if weiter == "J\n" || weiter == "j\n" {
		erstelleNutzer()
	} else {
		fmt.Println("Erstellen der Benutzer abgeschlossen")
	}
}

func main() {
	ladeProfile()
	port := flag.Int("port", const_port, "Port für den Webserver")
	flag.IntVar(&timeout, "timeout", const_timeout, "Timeout von Sitzungen in Minuten")
	var neuerNutzer bool
	flag.BoolVar(&neuerNutzer, "nutzer", false, "Neue Benutzer anlegen")
	flag.Parse()
	if neuerNutzer {
		erstelleNutzer()
	}
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("css"))))
	http.HandleFunc("/", startseite)
	http.HandleFunc("/login", login)
	http.HandleFunc("/passwort", passwort)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/profil", profil)
	http.HandleFunc("/seite/", seite)
	err := http.ListenAndServe(":"+strconv.Itoa(*port), nil)
	if err != nil {
		fmt.Println(err)
	}
}
