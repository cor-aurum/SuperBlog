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
	"testing"
	"time"
)

func TestCreateTestUser(t *testing.T) {
	t.Log("adding a loaded user to global variable.")
	name := "Max Mustermann\n"
	pass := "1234567890\n"
	t.Log("mockup name: ", name)
	t.Log("mockup pass: ", pass)
	appendUser(name, pass)
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

func TestGebeSitzung(t *testing.T) {
	t.Log("reading Session variable. ... (Expecting (false, \"\"))")
	t.Log("create a request")
	r, err := http.NewRequest("get", "", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}

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

	if sessionExists == false && len(name) <= 0 {
		t.Log("no session running as expected")
	} else {
		t.Errorf("there is a Session running! : %v, %v ", sessionExists, name)
	}
}

/* loescheSitzung deletes a session with an specifiic id. Void Method no test*/

/* creates an array of the menu items */
func TestMachMenu(t *testing.T) {

}

/* creates a cookie */
func TestKekse(t *testing.T) {

}

/* kommentarKeks void method. no test */

/* startseite  http template method ->  no test */

func TestStartseiteWithoutData(t *testing.T) {
	t.Log("Test with empty credentials. ... (Expecting no redirect)")
	templateIndex = template.Must(template.ParseFiles("res/index.html"))
	t.Log("create a request")
	r, err := http.NewRequest("post", "/login", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; params=value")
	if err != nil {
		t.Fatal(err)
	}

	/*  ResponseRecorder instead of ResponseWriter to record the response. */
	t.Log("create a ResponseRecorder and serve the Handler")
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(startseite)
	handler.ServeHTTP(w, r)

	t.Log("Check the HTTP StatusCode. ... (Expecting no redirect)")
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

/* logout  http template method ->  no test */

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
	r, err := http.NewRequest("post", "hjkl.json", nil)
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
	b := []byte(`{
	"Menu": null,
	"Dateiname": "",
	"Titel": "Testseite",
	"Inhalt": "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.",
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
	t.Log("Test creation of a comment. ... (Expecting it gets created)")
	b := []byte(`{
	"Menu": null,
	"Dateiname": "",
	"Titel": "Testseite",
	"Inhalt": "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.",
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

/* erstelleNutzer is a command line function. Not tested, because of missing input :( */

/*neu Create a new page. Not possible without user created on cli.*/

/*bearbeiten Alter an existing page. Not possible without any created page.*/

/*loeschen Delete an existing page. Not possible without any created page.*/

/* main Is not tested because of using hardcoded paths. */
