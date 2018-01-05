package main

/* Autoren: 3818468, 6985153, 9875672 */

/*
Wegen vielen in unserem Design hardcoded eingebauten Pfaden oder Überprüfungen,
sind viele Tests nicht oder nur mit einem erwarteten Fehler möglich. Die Vorles-
ung zur Nutzung von übergebenen Providern kam dafür zu spät.
*/

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func TestCreateTestUser(t *testing.T) {
	t.Log("adding a loaded user to global variable.")
	name := "Max Mustermann\n"
	pass := "1234567890\n"
	t.Log("mockup name: ", name)
	t.Log("mockup pass: ", pass)
	appendUser(name, pass)
	t.Log("added %v successfully", name)
}

/* erstelleVerzeichnis File Sytem function. */
func TestErstelleVerzeichnisSeite(t *testing.T) {
	t.Log("Create a Folder. ... (Expecting the Folder is existing afterwards)")
	erstelleVerzeichnis("./seite")
	if _, err := os.Stat("./seite"); err != nil {
		if os.IsNotExist(err) {
			t.Error("File has not been created")
		} else {
			t.Log("File exists as expected")
		}
	}
}

func TestErstelleVerzeichnisExampleFolder(t *testing.T) {
	defer func() {
		var err = os.Remove("examplefolder")
		if err != nil {
			return
		}
		t.Log("==> done deleting file")
	}()
	t.Log("Create a Folder. ... (Expecting the Folder is existing afterwards)")
	erstelleVerzeichnis("examplefolder")
	if _, err := os.Stat("examplefolder"); err != nil {
		if os.IsNotExist(err) {
			t.Error("File has not been created")
		} else {
			t.Log("File exists as expected")
		}
	}
}

/* convert username + password to a sha256 hashed base64 encoded string*/

func TestSalzHash(t *testing.T) {
	t.Log("check if the generated salted Hash is correct. ... (Expect it matches)")
	name := "Max Mustermann"
	pass := "1234567890"
	expected := "TWF4IE11c3Rlcm1hbm4xMjM0NTY3ODkwz4PhNX7vuL3xVChQ1m2AB9Yg5AULVxXcg_SpIdNs6c5H0NE8XYXysP-DGNKHfuwvY7kxvUdBeoGlODJ6-SfaPg=="
	hashed := SalzHash(name, pass)
	t.Log("mockup name: ", name)
	t.Log("mockup pass: ", pass)
	t.Log("expected output: ", expected)
	t.Log("actual output:   ", hashed)
	if expected == hashed {
		t.Log("matching")
	} else {
		t.Error("Not matching! Wrong hashed String")
	}

}

/* loads the user file to global variable. Just test if not existing in current architecture. Expect Failure */
func TestLadeProfileNotExisting(t *testing.T) {
	t.Log("test loading of not existing user.json file. ... (Expect code panic)")
	defer func() {
		if r := recover(); r == nil {
			t.Error("The code did not panic")
		} else {
			t.Log("ladeProfile failed as expected")
		}
	}()
	ladeProfile("notuser.json")
}

func TestLadeProfile(t *testing.T) {
	if _, err := os.Stat("user.json"); err != nil {
		if os.IsNotExist(err) {
			LadeProfileNotCreated(t)
		} else {
			LadeProfileAlreadyCreated(t)
		}
	}
}

func LadeProfileNotCreated(t *testing.T) {
	t.Log("test loading of not existing user.json file. ... (Expect code panic)")
	defer func() {
		if r := recover(); r == nil {
			t.Error("The code did not panic")
		} else {
			t.Log("ladeProfile failed as expected")
		}
	}()
	ladeProfile("notuser.json")
}

func LadeProfileAlreadyCreated(t *testing.T) {
	t.Log("test loading of existing user.json file. ... (do not Expect code panic)")
	defer func() {
		if r := recover(); r != nil {
			t.Error("The code did panic")
		} else {
			t.Log("The file got loaded")
		}
	}()
	ladeProfile("user.json")
}

/* searches for an profile by name. Will fail because of no profiles loaded. */

func TestGebeProfilNotExisting(t *testing.T) {
	t.Log("try to get a profile. ... (Expect nil return)")
	name := "Nicht Max Mustermann"
	t.Log("mockup name: ", name)
	result := gebeProfil(name)
	if result == nil {
		t.Log("no profile as expected")
	} else {
		t.Error("returned a non existing profile")
	}
}

func TestGebeProfilExisting(t *testing.T) {
	t.Log("try to get a profile. ... (Expect the profile returned)")
	name := "Max Mustermann"
	t.Log("tryout name: ", name)
	result := gebeProfil(name)
	if result == nil {
		t.Error("does not return an existing profile")
	} else {
		t.Log("returned the profile of Max Mustermann")
	}
}

/* checks if username and password are matching correctly. Will fail because of empty user.json.*/

func TestPruefeLoginNotExistingUser(t *testing.T) {
	t.Log("try to login with a not existing user. ... (Expect login failure)")
	name := "Nicht Max Mustermann"
	pass := "1234567890"
	t.Log("mockup name: ", name)
	t.Log("mockup pass: ", pass)
	result := pruefeLogin(name, pass)
	if result == true {
		t.Error("logged in with not existing user credentials!")
	} else {
		t.Log("login failed as expected")
	}
}

