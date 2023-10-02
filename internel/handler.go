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
	whitelist map[string]struct{}
	quota     map[string]int

	db   *dbOperator
	db_n *dbOperator
	mail *mailSender
}

type request struct {
	Name string `json:"name"`
	Sno  string `json:"id"`
	Pn   string `json:"phone"`
	Date string `json:"date"`
	Info string `json:"issue"`
}

func getDuration() time.Duration {
	now := time.Now()
	tmp := now.Format("20060102")
	today, _ := time.Parse("20060102", tmp)
	ret := today.Add(time.Second * 10).Sub(now)
	return ret
}

func NewHandler(conf *Config, db *sql.DB, mail *mailSender) (*Handler, error) {
	now := time.Now()

	dbOperator, err := newDbOperator(db, now.Format("200601"))
	if err != nil {
		return nil, err
	}

	dbOperator_n, err := newDbOperator(db, now.AddDate(0, 1, 0).Format("200601"))
	if err != nil {
		return nil, err
	}

	quota := make(map[string]int)

	fmt.Println("1")

	for i := conf.Quota; i >= 0; i-- {
		date := now.Format("2006-01-02")

		a, err := dbOperator.count(date)
		if err != nil {
			return nil, err
		}

		b, err := dbOperator_n.count(date)
		if err != nil {
			return nil, err
		}

		if a >= b {
			quota[date] = conf.MaxLength - a
		} else {
			quota[date] = conf.MaxLength - b
		}

		now = now.AddDate(0, 0, 1)
	}

	whiteList := make(map[string]struct{})

	for i := 0; i < len(conf.WhiteList); i++ {
		whiteList[conf.WhiteList[i]] = struct{}{}
	}

	fmt.Println(conf.WhiteList)

	return &Handler{
		maxLength: conf.MaxLength,
		quota:     quota,
		whitelist: whiteList,
		db:        dbOperator,
		db_n:      dbOperator_n,
		mail:      mail,
	}, nil
}

func (h *Handler) UpdateTask() {
	timer := time.NewTimer(0)
	for {
		<-timer.C
		h.db.createTable(time.Now().Format("200601"))
		timer.Reset(getDuration())
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

	if h.quota[tmp.Date] > 0 {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	err = h.db.insert(&tmp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.quota[tmp.Date]--

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

	if h.quota[tmp.Date] > 0 {
		resp["status"] = false
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}
