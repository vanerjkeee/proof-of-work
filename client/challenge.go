package main

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func SolveChallenge(ch Challenge) (int64, bool) {
	var nonce int64
	nonce = 0
	start := time.Now()

	for {
		h := sha256.New()
		h.Write([]byte(ch.Hash + strconv.FormatInt(ch.Timestamp, 10) + strconv.FormatInt(nonce, 10)))
		hash := fmt.Sprintf("%x", h.Sum(nil))
		need := strings.Repeat("0", ch.Difficulty)
		if need == hash[:ch.Difficulty] {
			break
		}

		nonce++

		if time.Since(start).Seconds() > float64(ch.Time) {
			return 0, false
		}

	}

	return nonce, true
}