/* checks if username and password are matching correctly. Using wrong password. Expecting failure. */
func TestPruefeLoginWrongPassword(t *testing.T) {

	t.Log("try to login with a wrong password. ... (Expect login failure)")
	name := "Max Mustermann"
	wrongpass := "asdfqwerty"
	t.Log("tryout name: ", name)
	t.Log("tryout pass: ", wrongpass)
	result := pruefeLogin(name, wrongpass)
	if result == true {
		t.Error("logged in with wrong user credentials!")
	} else {
		t.Log("login failed as expected")
	}
}

/* checks if username and password are matching correctly. Expect success*/
func TestPruefeLoginSuccessful(t *testing.T) {

	t.Log("try to login with real credentials. ... (Expect login success)")
	name := "Max Mustermann"
	pass := "1234567890"
	t.Log("tryout name: ", name)
	t.Log("tryout pass: ", pass)
	result := pruefeLogin(name, pass)
	if result == false {
		t.Error("Not logged in!")
	} else {
		t.Log("login succesfull as expected")
	}
}

/* generate a random UUID. Just testing if it is as long as requiered. Collision Test is stupid. */
/* The Byte array contains numbers with one to three digits. This makes the UUID String have a variable size between 1*128 and 3*128. */
func TestGebeUUID(t *testing.T) {
	t.Log("test if the generated UUID is as long a requested. ... (expecting it is)")
	laenge := 128
	t.Log("length: ", laenge)
	UUID := gebeUUID(laenge)
	t.Log(UUID)
	if len(UUID) >= 128 && len(UUID) <= 3*128 {
		t.Log("UUID like it schould be ", len(UUID))
	} else {
		t.Error("UUID length not like specified ", len(UUID))
	}
}

func TestGebeUUIDFail(t *testing.T) {
	t.Log("test if negative length is allowed. ... (expecting it is not)")
	defer func() {
		if r := recover(); r == nil {
			t.Error("negative length is allowed")
		} else {
			t.Log("negative length is not allowed")
		}
	}()
	laenge := -128
	t.Log("length: ", laenge)
	gebeUUID(laenge)

}

/* checks if given values are not empty. uses pruefeLogin. appends new session to sesionslist. Will fail because of empty user.json*/

func TestMachLoginEmptyData(t *testing.T) {
	t.Log("Test with both values empty. ... (Expecting (false, false))")
	t.Log("create a request")
	r, err := http.NewRequest("post", "Login", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(r)

	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder and serve the Handler")
	w := httptest.NewRecorder()
	var login bool
	var erfolg bool
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login, erfolg = machLogin(w, r)
		w.Header().Set("Content-Type", "application/json")
	})
	handler.ServeHTTP(w, r)

	t.Log("Check the HTTP StatusCode. ... (Expecting OK)")
	if status := w.Code; status != http.StatusOK {
		t.Errorf("unexpected status code: %v instead of %v",
			status, http.StatusOK)
	}
	// Check the returned values
	if login == false && erfolg == false {
		t.Logf("No Login like Expected: (%v and %v)", login, erfolg)
	} else {
		t.Errorf("unexpected response: (%v and %v) anstatt (false and false)", login, erfolg)
	}
}

func TestMachLoginRealData(t *testing.T) {

	t.Log("Test with a Testlogin. ... (Expecting (true, true))")
	t.Log("create a request")
	r, err := http.NewRequest("post", "/login", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}
	q := r.URL.Query()
	q.Add("name", "Max Mustermann")
	q.Add("pass", "1234567890")
	r.URL.RawQuery = q.Encode()
	t.Log(r.FormValue("name"))

	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder and serve the Handler")
	w := httptest.NewRecorder()
	var login bool
	var erfolg bool
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login, erfolg = machLogin(w, r)
		w.Header().Set("Content-Type", "application/json")
	})
	handler.ServeHTTP(w, r)

	t.Log("Check the HTTP StatusCode. ... (Expecting OK)")
	if status := w.Code; status != http.StatusOK {
		t.Errorf("unexpected status code: %v instead of %v",
			status, http.StatusOK)
	}
	// Check the returned values
	if login == true && erfolg == true {
		t.Logf("Login like Expected: (%v and %v)", login, erfolg)
	} else {
		t.Errorf("unexpected response: (%v and %v) anstatt (true and true)", login, erfolg)
	}
}

func TestMachLoginWrongData(t *testing.T) {
	t.Log("Test with a wrong Testlogin. ... (Expecting (true, false))")
	t.Log("create a request")
	r, err := http.NewRequest("post", "/login", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}
	q := r.URL.Query()
	q.Add("name", "Max Mustermann")
	q.Add("pass", "false")
	r.URL.RawQuery = q.Encode()
	t.Log(r.FormValue("name"))

	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder and serve the Handler")
	w := httptest.NewRecorder()
	var login bool
	var erfolg bool
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login, erfolg = machLogin(w, r)
		w.Header().Set("Content-Type", "application/json")
	})
	handler.ServeHTTP(w, r)

	t.Log("Check the HTTP StatusCode. ... (Expecting OK)")
	if status := w.Code; status != http.StatusOK {
		t.Errorf("unexpected status code: %v instead of %v",
			status, http.StatusOK)
	}
	// Check the returned values
	if login == true && erfolg == false {
		t.Logf("Login like Expected: (%v and %v)", login, erfolg)
	} else {
		t.Errorf("unexpected response: (%v and %v) anstatt (true and true)", login, erfolg)
	}
}

/* returns a session with an specific id. Will fail because no sessions are running*/

