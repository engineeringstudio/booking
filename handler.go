package main

import (
	"database/sql"
	"net/http"
)

type handler struct {
	name      string
	whitelist map[string]struct{}
}

func NewHandler(name string, db *sql.DB) *handler {
	return &handler{
		name:      name,
		whitelist: make(map[string]struct{}),
	}
}

func (h *handler) checkOrigin(w http.ResponseWriter, r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}

	_, ok := h.whitelist[origin]
	if !ok {
		return false
	}
	w.Header().Add("Access-Control-Allow-Origin", origin)
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Content-Type", "text/plain; charset=UTF-8")
	return true
}

func (h *handler) add(w http.ResponseWriter, r *http.Request) {
	ok := h.checkOrigin(w, r)
	if r.Method != "POST" || !ok {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.ContentLength > int64(conf.MaxLength)*1024+64 {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}

	key := ""

	http.SetCookie(w, &http.Cookie{
		Name:     "token_" + key,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		Value:    key,
		MaxAge:   3600,
	})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(key))
}
