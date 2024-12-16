package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"sync"
)

// Regex files for path's URL
var (
	listUsersRe    = regexp.MustCompile(`^\/users[\/]*$`)
	spesificUserRe = regexp.MustCompile(`^\/users\/(\d+)$`)
)

type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Since it storing data at memory, will
// gone after terminated the program.
type datastore struct {
	m map[string]user
	*sync.RWMutex
}

type userHandler struct {
	store *datastore
}

func notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("User not found"))
}

func internalServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Server down, please restart !!!"))
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(`{"error":"}` + message))
}

func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not encode response")
	}
}

func (h *userHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodGet && listUsersRe.MatchString(r.URL.Path):
		h.List(w, r)
		return
	case r.Method == http.MethodGet && spesificUserRe.MatchString(r.URL.Path):
		h.Get(w, r)
		return
	case r.Method == http.MethodPost && listUsersRe.MatchString(r.URL.Path):
		h.Create(w, r)
		return
	case r.Method == http.MethodPut && spesificUserRe.MatchString(r.URL.Path):
		h.Update(w, r)
	case r.Method == http.MethodDelete && spesificUserRe.MatchString(r.URL.Path):
		h.Delete(w, r)
	default:
		notFound(w)
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
	log.Println("Lists are sets.")
}

func (h *userHandler) Get(w http.ResponseWriter, r *http.Request) {
	matches := spesificUserRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		notFound(w)
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
		log.Println("Data not satify the requirements please check !!!")
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

func (h *userHandler) Update(w http.ResponseWriter, r *http.Request) {
	matches := spesificUserRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		notFound(w)
		return
	}

	userID := matches[1]
	h.store.Lock()
	defer h.store.Unlock()

	// Check if the user exists
	u, exists := h.store.m[userID]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		return
	}

	// Decode the updated user details
	var updatedUser user
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		//internalServerError(w)
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Update the user fields, here we are just replacing the entire user object
	u.Name = updatedUser.Name
	h.store.m[userID] = u

	// Respond with the updated user
	if err := json.NewEncoder(w).Encode(u); err != nil {
		//internalServerError(w)
		respondWithJSON(w, http.StatusOK, u)
	}
}

func (h *userHandler) Delete(w http.ResponseWriter, r *http.Request) {
	matches := spesificUserRe.FindStringSubmatch(r.URL.Path)

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