func TestGebeSitzungSuccess(t *testing.T) {
	t.Log("reading Session variable. ... (Expecting return of a session")
	t.Logf("%v sessions are running", len(sitzungen))
	t.Log(sitzungen)
	t.Log("Create a request in the running session")
	session := sitzungen[0]
	cookie := session.Keks
	r, err := http.NewRequest("get", "", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}
	r.AddCookie(&cookie)

	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder and serve the Handler")
	w := httptest.NewRecorder()
	var sessionExists bool
	var name string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionExists, name = gebeSitzung(r)
		w.Header().Set("Content-Type", "application/json")
	})
	handler.ServeHTTP(w, r)

	if sessionExists == true && name == "Max Mustermann" {
		t.Logf("there is a Session running! : %v, %v ", sessionExists, name)
	} else {
		t.Errorf("No session returned !")
	}
}

/* loescheSitzung deletes a session with an specifiic id.*/
func TestLoescheSitzung(t *testing.T) {
	t.Log("Test deletion of a session. ... (Expect deletion of session)")
	t.Logf("%v sessions are running", len(sitzungen))
	t.Log(sitzungen)
	t.Log("Create a request in the running session")
	session := sitzungen[0]
	cookie := session.Keks
	r, err := http.NewRequest("get", "", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}
	r.AddCookie(&cookie)

	loescheSitzung(r)
	t.Log("Test if session is removed")
	if len(sitzungen) > 0 {
		t.Error("Sessions still active")
	} else {
		t.Log("Session terminated")
	}

}

/* kommentarKeks void method. */
func TestKommentarKeks(t *testing.T) {
	t.Log("Test creation of a Kommentar Cookie. ... (Expecting creation of cookie)")

	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder")
	w := httptest.NewRecorder()

	name := "Test-Cookie"

	kommentarKeks(w, name)

	request := &http.Request{Header: http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}}

	// Extract the dropped cookie from the request.
	cookie, err := request.Cookie("Name")
	if err != nil {
		t.Errorf("Failed to read 'name' Cookie: %v", err)
	}
	if cookie.String() == "Name=Test-Cookie" {
		t.Log(cookie)
	} else {
		t.Errorf("unexpected cookie: %v", cookie)
	}
}

/* startseite  http template method */

func TestStartseiteNoFolder(t *testing.T) {
	t.Log("Test Startseite with no Folder. ... (Expecting http OK)")
	templateIndex = template.Must(template.ParseFiles("res/index.html"))
	t.Log("make sure folder 'seite' does not exist.")
	if _, err := os.Stat("seite"); err == nil {
		t.Log("seite does exist")
		RemoveContents("seite")
		t.Log("delete folder")
		var err = os.Remove("seite")
		if err != nil {
			t.Error("Could not delete 'seite'")
		}
		t.Log("==> done deleting 'seite'")

		t.Log("create a request")
		r, err := http.NewRequest("post", "/", nil)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
		if err != nil {
			t.Fatal(err)
		}

		/*  ResponseRecorder instead of ResponseWriter to record the response. */
		t.Log("create a ResponseRecorder and serve the Handler")
		w := httptest.NewRecorder()
		handler := http.HandlerFunc(startseite)
		handler.ServeHTTP(w, r)

		t.Log("Check the HTTP StatusCode. ... (Expecting Http OK)")
		if status := w.Code; status == 302 {
			t.Errorf("unexpected status code: %v instead of 200",
				status)
		}
	}
}

func TestStartseiteEmptyFolder(t *testing.T) {
	defer func() {
		t.Log("make sure folder 'seite' gets removed.")
		if _, err := os.Stat("seite"); err == nil {
			t.Log("seite does exist")
			RemoveContents("seite")
			t.Log("delete folder")
			var err = os.Remove("seite")
			if err != nil {
				t.Error("Could not delete 'seite'")
			}
			t.Log("==> done deleting 'seite'")
		}
	}()

	t.Log("Test Startseite with empty folder. ... (Expecting http OK)")
	templateIndex = template.Must(template.ParseFiles("res/index.html"))

	t.Log("make sure folder 'seite' does exist.")
	if _, err := os.Stat("seite"); os.IsNotExist(err) {
		t.Log("seite does not exist")
		t.Log("create folder")
		var err = os.MkdirAll("seite", 0711)
		if err != nil {
			t.Error("Could not create 'seite'")
		}
		t.Log("==> done creating 'seite'")
	} else {
		RemoveContents("seite")
	}
	t.Log("create a request")
	r, err := http.NewRequest("post", "/", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}
	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder and serve the Handler")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(startseite)
	handler.ServeHTTP(w, r)

	t.Log("Check the HTTP StatusCode. ... (Expecting Http OK)")
	if status := w.Code; status == 302 {
		t.Errorf("unexpected status code: %v instead of 200",
			status)
	}
}

