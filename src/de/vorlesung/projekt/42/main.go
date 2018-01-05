/* Autoren: 3818468, 6985153, 9875672 */

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
const const_host string = ""

/*
Typen
*/

/*Eigenschaften für die Index-Seite*/
type Startseite struct {
	Menu   []Menuitem
	Seiten []Seite
}

/*Menüeinträge*/
type Menuitem struct {
	Ziel string
	Text string
}

/*einzelner Kommentar*/
type Kommentar struct {
	Autor  string
	Datum  time.Time
	Inhalt string
}

/*Eigenschaften für eine Blog-Seite*/
type Seite struct {
	Menu       []Menuitem
	Dateiname  string
	Titel      string
	Inhalt     string
	Datum      time.Time
	Autor      string
	Kommentare []Kommentar
	Bearbeitet time.Time
}

/*Nutzerdaten*/
type Nutzerdaten struct {
	Profile []Profil
}

/*einzelnes Profil*/
type Profil struct {
	Menu     []Menuitem
	Name     string
	Passwort string
	Meldung  string
}

/*einzelne Sitzung*/
type Sitzung struct {
	Keks  http.Cookie
	Datum time.Time
	Name  string
}

/*Eigenschaften für die Bearbeiten-Seite*/
type Bearbeiten struct {
	Menu      []Menuitem
	Meldung   string
	Titel     string
	Inhalt    string
	Dateiname string
}

/*
Globale Variablen
*/

var sitzungen []Sitzung
var timeout int
var profile Nutzerdaten
var templateProfil *template.Template
var templateSeite *template.Template
var templateLogin *template.Template
var templateNeu *template.Template
var templateLoeschen *template.Template
var templateIndex *template.Template

/*
Funktionen
*/

/*
Gibt an, ob eine Sitzung vorhanden ist
Falls Ja: Gibt Name des eingeloggten Users als zweiten Rückgabewert zurück
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

/*
Entfernt eine Sitzung
*/
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

/*
Gibt ein Profil anhand des übergebenen Profilnamens zurück
*/
func gebeProfil(name string) *Profil {
	for i, profil := range profile.Profile {
		if profil.Name == name {
			return &profile.Profile[i]
		}
	}
	return nil
}

/*
Gibt eine einzigartige ID mit der Länge laenge zurück
*/
func gebeUUID(laenge int) string {
	b := make([]byte, laenge)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println(err)
	}
	s := make([]string, len(b))
	for i := range b {
		s[i] = strconv.Itoa(int(b[i]))
	}
	return strings.Join(s, "")
}

/*
Setzt einen Cookie
*/
func kekse(w http.ResponseWriter) http.Cookie {
	ablauf := time.Now().Add(time.Minute * time.Duration(timeout))
	c := http.Cookie{Name: "id", Value: gebeUUID(128), Expires: ablauf}
	http.SetCookie(w, &c)
	return c
}

/*
Setzt einen Cookie mit dem Namen des Kommentators eines Kommentares
*/
func kommentarKeks(w http.ResponseWriter, name string) {
	ablauf := time.Now().Add(time.Minute * time.Duration(timeout))
	c := http.Cookie{Name: "Name", Value: name, Expires: ablauf}
	http.SetCookie(w, &c)
}

/*
Erstellt einen Hash aus Name und Passwort des Users
*/
func SalzHash(name string, pass string) string {
	h := sha512.New()
	salz := name + pass
	return base64.URLEncoding.EncodeToString(h.Sum([]byte(salz)))
}

