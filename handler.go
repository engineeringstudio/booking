package main

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"net/http"
	"time"
)

type handler struct {
	name      string
	whitelist map[string]struct{}

	db   *dbOperator
	mail *mailSender
}

type request struct {
	Name string `json:"name"`
	Sno  string `json:"id"`
	Pn   string `json:"phone"`
	Date string `json:"date"`
	Info string `json:"issue"`
}

func NewHandler(name string, db *sql.DB, mail *mailSender) *handler {
	dbOperator, err := newDbOperator(db, name)
	if err != nil {
		return nil
	}

	return &handler{
		name:      name,
		whitelist: make(map[string]struct{}),
		db:        dbOperator,
		mail:      mail,
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

	if r.ContentLength > int64(conf.MaxLength)*1024 {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}

	var tmp request

	err := json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = h.db.insert(&tmp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	hasher := fnv.New64()
	key := hex.EncodeToString(hasher.Sum([]byte(tmp.Name)))

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

func (h *handler) send(w http.ResponseWriter, r *http.Request) {
	ok := h.checkOrigin(w, r)
	if r.Method != "GET" || !ok {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	date := time.Now()
	tmp, err := h.db.query(date.Format("2006-01-02"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.mail.send(fmt.Sprint(tmp))
	w.WriteHeader(http.StatusOK)
}