func TestStartseiteAll(t *testing.T) {
	defer func() {
		t.Log("make sure folder 'seite' gets removed.")
		if _, err := os.Stat("seite"); err == nil {
			t.Log("seite does exist")
			var err = RemoveContents("seite")
			if err != nil {
				return
			}
			t.Log("==> done deleting file")
			t.Log("delete folder")
			err = os.Remove("seite")
			if err != nil {
				t.Error("Could not delete 'seite'")
			}
			t.Log("==> done deleting 'seite'")
		}

	}()

	t.Log("Test Startseite All. ... (Expecting http OK)")
	templateIndex = template.Must(template.ParseFiles("res/index.html"))

	t.Log("make sure folder 'seite' does exist.")
	if _, err := os.Stat("seite"); os.IsNotExist(err) {
		t.Log("seite does not exist")
		t.Log("create folder")
		var err = os.MkdirAll("seite", 0711)
		if err != nil {
			t.Error("Could not create 'seite'")
		}
		t.Log("==> done creating 'seite'")
		t.Log("Create file for the test. Will be deleted afterwards.")
		b := []byte(`{
				"Menu": null,
				"Dateiname": "",
				"Titel": "Testseite",
				"Inhalt": "Lorem ipsum dolor sit amet.",
				"Datum": "2017-10-31T21:51:53.15612987+01:00",
				"Autor": "Nils Holgerson",
				"Kommentare": [
					{
						"Autor": "Hallo Welt",
						"Datum": "2017-10-31T22:50:32.283095133+01:00",
						"Inhalt": "Programming is fun"
					}
				]
			}`)
		ioutil.WriteFile("seite/test.json", b, 0644)
	}
	t.Log("create a request")
	r, err := http.NewRequest("post", "/", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}

	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder and serve the Handler")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(startseite)
	handler.ServeHTTP(w, r)

	t.Log("Check the HTTP StatusCode. ... (Expecting Http OK)")
	if status := w.Code; status == 302 {
		t.Errorf("unexpected status code: %v instead of 200",
			status)
	}
}

func TestSeite(t *testing.T) {
	defer func() {
		t.Log("make sure folder 'seite' gets removed.")
		if _, err := os.Stat("seite"); err == nil {
			t.Log("seite does exist")
			var err = RemoveContents("seite")
			if err != nil {
				return
			}
			t.Log("==> done deleting file")
			t.Log("delete folder")
			err = os.Remove("seite")
			if err != nil {
				t.Error("Could not delete 'seite'")
			}
			t.Log("==> done deleting 'seite'")
		}

	}()

	t.Log("Test seite. ... (Expecting http OK)")
	templateSeite = template.Must(template.ParseFiles("res/template.html"))

	t.Log("make sure folder 'seite' does exist.")
	if _, err := os.Stat("seite"); os.IsNotExist(err) {
		t.Log("seite does not exist")
		t.Log("create folder")
		var err = os.MkdirAll("seite", 0711)
		if err != nil {
			t.Error("Could not create 'seite'")
		}
		t.Log("==> done creating 'seite'")
		t.Log("Create file for the test. Will be deleted afterwards.")
		b := []byte(`{
				"Menu": null,
				"Dateiname": "",
				"Titel": "Testseite",
				"Inhalt": "Lorem ipsum dolor sit amet.",
				"Datum": "2017-10-31T21:51:53.15612987+01:00",
				"Autor": "Nils Holgerson",
				"Kommentare": [
					{
						"Autor": "Hallo Welt",
						"Datum": "2017-10-31T22:50:32.283095133+01:00",
						"Inhalt": "Programming is fun"
					}
				]
			}`)
		ioutil.WriteFile("seite/test.json", b, 0644)
	}
	t.Log("create a request")
	r, err := http.NewRequest("post", "/seite/test", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}

	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder and serve the Handler")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(seite)
	handler.ServeHTTP(w, r)

	t.Log("Check the HTTP StatusCode. ... (Expecting Http OK)")
	if status := w.Code; status == 302 {
		t.Errorf("unexpected status code: %v instead of 200",
			status)
	}
}

/* login  http template method */
func TestLogin(t *testing.T) {
	t.Log("Test with a Testlogin. ... (Expecting Redirect)")
	t.Log("create a request")
	r, err := http.NewRequest("post", "/login", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}

	q := r.URL.Query()
	q.Add("name", "Max Mustermann")
	q.Add("pass", "1234567890")
	r.URL.RawQuery = q.Encode()

	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder and serve the Handler")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(login)
	handler.ServeHTTP(w, r)

	t.Log("Check the HTTP StatusCode. ... (Expecting 302)")
	if status := w.Code; status != 302 {
		t.Errorf("unexpected status code: %v instead of %v",
			status, 302)
	}
}

func TestLoginEmptyCredentials(t *testing.T) {
	t.Log("Test with empty credentials. ... (Expecting no redirect)")
	templateLogin = template.Must(template.ParseFiles("res/login.html"))
	t.Log("create a request")
	r, err := http.NewRequest("post", "/login", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}

	q := r.URL.Query()
	q.Add("name", "")
	q.Add("pass", "")
	r.URL.RawQuery = q.Encode()

	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder and serve the Handler")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(login)
	handler.ServeHTTP(w, r)

	t.Log("Check the HTTP StatusCode. ... (Expecting no redirect)")
	if status := w.Code; status == 302 {
		t.Errorf("unexpected status code: %v instead of 200",
			status)
	}
}

func TestLoginWrongData(t *testing.T) {
	t.Log("Test with a Wrong Testlogin. ... (Expecting no redirect)")
	templateLogin = template.Must(template.ParseFiles("res/login.html"))
	t.Log("create a request")
	r, err := http.NewRequest("post", "/login", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}

	q := r.URL.Query()
	q.Add("name", "Max Mustermann")
	q.Add("pass", "wrong")
	r.URL.RawQuery = q.Encode()

	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder and serve the Handler")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(login)
	handler.ServeHTTP(w, r)

	t.Log("Check the HTTP StatusCode. ... (Expecting o redirect)")
	if status := w.Code; status == 302 {
		t.Errorf("unexpected status code: %v instead of 200",
			status)
	}
}

/* logout  http template method */

func TestLogout(t *testing.T) {
	t.Log("Test logout. ... (Expecting Redirect)")
	t.Log("create a request")
	r, err := http.NewRequest("post", "/logout", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}

	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder and serve the Handler")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(logout)
	handler.ServeHTTP(w, r)

	t.Log("Check the HTTP StatusCode. ... (Expecting 302)")
	if status := w.Code; status != 302 {
		t.Errorf("unexpected status code: %v instead of %v",
			status, 302)
	}
}

