package internel

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"net/http"
	"time"
)

type Handler struct {
	maxLength int
	quota     int
	whitelist map[string]struct{}
	nums      map[string]int

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

func NewHandler(conf *Config, db *sql.DB, mail *mailSender) *Handler {

	now := time.Now()

	dbOperator, err := newDbOperator(db, now.Format("2006-01-02"))
	if err != nil {
		return nil
	}

	quota := make(map[string]int)

	for i := conf.Quota; i >= 0; i-- {
		date := now.Format("2006-01-02")
		quota[date], _ = dbOperator.count(date)
		now = now.Add(time.Hour * 24)
	}

	whiteList := make(map[string]struct{})

	for i := 0; i < len(conf.WhiteList); i++ {
		whiteList[conf.WhiteList[i]] = struct{}{}
	}

	return &Handler{
		maxLength: conf.MaxLength,
		quota:     conf.Quota,
		nums:      quota,
		whitelist: whiteList,
		db:        dbOperator,
		mail:      mail,
	}
}

func (h *Handler) checkOrigin(w http.ResponseWriter, r *http.Request) bool {
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

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	ok := h.checkOrigin(w, r)
	if r.Method != "POST" || !ok {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.ContentLength > int64(h.maxLength)*1024 {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}

	var tmp request

	err := json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if h.nums[tmp.Date] >= h.maxLength {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	err = h.db.insert(&tmp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.nums[tmp.Date]++

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

func (h *Handler) Send(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) Check(w http.ResponseWriter, r *http.Request) {
	ok := h.checkOrigin(w, r)
	if r.Method != "POST" || !ok {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var tmp struct {
		Date string
	}

	err := json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]bool)
	resp["status"] = true

	if h.nums[tmp.Date] >= h.maxLength {
		resp["status"] = false
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}
