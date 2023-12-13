package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"sync"
)

var (
	listUserRe   = regexp.MustCompile(`^\/users[\/]*$`)
	getUserRe    = regexp.MustCompile(`^\/users\/(\d+)$`)
	createUserRe = regexp.MustCompile(`^\/users[\/]*$`)
)

type UserHandler struct {
	store *DataStore
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DataStore struct {
	m map[string]User
	*sync.RWMutex
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	h.store.RLock()
	users := make([]User, 0, len(h.store.m))
	for _, v := range h.store.m {
		users = append(users, v)
	}

	h.store.RUnlock()
	jsonBytes, err := json.Marshal(users)
	if err != nil {
		internalServerError(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	matches := getUserRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		notFound(w, r)
		return
	}

	h.store.RLock()
	u, ok := h.store.m[matches[1]]
	h.store.RUnlock()

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("USER NOT FOUND"))
		return
	}

	jsonBytes, err := json.Marshal(u)
	if err != nil {
		internalServerError(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		internalServerError(w, r)
		return
	}

	h.store.Lock()
	h.store.m[user.ID] = user
	h.store.Unlock()
	jsonBytes, err := json.Marshal(user)

	if err != nil {
		internalServerError(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func internalServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("INTERNAL SERVER ERROR"))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("NOT FOUND"))
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	switch {
	case r.Method == http.MethodGet && listUserRe.MatchString(r.URL.Path):
		h.List(w, r)
		return
	case r.Method == http.MethodGet && getUserRe.MatchString(r.URL.Path):
		h.Get(w, r)
		return
	case r.Method == http.MethodPost && createUserRe.MatchString(r.URL.Path):
		h.Create(w, r)
		return
	default:
		notFound(w, r)
		return
	}
}

func main() {
	mux := http.NewServeMux()
	userH := &UserHandler{
		store: &DataStore{
			m: map[string]User{
				"1": {ID: "1", Name: "Lolen"},
			},
			RWMutex: &sync.RWMutex{},
		},
	}
	// mux.Handle("/users", &UserHandler{})
	mux.Handle("/users", userH)
	mux.Handle("/users/", userH)
	http.ListenAndServe(":8080", mux)
}
