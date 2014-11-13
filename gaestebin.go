package gaestebin

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"appengine/user"

	"github.com/gorilla/mux"
)

type Paste struct {
	Id        string    `datastore:"id" json:"id"`
	Timestamp time.Time `datastore:"timestamp" json:"timestamp"`
	Content   string    `datastore:"content,noindex" json:"content"`
	Email     string    `datastore:"email" json:"email"`
	Title     string    `datastore:"title" json:"title"`
	Language  string    `datastore:"language" json:"language"`
	IsOwner   bool      `datastore:"-" json:"isOwner"`
}

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "=")
}

func GetPaste(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		c.Infof("%v Login required", appengine.RequestID(c))
		w.WriteHeader(http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	paste := Paste{}
	if _, err := memcache.JSON.Get(c, id, &paste); err == memcache.ErrCacheMiss {
		key := datastore.NewKey(c, "Paste", id, 0, nil)
		if err := datastore.Get(c, key, &paste); err != nil {
			c.Infof(err.Error())
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		item := &memcache.Item{Key: id, Object: paste}
		memcache.JSON.Set(c, item)
	}
	paste.IsOwner = (paste.Email == u.Email)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&paste); err != nil {
		c.Infof(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func CreatePaste(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		c.Infof("%v Login required", appengine.RequestID(c))
		w.WriteHeader(http.StatusForbidden)
		return
	}

	paste := Paste{}
	if err := json.NewDecoder(r.Body).Decode(&paste); err != nil {
		c.Infof(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	paste.Id = GenerateRandomString(8)
	paste.Timestamp = time.Now()
	paste.Email = u.Email
	paste.IsOwner = (paste.Email == u.Email)

	key := datastore.NewKey(c, "Paste", paste.Id, 0, nil)
	if _, err := datastore.Put(c, key, &paste); err != nil {
		c.Infof(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	item := &memcache.Item{Key: paste.Id, Object: paste}
	memcache.JSON.Set(c, item)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&paste); err != nil {
		c.Infof(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Infof("Created new paste")
}

func DeletePaste(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		c.Infof("%v Login required", appengine.RequestID(c))
		w.WriteHeader(http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	key := datastore.NewKey(c, "Paste", id, 0, nil)
	paste := Paste{}
	if err := datastore.Get(c, key, &paste); err != nil {
		c.Infof(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if paste.Email != u.Email {
		c.Infof("Bad owner")
		http.Error(w, "Bad Owner", http.StatusForbidden)
		return
	}

	if err := datastore.Delete(c, key); err != nil {
		c.Infof(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	memcache.Delete(c, id)
}

func Router() *mux.Router {
	r := mux.NewRouter()

	s := r.PathPrefix("/api").Subrouter()
	s.HandleFunc("/paste/{id:[a-zA-Z0-9_-]+}", GetPaste).Methods("GET")
	s.HandleFunc("/paste", CreatePaste).Methods("POST")
	s.HandleFunc("/paste/{id:[a-zA-Z0-9_-]+}", DeletePaste).Methods("DELETE")

	return r
}

func init() {
	http.Handle("/", Router())
}
