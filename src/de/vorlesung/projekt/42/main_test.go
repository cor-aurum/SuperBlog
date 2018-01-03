package main

/* Autoren: 3818468, 6985153, 9875672 */

/*
Wegen vielen in unserem Design hardcoded eingebauten Pfaden oder Überprüfungen,
sind viele Tests nicht oder nur mit einem erwarteten Fehler möglich. Die Vorles-
ung zur Nutzung von übergebenen Providern kam dafür zu spät.
*/

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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

/* searches for an profile by name. Will fail because of no profiles loaded. */

func TestGebeProfilNotExisting(t *testing.T) {
	t.Log("try to get a profile. ... (Expect nil return)")
	name := "Max Mustermann"
	t.Log("mockup name: ", name)
	result := gebeProfil(name)
	if result == nil {
		t.Log("no profile as expected")
	} else {
		t.Error("returned a non existing profile")
	}
}

func TestGebeProfilExisting(t *testing.T) {
	t.Log("adding a loaded user to global variable.")
	name := "Max Mustermann\n"
	pass := "1234567890\n"
	t.Log("mockup name: ", name)
	t.Log("mockup pass: ", pass)
	appendUser(name, pass)
	fmt.Println(profile.Profile)

	t.Log("try to get a profile. ... (Expect the profile returned)")
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
	name := "Max Mustermann"
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
	t.Log("adding a loaded user to global variable.")
	name := "Max Mustermann\n"
	pass := "1234567890\n"
	t.Log("mockup name: ", name)
	t.Log("mockup pass: ", pass)
	appendUser(name,pass)

	t.Log("try to login with a wrong password. ... (Expect login failure)")
	wrong_pass := "asdfqwerty"
	t.Log("tryout name: ", name)
	t.Log("tryout pass: ", wrong_pass)
	result := pruefeLogin(name, SalzHash(name, pass))
	if result == true {
		t.Error("logged in with wrong user credentials!")
	} else {
		t.Log("login failed as expected")
	}
}

/* checks if username and password are matching correctly. Expect success*/
func TestPruefeLoginSuccessful(t *testing.T) {
	t.Log("adding a loaded user to global variable.")
	name := "Max Mustermann\n"
	pass := "1234567890\n"
	t.Log("mockup name: ", name)
	t.Log("mockup pass: ", pass)
	appendUser(name, pass)

	t.Log("try to login with real credentials. ... (Expect login success)")
	t.Log("tryout name: ", name)
	t.Log("tryout pass: ", pass)
	result := pruefeLogin(name, SalzHash(name, pass))
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

	t.Log("Check the HTTP Staus Code. ... (Expecting OK)")
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
	t.Log(login)
	t.Log(erfolg)
}

func TestMachLoginRealData(t *testing.T) {
	t.Log("adding a loaded user to global variable.")
	name := "Max Mustermann\n"
	pass := "1234567890\n"
	t.Log("mockup name: ", name)
	t.Log("mockup pass: ", pass)
	appendUser(name, pass)

	t.Log("Test with a Testlogin. ... (Expecting (true, true))")
	t.Log("create a request")
	r, err := http.NewRequest("post", "Login", strings.NewReader("{name: 'Max Mustermann', pass: 1234567890}"))
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

	t.Log("Check the HTTP Staus Code. ... (Expecting OK)")
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

/* returns a session with an specific id. Will fail because no sessions are running*/

func TestGebeSitzung(t *testing.T) {

}

/* deletes a session with an specifiic id. Will fail because no sessions are running*/

func TestLoescheSitzung(t *testing.T) {

}

/* creates an array of the menu items */
func TestMachMenu(t *testing.T) {

}

/* creates a cookie */
func TestKekse(t *testing.T) {

}

/* kommentarKeks void method. no test */

/* startseite  http template method ->  no test */

/* login  http template method ->  no test */

/* logout  http template method ->  no test */

/* test if comment is a duplicate. Will fail because of no comments*/

func TestEnthaeltNoComments(t *testing.T) {

}

/* writes a comment to the file system. Even possible without a user*/

func TestErstelleKommentar(t *testing.T) {

}

/* test if comment is a duplicate. */

func TestEnthaeltduplicateComments(t *testing.T) {

}

/* erstelleNutzer is a command line function. Not tested, because of missing input :( */

/* erstelleVerzeichnis File Sytem function. Not tested */

/*neu Create a new page. Not possible without user created on cli.*/

/*bearbeiten Alter an existing page. Not possible without any created page.*/

/*loeschen Delete an existing page. Not possible without any created page.*/

/* main Is not tested because of using hardcoded paths. */
