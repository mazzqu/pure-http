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

type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type datastore struct {
	m map[string]user
	*sync.RWMutex
}

type userHandler struct {
	store *datastore
}

func (h *userHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
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
	case r.Method == http.MethodDelete && getUserRe.MatchString(r.URL.Path):
		h.Delete(w, r)
	default:
		notFound(w, r)
		return
	}
}

func (h *userHandler) List(w http.ResponseWriter, r *http.Request) {
	h.store.RLock()
	users := make([]user, 0, len(h.store.m))
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

func (h *userHandler) Get(w http.ResponseWriter, r *http.Request) {
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
		w.Write([]byte("user not found"))
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

func (h *userHandler) Create(w http.ResponseWriter, r *http.Request) {
	var u user
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		internalServerError(w, r)
		return
	}
	h.store.Lock()
	h.store.m[u.ID] = u
	h.store.Unlock()
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		internalServerError(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// -- got wrong approach
/*
	In Go, when you access a map, you receive two values: the value and a boolean
	indicating whether the key exists in the map. In your case, you're trying
	to assign this boolean (err) to a variable that you've named err,
	which is misleading. This can be fixed by renaming it
	to a more appropriate name (e.g., exists), and handling the case correctly.

*/
// func (h *userHandler) Delete(w http.ResponseWriter, r *http.Request) {
// 	matches := getUserRe.FindStringSubmatch(r.URL.Path)

// 	_, err := h.store.m[matches[1]]
// 	if err {
// 		w.WriteHeader(http.StatusNotFound)
// 		w.Write([]byte("user not found"))
// 		return
// 	}

// 	h.store.Lock()
// 	delete(h.store.m, matches[1])
// 	h.store.Unlock()

// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("user deleted"))
// }

// new approach DELETE method
/*
	Explanation of the changes:
1.Check for existence in the map:

--The expression _, err := h.store.m[matches[1]] is incorrect because err is not
the right name for the second value returned from map access.
--I renamed the variable to exists, which makes more sense, and the check now correctly verifies if the user exists in the map using if _, exists := h.store.m[userID]; !exists.

2.Using the user ID correctly:

--userID := matches[1] captures the user ID, so we can use it in both
the existence check and for deletion.

3.Deleting the user:

--The user is only deleted after confirming that it exists,
and the map is locked during the operation to prevent concurrent modifications.

This should resolve the issue where the DELETE method doesn't work as expected.

*/
func (h *userHandler) Delete(w http.ResponseWriter, r *http.Request) {
	matches := getUserRe.FindStringSubmatch(r.URL.Path)

	// Check if the user exists in the store
	userID := matches[1]
	if _, exists := h.store.m[userID]; !exists {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		return
	}

	// Lock the store, delete the user, and unlock
	h.store.Lock()
	delete(h.store.m, userID)
	h.store.Unlock()

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("user deleted"))
}

func internalServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("internal server error"))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("not found"))
}

func main() {
	mux := http.NewServeMux()
	userH := &userHandler{
		store: &datastore{
			m: map[string]user{
				// "1": user{ID: "1", Name: "bob"}, this typing redundant reason
				"1": {ID: "1", Name: "bob"},
			},
			RWMutex: &sync.RWMutex{},
		},
	}
	// mux.Handle("/users", &UserHandler{})
	mux.Handle("/users", userH)
	mux.Handle("/users/", userH)

	http.ListenAndServe("localhost:8080", mux)
}
