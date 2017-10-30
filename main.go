package main

import (
    "io"
    "net/http"
    "html/template"
)

func kekse(w http.ResponseWriter) {
    c := http.Cookie{Name: "name", Value: "value"}
    http.SetCookie(w, &c)
}

func startseite(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w,r,"seite/index.html",300)
}

func login(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "TODO Formular für Login und hinterstehende Logik" )
}

type Seite struct {
    Titel string
    Inhalt string
}

func seite(w http.ResponseWriter, r *http.Request) {
    t,_ := template.ParseFiles("template.html")
    s:=Seite{Titel:"Testseite", Inhalt:"Automatisch eingesetzter <br>Inhalt"}
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