/*
Liefert die Index-Seite mit einer Liste von (gekürzten) Blog-Einträgen aus
*/
func startseite(w http.ResponseWriter, r *http.Request) {
	seiten, err := ioutil.ReadDir("seite")
	start := Startseite{}
	if err != nil {
		start.Seiten = append(start.Seiten, Seite{Titel: "Noch keine Seiten vorhanden", Datum: time.Now(), Inhalt: "Schau später nochmal vorbei", Dateiname: "/", Autor: "SuperBlog"})
	}
	for _, seite := range seiten {
		var s Seite
		dat, err := ioutil.ReadFile("seite/" + seite.Name())
		if err != nil {
			fmt.Println(err)
			return
		}
		err = json.Unmarshal(dat, &s)
		if len(s.Inhalt) > 1000 {
			s.Inhalt = s.Inhalt[:1000] + "..."
		}
		s.Dateiname = "seite/" + seite.Name()[:len(seite.Name())-5]
		start.Seiten = append(start.Seiten, s)
	}
	if len(start.Seiten) == 0 {
		start.Seiten = append(start.Seiten, Seite{Titel: "Keine Seiten mehr vorhanden", Datum: time.Now(), Inhalt: "Schau später nochmal vorbei", Dateiname: "/", Autor: "SuperBlog"})
	}
	sort.Slice(start.Seiten, func(i, j int) bool { return start.Seiten[i].Datum.After(start.Seiten[j].Datum) })
	start.Menu = machMenu(start.Menu, r, 2)
	templateIndex.Execute(w, start)
}

/*
Prüft, ob der eingegebene Login gültig ist
*/
func pruefeLogin(name string, pass string) bool {
	pass = SalzHash(name, pass)
	for _, profil := range profile.Profile {
		if profil.Name == name && profil.Passwort == pass {
			return true
		}
	}
	return false
}

/*
Loggt einen User mit den eingetragenen Nutzerinformationen ein und erstellt eine neue Sitzung
*/
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

/*
Liefert die Login-Seite aus
*/
func login(w http.ResponseWriter, r *http.Request) {
	login, erfolg := machLogin(w, r)
	if !login {
		templateLogin.Execute(w, nil)
	} else {
		if erfolg {
			http.Redirect(w, r, "/", 302)
		} else {
			templateLogin.Execute(w, "Login fehlgeschlagen")
		}
	}
}

/*
Loggt einen User aus, löscht seine Sitzung und leitet ihn auf die Index-Seite um
*/
func logout(w http.ResponseWriter, r *http.Request) {
	loescheSitzung(r)
	http.Redirect(w, r, "/", 302)
}

/*
Prüft, ob ein Kommentar in einer Liste von Kommentaren bereits vorhanden ist
*/
func enthaelt(k []Kommentar, e Kommentar) bool {
	for _, a := range k {
		if a.Autor == e.Autor && a.Inhalt == e.Inhalt {
			return true
		}
	}
	return false
}

