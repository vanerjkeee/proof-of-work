package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
)

type Controller struct {
	challengeManager ChallengeManager
}

func main() {
	difficulty := flag.Int("difficulty", 4, "challenge difficulty")
	time := flag.Int("time", 60, "time to solve challenge")
	flag.Parse()

	c := Controller{NewChallengeManager(*difficulty, *time)}
	http.HandleFunc("/challenge", c.GetChallenge)
	http.HandleFunc("/quote", c.Quote)

	http.ListenAndServe(fmt.Sprintf(":%d", 80), nil)
}

func (c *Controller) GetChallenge(w http.ResponseWriter, r *http.Request) {
	ch := c.challengeManager.NewChallenge()
	w.Header().Add("Content-Type", "application/json")
	resp := GetChallengeResponse{ch.Hash, ch.Timestamp, ch.Difficulty, ch.Time}
	json.NewEncoder(w).Encode(resp)
}

func (c *Controller) Quote(w http.ResponseWriter, r *http.Request) {
	var request QuoteRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	checked := c.challengeManager.CheckChallenge(request.Hash, request.Nonce)
	if !checked {
		http.Error(w, "Wrong answer", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	resp := QuoteResponse{GetQuote()}
	json.NewEncoder(w).Encode(resp)
}

type GetChallengeResponse struct {
	Hash       string `json:"hash"`
	Timestamp  int64  `json:"timestamp"`
	Difficulty int    `json:"difficulty"`
	Time       int    `json:"time"`
}

type QuoteRequest struct {
	Hash  string `json:"hash"`
	Nonce int64  `json:"nonce"`
}

type QuoteResponse struct {
	Quote string `json:"quote"`
}
