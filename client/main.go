package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Challenge struct {
	Hash       string `json:"hash"`
	Timestamp  int64  `json:"timestamp"`
	Difficulty int    `json:"difficulty"`
	Time       int    `json:"time"`
}

type QuoteRequest struct {
	Hash  string `json:"hash"`
	Nonce int64  `json:"nonce"`
}

func main() {
	host := flag.String("host", "http://localhost", "Server host")
	flag.Parse()

	resp, err := http.Get(*host + "/challenge")
	if err != nil {
		fmt.Println(err)
		return
	}

	var challenge Challenge
	err = json.NewDecoder(resp.Body).Decode(&challenge)
	if err != nil {
		fmt.Println(err)
		return
	}

	nonce, ok := SolveChallenge(challenge)
	if !ok {
		fmt.Println("Timeout")
		return
	}

	request := QuoteRequest{challenge.Hash, nonce}

	marshalled, err := json.Marshal(request)
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest("GET", *host+"/quote", bytes.NewReader(marshalled))

	if err != nil {
		fmt.Println(err)
		return
	}

	client := http.Client{Timeout: 1 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(resBody))
}