/* test if comment is a duplicate. Will fail because of no comments*/

func TestEnthaeltNoComments(t *testing.T) {
	var commentsList []Kommentar
	newComment := Kommentar{Autor: "AUTHOR NAME", Inhalt: "INHALT", Datum: time.Now()}
	result := enthaelt(commentsList, newComment)
	if result == false {
		t.Log("no match found as expected")
	} else {
		t.Error("Match Found!")
	}
}

/* test if comment is a duplicate. */

func TestEnthaeltduplicateComments(t *testing.T) {
	var commentsList []Kommentar
	k1 := Kommentar{Autor: "AUTHOR NAME", Inhalt: "INHALT X", Datum: time.Now()}
	k2 := Kommentar{Autor: "AUTHOR NAME", Inhalt: "INHALT Y", Datum: time.Now()}
	k3 := Kommentar{Autor: "AUTHOR NAME", Inhalt: "INHALT Z", Datum: time.Now()}
	commentsList = append([]Kommentar{k1}, commentsList...)
	commentsList = append([]Kommentar{k2}, commentsList...)
	commentsList = append([]Kommentar{k3}, commentsList...)

	newComment := Kommentar{Autor: "AUTHOR NAME", Inhalt: "INHALT X", Datum: time.Now()}
	result := enthaelt(commentsList, newComment)
	if result != false {
		t.Log("Match found as expected")
	} else {
		t.Error("No Match Found!")
	}
}

/* writes a comment to the file system.*/
func TestErstelleKommentarFailure(t *testing.T) {
	t.Log("Test creation of a comment. ... (Expecting it gets not created)")
	t.Log("create a request")
	r, err := http.NewRequest("post", "/hjkl", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}

	q := r.URL.Query()
	q.Add("autor", "Max Mustermann")
	q.Add("inhalt", "Einfach ein bisschen Inhalt")
	r.URL.RawQuery = q.Encode()

	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder and serve the Handler")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(erstelleKommentar)
	handler.ServeHTTP(w, r)

	t.Log("Check the HTTP Status Code. ... (Expecting status OK)")
	if status := w.Code; status != http.StatusOK {
		t.Errorf("unexpected status code: %v instead of %v",
			status, http.StatusOK)
	}
}

func TestErstelleKommentar(t *testing.T) {
	defer func() {
		var err = os.Remove("test.json")
		if err != nil {
			return
		}
		t.Log("==> done deleting file")
	}()
	t.Log("Test creation of a comment. ... (Expecting it gets created)")
	t.Log("Create file for the test. Will be deleted afterwards.")
	b := []byte(`{
	"Menu": null,
	"Dateiname": "",
	"Titel": "Testseite",
	"Inhalt": "Lorem ipsum dolor sit amet.",
	"Datum": "2017-10-31T21:51:53.15612987+01:00",
	"Autor": "Nils Holgerson",
	"Kommentare": [
		{
			"Autor": "Hallo Welt",
			"Datum": "2017-10-31T22:50:32.283095133+01:00",
			"Inhalt": "Programming is fun"
		}
	]
}`)
	ioutil.WriteFile("test.json", b, 0644)
	t.Log("create a request")
	r, err := http.NewRequest("post", "/test", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}

	q := r.URL.Query()
	q.Add("autor", "Max Mustermann")
	q.Add("inhalt", "Einfach ein bisschen Inhalt")
	r.URL.RawQuery = q.Encode()

	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder and serve the Handler")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(erstelleKommentar)
	handler.ServeHTTP(w, r)

	t.Log("Check the HTTP Status Code. ... (Expecting status OK)")
	if status := w.Code; status != http.StatusOK {
		t.Errorf("unexpected status code: %v instead of %v",
			status, http.StatusOK)
	}
}

func TestErstelleKommentarDuplicate(t *testing.T) {
	defer func() {
		var err = os.Remove("test.json")
		if err != nil {
			return
		}
		t.Log("==> done deleting file")
	}()
	t.Log("Test creation of a comment. ... (Expecting it gets not created)")
	t.Log("Create file for the test. Will be deleted afterwards.")
	b := []byte(`{
	"Menu": null,
	"Dateiname": "",
	"Titel": "Testseite",
	"Inhalt": "Lorem ipsum dolor sit amet.",
	"Datum": "2017-10-31T21:51:53.15612987+01:00",
	"Autor": "Nils Holgerson",
	"Kommentare": [
		{
			"Autor": "Hallo Welt",
			"Datum": "2017-10-31T22:50:32.283095133+01:00",
			"Inhalt": "Programming is fun"
		}
	]
}`)
	ioutil.WriteFile("test.json", b, 0644)
	t.Log("create a request")
	r, err := http.NewRequest("post", "/test", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}

	q := r.URL.Query()
	q.Add("autor", "Hallo Welt")
	q.Add("inhalt", "Programming is fun")
	r.URL.RawQuery = q.Encode()

	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder and serve the Handler")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(erstelleKommentar)
	handler.ServeHTTP(w, r)

	t.Log("Check the HTTP Status Code. ... (Expecting status OK)")
	if status := w.Code; status != http.StatusOK {
		t.Errorf("unexpected status code: %v instead of %v",
			status, http.StatusOK)
	}
}

