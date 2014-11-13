package gaestebin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"appengine"
	"appengine/aetest"
	"appengine/datastore"
	"appengine/user"
)

func TestGetPaste(t *testing.T) {
	inst, err := aetest.NewInstance(nil)
	if err != nil {
		t.Errorf("Failed to create instance: %v", err)
	}
	defer inst.Close()

	id := GenerateRandomString(8)

	r, err := inst.NewRequest("GET", "/paste/v1/"+id, nil)
	if err != nil {
		t.Errorf("Failed to create req: %v", err)
	}
	c := appengine.NewContext(r)

	p := Paste{Id: id, Email: "name@domain.com"}

	key := datastore.NewKey(c, "Paste", p.Id, 0, nil)
	if _, err := datastore.Put(c, key, &p); err != nil {
		t.Fatal(err)
	}

	u := user.User{Email: "name@domain.com"}
	aetest.Login(&u, r)

	w := httptest.NewRecorder()

	Router().ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fail()
	}
}

func TestCreatePaste(t *testing.T) {
	inst, err := aetest.NewInstance(nil)
	if err != nil {
		t.Fatalf("Failed to create instance: %v", err)
	}
	defer inst.Close()

	body := strings.NewReader(`{"Title":"","Content":"","Language":""}`)

	r, err := inst.NewRequest("POST", "/paste/v1", body)
	if err != nil {
		t.Fatalf("Failed to create req: %v", err)
	}

	u := user.User{Email: "name@domain.com"}
	aetest.Login(&u, r)

	w := httptest.NewRecorder()

	Router().ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fail()
	}

	var p Paste
	err = json.Unmarshal(w.Body.Bytes(), &p)
	if err != nil {
		t.Fail()
	}
}

func TestDeletePaste(t *testing.T) {
	inst, err := aetest.NewInstance(nil)
	if err != nil {
		t.Fatalf("Failed to create instance: %v", err)
	}
	defer inst.Close()

	id := GenerateRandomString(8)

	r, err := inst.NewRequest("DELETE", "/paste/v1/"+id, nil)
	if err != nil {
		t.Fatalf("Failed to create req: %v", err)
	}
	c := appengine.NewContext(r)

	p := Paste{Id: id, Email: "name@domain.com"}

	key := datastore.NewKey(c, "Paste", p.Id, 0, nil)
	if _, err := datastore.Put(c, key, &p); err != nil {
		t.Fatal(err)
	}

	u := user.User{Email: "name@domain.com"}
	aetest.Login(&u, r)

	w := httptest.NewRecorder()

	Router().ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fail()
	}

	err = datastore.Get(c, key, &p)
	if err != nil && err != datastore.ErrNoSuchEntity {
		t.Fail()
	}
}
