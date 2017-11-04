# SuperBlog
Ein Webblog, geschrieben in der Programmiersprache Go.

## Benutzung
go run /src/main/main.go [Optionen]
oder:
go build -o SuperBlog src/main/main.go
./SuperBlog [Optionen]

Mögliche Optionen sind:  
-nutzer  
    *Neue Benutzer anlegen*  
-port int  
    *Port für den Webserver (default 8000)*  
-timeout int  
    *Timeout von Sitzungen in Minuten (default 15)*  
        
        
Beim ersten Start des Programmes kann ohne Angabe von Optionen ein oder mehrere Benutzer angelegt werden.  

Für den Leser existiert eine übersichtliche Oberfläche, in der er ohne Anmeldung Artikel kommentieren, 
auf der Startseite über neue Artikel scrollen und über gelesene Artikel mit anderen Lesern diskutieren.  

Für die Kommentarfunktion ist ein einfacher Spamschutz integriert, der verhindert, dass ein Leser einen Artikel 
mehrfach mit gleichem Text kommentiert.  

Sowohl Kommentare als auch die Artikel auf der Startseite sind chronologisch angeordnet, sodass kein Eintrag übersehen wird.  