func TestErstelleKommentarEmptyData(t *testing.T) {
	t.Log("Test creation of a comment. ... (Expecting it gets created)")
	t.Log("create a request")
	r, err := http.NewRequest("post", "/kommentar", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}

	q := r.URL.Query()
	q.Add("autor", "")
	q.Add("inhalt", "")
	r.URL.RawQuery = q.Encode()

	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder and serve the Handler")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(erstelleKommentar)
	handler.ServeHTTP(w, r)

	t.Log("Check the HTTP Status Code. ... (Expecting status OK)")
	if status := w.Code; status != http.StatusOK {
		t.Errorf("unexpected status code: %v instead of 200",
			status)
	}
}

/* profil */
func TestProfilNoLogin(t *testing.T) {
	t.Log("Test getting profile page with no Login. ... (Expecting Http redirect)")
	t.Log("create a request")
	r, err := http.NewRequest("post", "/profil", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}

	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder and serve the Handler")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(profil)
	handler.ServeHTTP(w, r)

	t.Log("Check the HTTP Status Code. ... (Expecting status redirect)")
	if status := w.Code; status != 302 {
		t.Errorf("unexpected status code: %v instead of 302",
			status)
	}
}
func TestProfilActiveLogin(t *testing.T) {
	t.Log("Test getting profile page with active Login. ... (Expecting Http OK)")
	templateProfil = template.Must(template.ParseFiles("res/profil.html"))
	t.Log("Create a request in the running session")
	t.Log(sitzungen)
	session := sitzungen[0]
	cookie := session.Keks
	r, err := http.NewRequest("post", "/profil", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}
	r.AddCookie(&cookie)

	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder and serve the Handler")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(profil)
	handler.ServeHTTP(w, r)

	t.Log("Check the HTTP Status Code. ... (Expecting status OK)")
	if status := w.Code; status != 200 {
		t.Errorf("unexpected status code: %v instead of 200",
			status)
	}
}

/* password */