/*
Erstellt einen Kommentar aus den HTTP-Get- Parametern
*/
func erstelleKommentar(w http.ResponseWriter, r *http.Request) {
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

/*
Liefert eine Blog-Seite aus und erlaubt de Kommentar-Erstellung
*/
func seite(w http.ResponseWriter, r *http.Request) {
	erstelleKommentar(w, r)
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
	_, name := gebeSitzung(r)
	if name != "" {
		kommentarKeks(w, name)
	} else {
		if r.URL.Query().Get("autor") != "" {
			kommentarKeks(w, r.URL.Query().Get("autor"))
		}
	}
	menuitem := 0
	if s.Autor == name {
		menuitem = 1
	}
	s.Menu = machMenu(s.Menu, r, menuitem)
	templateSeite.Execute(w, s)
}

/*
Erstellt dynamisch ein Hauptmenü
seite:
1: Artikeldetailansicht
2: Startseite
3: Profil
0: Sonstiges
*/
func machMenu(m []Menuitem, r *http.Request, seite int) []Menuitem {
	m = append(m, Menuitem{Ziel: "/", Text: "Startseite"})
	login, name := gebeSitzung(r)
	if login {
		m = append(m, Menuitem{Ziel: "/profil", Text: name})
		switch seite {
		case 1:
			m = append(m, Menuitem{Ziel: "/bearbeiten/" + r.URL.Path[7:], Text: "Artikel bearbeiten"})
			m = append(m, Menuitem{Ziel: "/bestaetigen/" + r.URL.Path[7:], Text: "Artikel löschen"})
		case 2:
			m = append(m, Menuitem{Ziel: "/neu", Text: "Artikel erstellen"})
		}
		m = append(m, Menuitem{Ziel: "/logout", Text: "Logout"})
	} else {
		m = append(m, Menuitem{Ziel: "/login", Text: "Login"})
	}
	return m
}

/*
Liefert die Profil-Seite aus
*/
func profil(w http.ResponseWriter, r *http.Request) {
	login, name := gebeSitzung(r)
	if !login {
		http.Redirect(w, r, "/", 302)
		return
	}
	p := Profil{Name: name}
	p.Menu = machMenu(p.Menu, r, 3)
	templateProfil.Execute(w, p)
}

/*
Liefert die Passwort-Änderungsseite aus
*/
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
		p.Menu = nil
		p.Meldung = ""
		if p.Passwort == SalzHash(name, passAlt) {
			if passNeu == passWdh {
				p.Passwort = SalzHash(name, passNeu)
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
		t, _ := template.ParseFiles("res/profil.html")
		p.Menu = machMenu(p.Menu, r, 0)
		t.Execute(w, p)
	}
}

/*
Lädt alle Profile aus der user.json
*/
func ladeProfile(pfad string) {
	var p Nutzerdaten
	dat, err := ioutil.ReadFile(pfad)
	if err != nil {
		fmt.Println("Lesen der Nutzerdaten fehlgeschlagen oder noch keine Nutzer vorhanden", err)
		fmt.Println("Erstelle neue Nutzer...")
		erstelleNutzer()
		return
	}
	err = json.Unmarshal(dat, &p)
	profile = p
}

/*
Fügt einen Nutzer zur Nutzerliste hinzu
*/
func appendUser(name string, pass string) {
	profile.Profile = append(profile.Profile, Profil{Name: name[:len(name)-1], Passwort: SalzHash(name[:len(name)-1], pass[:len(pass)-1])})
}

/*
Ermöglicht es, beliebig viele Nutzer über die Kommandozeile anzulegen und prüft dabei auf bereits vorhanden sein
*/
func erstelleNutzer() {
	input := bufio.NewReader(os.Stdin)
	fmt.Print("Namen des neuen Benutzers eingeben: ")
	name, _ := input.ReadString('\n')
	if gebeProfil(name[:len(name)-1]) != nil {
		fmt.Print("Benutzername bereits vorhanden\n")
		erstelleNutzer()
		return
	}
	fmt.Print("Passwort des neuen Benutzers eingeben: ")
	pass, _ := input.ReadString('\n')
	appendUser(name, pass)
	b, err := json.Marshal(profile)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile("user.json", b, 0644)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print("Weitere Benutzer anlegen? (j/N): ")
	weiter, _ := input.ReadString('\n')
	if weiter == "J\n" || weiter == "j\n" {
		erstelleNutzer()
	} else {
		fmt.Println("Erstellen der Benutzer abgeschlossen")
	}
}

/*
Erstellt ein Verzeichnis mit dem angegebenen Pfad
*/
func erstelleVerzeichnis(pfad string) {
	if _, err := os.Stat(pfad); os.IsNotExist(err) {
		err := os.MkdirAll(pfad, 0711)
		if err != nil {
			fmt.Println(err)
		}
	}
}

/*
Liefert die Blogeintrag-Erstellenseite aus
*/
func erstelleSeite(w http.ResponseWriter, r *http.Request, altSeite Seite) {
	b := Bearbeiten{}
	b.Inhalt = altSeite.Inhalt
	b.Titel = altSeite.Titel
	titel := r.FormValue("titel")
	inhalt := r.FormValue("inhalt")
	if titel != "" || inhalt != "" {
		if titel == "" || inhalt == "" {
			b.Meldung = "Bitte Titel und Inhalt eintragen"
		} else {
			altSeite.Titel = titel
			altSeite.Inhalt = inhalt
			seitenjson, _ := json.Marshal(altSeite)
			erstelleVerzeichnis("seite")
			ioutil.WriteFile("seite/"+altSeite.Dateiname+".json", seitenjson, 0644)
			http.Redirect(w, r, "/seite/"+altSeite.Dateiname, 302)
		}
	}
	b.Menu = machMenu(b.Menu, r, 0)
	templateNeu.Execute(w, b)
}

/*
Erstellt einen neuen Blogeintrag
*/
func neu(w http.ResponseWriter, r *http.Request) {
	login, name := gebeSitzung(r)
	if !login {
		http.Redirect(w, r, "/", 302)
		return
	}
	dateiname := gebeUUID(8)
	erstelleSeite(w, r, Seite{Autor: name, Datum: time.Now(), Dateiname: dateiname})
}

/*
Bearbeitet einen bestehenden Blogeintrag
*/
func bearbeiten(w http.ResponseWriter, r *http.Request) {
	dateiname := r.URL.Path[len("/bearbeiten/"):]
	var s Seite
	dat, err := ioutil.ReadFile("seite/" + dateiname + ".json")
	if err != nil {
		fmt.Println("Seite zum Bearbeiten konnte nicht geöffnet werden", err)
		return
	}
	err = json.Unmarshal(dat, &s)
	login, name := gebeSitzung(r)
	if !login || s.Autor != name {
		http.Redirect(w, r, "/", 302)
		return
	}
	s.Dateiname = dateiname
	s.Bearbeitet = time.Now()
	erstelleSeite(w, r, s)
}

/*
Löscht einen Blogeintrag
*/
func loeschen(w http.ResponseWriter, r *http.Request) {
	dateiname := r.URL.Path[len("/loeschen/"):]
	var s Seite
	dat, err := ioutil.ReadFile("seite/" + dateiname + ".json")
	if err != nil {
		fmt.Println("Seite zum Löschen konnte nicht geöffnet werden", err)
		return
	}
	err = json.Unmarshal(dat, &s)
	login, name := gebeSitzung(r)
	if !login || s.Autor != name {
		http.Redirect(w, r, "/", 302)
		return
	}
	s.Dateiname = dateiname
	err = os.Remove("seite/" + dateiname + ".json")
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/", 302)
}

/*
Fragt nach einer Bestätigung, ob ein Blogeintrag gelöscht werden soll
*/
func bestaetigen(w http.ResponseWriter, r *http.Request) {
	dateiname := r.URL.Path[len("/bestaetigen/"):]
	b := Bearbeiten{Dateiname: dateiname}
	b.Menu = machMenu(b.Menu, r, 0)
	templateLoeschen.Execute(w, b)
}

/*
Startet das Programm, initialisiert die Variablen mit übergebenen Flags, parsed die Templates und startet den Webserver
*/
func main() {
	ladeProfile("user.json")
	port := flag.Int("port", const_port, "Port für den Webserver")
	host := flag.String("host", const_host, "Host für den Webserver")
	flag.IntVar(&timeout, "timeout", const_timeout, "Timeout von Sitzungen in Minuten")
	var neuerNutzer bool
	flag.BoolVar(&neuerNutzer, "nutzer", false, "Neue Benutzer anlegen")
	flag.Parse()
	if neuerNutzer {
		erstelleNutzer()
	}
	templateProfil = template.Must(template.ParseFiles("res/profil.html"))
	templateSeite = template.Must(template.ParseFiles("res/template.html"))
	templateNeu = template.Must(template.ParseFiles("res/neu.html"))
	templateLoeschen = template.Must(template.ParseFiles("res/loeschen.html"))
	templateLogin = template.Must(template.ParseFiles("res/login.html"))
	templateIndex = template.Must(template.ParseFiles("res/index.html"))
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("res/css"))))
	http.HandleFunc("/", startseite)
	http.HandleFunc("/login", login)
	http.HandleFunc("/passwort", passwort)
	http.HandleFunc("/neu", neu)
	http.HandleFunc("/loeschen/", loeschen)
	http.HandleFunc("/bestaetigen/", bestaetigen)
	http.HandleFunc("/bearbeiten/", bearbeiten)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/profil", profil)
	http.HandleFunc("/seite/", seite)
	err := http.ListenAndServeTLS(*host+":"+strconv.Itoa(*port), "server.crt", "server.key", nil)
	if err != nil {
		fmt.Println(err)
	}
}
