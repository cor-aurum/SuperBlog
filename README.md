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


Für die Benutzung von SuperBlog wird ein TLS Zertifikat benötigt. Dieses kann entweder von einer Zertifizierungsstelle der Wahl erstellt oder selbst generiert werden. Bei dieser Variante ist allerdings zu beachten, dass der Browser des Endnutzers eine Warnung aussprechen wird.
Die Generierung wird unter Linux folgendermaßen initiiert:
openssl genrsa -out server.key 4096
openssl ecparam -genkey -name secp384r1 -out server.key
openssl req -new -x509 -sha512 -key server.key -out server.crt -days 3650

Es ist darauf zu achten, dass die Namen der Zertifikate bei server.key und server.crt belassen werden, da eine Zuordung sonst nicht stattfinden kann. Darüber hinaus müssen die Zertifikate im gleichen Verzeichnis liegen wie die im Repository vorhandenen .html Dateien.

## Theming
SuperBlog ist für ein einfaches Theming konzipiert. Im Verzeichnis css/ liegt eine Datei theme.css in der alle Farben der Weboberfläche definiert werden. Empfohlen seien beispielsweise folgende erste Zeilen (bestehende Einträge ersetzen):
  --element-width: 96vw;  
  --body-background-color: #ffff00;  
  --accent-color: #0000ff;  
  --background-color: #ff00b1;  
  --body-color: #ffff00;  
  --link-color: var(--body-color);  