func TestPasswortNoLogin(t *testing.T) {
	t.Log("Test passwort change. ... (Expecting Http redirect)")
	templateProfil = template.Must(template.ParseFiles("res/profil.html"))
	t.Log("No request in a running session")
	r, err := http.NewRequest("post", "/passwort", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("add form data")
	q := r.URL.Query()
	q.Add("neu_pass", "asdf")
	q.Add("pass_wdh", "asdf")
	q.Add("alt_pass", "1234567890")
	r.URL.RawQuery = q.Encode()
	t.Log(r.FormValue("neu_pass"))
	t.Log(r.FormValue("pass_wdh"))
	t.Log(r.FormValue("alt_pass"))
	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder and serve the Handler")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(passwort)
	handler.ServeHTTP(w, r)

	t.Log("Check the HTTP Status Code. ... (Expecting status 302)")
	if status := w.Code; status != 302 {
		t.Errorf("unexpected status code: %v instead of 302",
			status)
	}

}

func TestPasswortActiveLoginInCorrect(t *testing.T) {
	t.Log("Test passwort change. ... (Expecting Http OK)")
	templateProfil = template.Must(template.ParseFiles("res/profil.html"))
	t.Log("Create a request in a running session")
	t.Log(sitzungen)
	session := sitzungen[0]
	cookie := session.Keks
	r, err := http.NewRequest("post", "/passwort", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}
	r.AddCookie(&cookie)
	t.Log("add form data")
	q := r.URL.Query()
	q.Add("neu_pass", "asdf")
	q.Add("pass_wdh", "asdf")
	q.Add("alt_pass", "qwerty")
	r.URL.RawQuery = q.Encode()
	t.Log(r.FormValue("neu_pass"))
	t.Log(r.FormValue("pass_wdh"))
	t.Log(r.FormValue("alt_pass"))

	p := gebeProfil("Max Mustermann")
	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder and serve the Handler")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(passwort)
	handler.ServeHTTP(w, r)

	t.Log("Check the HTTP Status Code. ... (Expecting status OK)")
	if status := w.Code; status != 200 {
		t.Errorf("unexpected status code: %v instead of 200",
			status)
	}
	t.Log("Check the Meldung Code. ... (Expecting 'Das Passwort ist falsch')")
	if note := p.Meldung; note != "Das Passwort ist falsch" {
		t.Errorf("unexpected status code: %v instead of 'Das Passwort ist falsch'",
			note)
	}
}

func TestPasswortActiveLoginDifferentPWD(t *testing.T) {
	t.Log("Test passwort change. ... (Expecting Http OK)")
	templateProfil = template.Must(template.ParseFiles("res/profil.html"))
	t.Log("Create a request in a running session")
	t.Log(sitzungen)
	session := sitzungen[0]
	cookie := session.Keks
	r, err := http.NewRequest("post", "/passwort", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}
	r.AddCookie(&cookie)
	t.Log("add form data")
	q := r.URL.Query()
	q.Add("neu_pass", "asdf")
	q.Add("pass_wdh", "qwerty")
	q.Add("alt_pass", "1234567890")
	r.URL.RawQuery = q.Encode()
	t.Log(r.FormValue("neu_pass"))
	t.Log(r.FormValue("pass_wdh"))
	t.Log(r.FormValue("alt_pass"))

	p := gebeProfil("Max Mustermann")
	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder and serve the Handler")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(passwort)
	handler.ServeHTTP(w, r)

	t.Log("Check the HTTP Status Code. ... (Expecting status OK)")
	if status := w.Code; status != 200 {
		t.Errorf("unexpected status code: %v instead of 200",
			status)
	}
	t.Log("Check the Meldung Code. ... (Expecting 'Die Passwörter stimmen nicht überein')")
	if note := p.Meldung; note != "Die Passwörter stimmen nicht überein" {
		t.Errorf("unexpected status code: %v instead of 'Die Passwörter stimmen nicht überein'",
			note)
	}
}

func TestPasswortActiveLoginCorrect(t *testing.T) {
	defer func() {
		t.Log("remove user from user.json after password change.")
		err := ioutil.WriteFile("user.json", []byte(`""`), 0644)
		if err != nil {
			t.Error(err)
		}
	}()
	t.Log("Test passwort change. ... (Expecting Http OK)")
	templateProfil = template.Must(template.ParseFiles("res/profil.html"))
	t.Log("Create a request in a running session")
	t.Log(sitzungen)
	session := sitzungen[0]
	cookie := session.Keks
	r, err := http.NewRequest("post", "/neu", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}
	r.AddCookie(&cookie)
	t.Log("add form data")
	q := r.URL.Query()
	q.Add("neu_pass", "asdf")
	q.Add("pass_wdh", "asdf")
	q.Add("alt_pass", "1234567890")
	r.URL.RawQuery = q.Encode()
	t.Log(r.FormValue("neu_pass"))
	t.Log(r.FormValue("pass_wdh"))
	t.Log(r.FormValue("alt_pass"))

	p := gebeProfil("Max Mustermann")
	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder and serve the Handler")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(passwort)
	handler.ServeHTTP(w, r)

	t.Log("Check the HTTP Status Code. ... (Expecting status OK)")
	if status := w.Code; status != 200 {
		t.Errorf("unexpected status code: %v instead of 200",
			status)
	}
	t.Log("Check the Meldung Code. ... (Expecting 'Das Passwort wurde geändert')")
	if note := p.Meldung; note != "Das Passwort wurde geändert" {
		t.Errorf("unexpected status code: %v instead of 'Das Passwort wurde geändert'",
			note)
	}
}

/* erstelleSeite */

/*neu Create a new page. Not possible without user created on cli.*/

func TestNeuEmpty(t *testing.T) {
	defer RemoveContents("seite")
	t.Log("Test create Neu. ... (Expecting Http OK)")
	templateNeu = template.Must(template.ParseFiles("res/neu.html"))
	t.Log("Create a request in a running session")
	t.Log(sitzungen)
	session := sitzungen[0]
	cookie := session.Keks
	r, err := http.NewRequest("post", "/neu", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}
	r.AddCookie(&cookie)
	t.Log("add form data")
	q := r.URL.Query()
	q.Add("titel", "")
	q.Add("inhalt", "")
	r.URL.RawQuery = q.Encode()
	t.Log(r.FormValue("titel"))
	t.Log(r.FormValue("inhalt"))
	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder and serve the Handler")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(neu)
	handler.ServeHTTP(w, r)

	t.Log("Check the HTTP Status Code. ... (Expecting status OK)")
	if status := w.Code; status != 200 {
		t.Errorf("unexpected status code: %v instead of 200",
			status)
	}
}

func TestNeuNoSession(t *testing.T) {
	defer RemoveContents("seite")
	t.Log("Test create Neu. ... (Expecting Http Redirect)")
	templateNeu = template.Must(template.ParseFiles("res/neu.html"))
	t.Log("Create a request outside a running session")
	r, err := http.NewRequest("post", "/neu", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("add form data")
	q := r.URL.Query()
	q.Add("titel", "")
	q.Add("inhalt", "")
	r.URL.RawQuery = q.Encode()
	t.Log(r.FormValue("titel"))
	t.Log(r.FormValue("inhalt"))
	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder and serve the Handler")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(neu)
	handler.ServeHTTP(w, r)

	t.Log("Check the HTTP Status Code. ... (Expecting status 302)")
	if status := w.Code; status != 302 {
		t.Errorf("unexpected status code: %v instead of 302",
			status)
	}
}
func TestNeu(t *testing.T) {
	t.Log("Test create Neu. ... (Expecting creation -> 302)")
	templateNeu = template.Must(template.ParseFiles("res/neu.html"))
	t.Log("Create a request in a running session")
	t.Log(sitzungen)
	session := sitzungen[0]
	cookie := session.Keks
	r, err := http.NewRequest("post", "/neu", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}
	r.AddCookie(&cookie)
	t.Log("add form data")
	q := r.URL.Query()
	q.Add("titel", "asdf")
	q.Add("inhalt", "asdf")
	r.URL.RawQuery = q.Encode()
	t.Log(r.FormValue("titel"))
	t.Log(r.FormValue("inhalt"))
	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder and serve the Handler")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(neu)
	handler.ServeHTTP(w, r)

	t.Log("Check the HTTP Status Code. ... (Expecting status 302)")
	if status := w.Code; status != 302 {
		t.Errorf("unexpected status code: %v instead of 302",
			status)
	}
}

/*bearbeiten Alter an existing page.*/
func TestBearbeitenNoSession(t *testing.T) {
	t.Log("Test Bearbeiten. ... (Expecting Http redirect)")
	templateNeu = template.Must(template.ParseFiles("res/neu.html"))
	t.Log("get a site UUID")
	d, err := os.Open("seite")
	if err != nil {
		t.Fatalf("could not get directory 'seite': %v", err)
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		t.Fatalf("could not get directory 'seite': %v", err)
	}
	sitename := strings.TrimRight(names[0], ".json")
	t.Log("Create a request outside a running session")
	r, err := http.NewRequest("post", "/bearbeiten/"+sitename, nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("add form data")
	q := r.URL.Query()
	q.Add("titel", "")
	q.Add("inhalt", "")
	r.URL.RawQuery = q.Encode()
	t.Log(r.FormValue("titel"))
	t.Log(r.FormValue("inhalt"))
	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder and serve the Handler")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(bearbeiten)
	handler.ServeHTTP(w, r)

	t.Log("Check the HTTP Status Code. ... (Expecting status 302)")
	if status := w.Code; status != 302 {
		t.Errorf("unexpected status code: %v instead of 302",
			status)
	}
}

func TestBearbeiten(t *testing.T) {
	t.Log("Test Bearbeiten. ... (Expecting Http OK)")
	templateNeu = template.Must(template.ParseFiles("res/neu.html"))
	t.Log("get a site UUID")
	d, err := os.Open("seite")
	if err != nil {
		t.Fatalf("could not get directory 'seite': %v", err)
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		t.Fatalf("could not get directory 'seite': %v", err)
	}
	sitename := strings.TrimRight(names[0], ".json")
	t.Log("Create a request outside a running session")
	r, err := http.NewRequest("post", "/bearbeiten/"+sitename, nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sitzungen)
	session := sitzungen[0]
	cookie := session.Keks
	r.AddCookie(&cookie)
	t.Log("add form data")
	q := r.URL.Query()
	q.Add("titel", "")
	q.Add("inhalt", "")
	r.URL.RawQuery = q.Encode()
	t.Log(r.FormValue("titel"))
	t.Log(r.FormValue("inhalt"))
	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder and serve the Handler")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(bearbeiten)
	handler.ServeHTTP(w, r)

	t.Log("Check the HTTP Status Code. ... (Expecting status 200)")
	if status := w.Code; status != 200 {
		t.Errorf("unexpected status code: %v instead of 200",
			status)
	}
}

/* bestaetigen */

func TestBestaetigen(t *testing.T) {
	t.Log("Test Bearbeiten. ... (Expecting Http OK)")
	templateLoeschen = template.Must(template.ParseFiles("res/loeschen.html"))
	t.Log("get a site UUID")
	d, err := os.Open("seite")
	if err != nil {
		t.Fatalf("could not get directory 'seite': %v", err)
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		t.Fatalf("could not get directory 'seite': %v", err)
	}
	t.Log(names)
	sitename := strings.TrimRight(names[0], ".json")
	t.Log("Create a request outside a running session")
	r, err := http.NewRequest("post", "/bestaetigen/"+sitename, nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sitzungen)
	session := sitzungen[0]
	cookie := session.Keks
	r.AddCookie(&cookie)
	t.Log("add form data")
	q := r.URL.Query()
	q.Add("titel", "")
	q.Add("inhalt", "")
	r.URL.RawQuery = q.Encode()
	t.Log(r.FormValue("titel"))
	t.Log(r.FormValue("inhalt"))
	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder and serve the Handler")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(bestaetigen)
	handler.ServeHTTP(w, r)

	t.Log("Check the HTTP Status Code. ... (Expecting status OK)")
	if status := w.Code; status != 200 {
		t.Errorf("unexpected status code: %v instead of 200",
			status)
	}
}

/*loeschen Delete an existing page. Not possible without any created page.*/

func TestLoeschenNoSession(t *testing.T) {
	t.Log("Test Loeschen. ... (Expecting Http Redirect)")
	templateNeu = template.Must(template.ParseFiles("res/neu.html"))
	t.Log("get a site UUID")
	d, err := os.Open("seite")
	if err != nil {
		t.Fatalf("could not get directory 'seite': %v", err)
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		t.Fatalf("could not get directory 'seite': %v", err)
	}
	sitename := strings.TrimRight(names[0], ".json")
	t.Log("Create a request outside a running session")
	r, err := http.NewRequest("post", "/loeschen/"+sitename, nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("add form data")
	q := r.URL.Query()
	q.Add("titel", "")
	q.Add("inhalt", "")
	r.URL.RawQuery = q.Encode()
	t.Log(r.FormValue("titel"))
	t.Log(r.FormValue("inhalt"))
	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder and serve the Handler")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(loeschen)
	handler.ServeHTTP(w, r)

	t.Log("Check the HTTP Status Code. ... (Expecting status 302)")
	if status := w.Code; status != 302 {
		t.Errorf("unexpected status code: %v instead of 302",
			status)
	}
}

func TestLoeschen(t *testing.T) {
	t.Log("Test Bearbeiten. ... (Expecting Http 302)")
	templateNeu = template.Must(template.ParseFiles("res/neu.html"))
	t.Log("get a site UUID")
	d, err := os.Open("seite")
	if err != nil {
		t.Fatalf("could not get directory 'seite': %v", err)
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		t.Fatalf("could not get directory 'seite': %v", err)
	}
	sitename := strings.TrimRight(names[0], ".json")
	t.Log("Create a request outside a running session")
	r, err := http.NewRequest("post", "/loeschen/"+sitename, nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sitzungen)
	session := sitzungen[0]
	cookie := session.Keks
	r.AddCookie(&cookie)
	t.Log("add form data")
	q := r.URL.Query()
	q.Add("titel", "")
	q.Add("inhalt", "")
	r.URL.RawQuery = q.Encode()
	t.Log(r.FormValue("titel"))
	t.Log(r.FormValue("inhalt"))
	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder and serve the Handler")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(loeschen)
	handler.ServeHTTP(w, r)

	t.Log("Check the HTTP Status Code. ... (Expecting status 302)")
	if status := w.Code; status != 302 {
		t.Errorf("unexpected status code: %v instead of 302",
			status)
	}
}

/* main */

/* erstelleNutzer is a command line function. Not tested, because of missing input :( */
