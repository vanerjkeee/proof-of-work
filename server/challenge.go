package main

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const salt = "MoQEnUWjtSzX"

type Challenge struct {
	Hash       string
	Timestamp  int64
	Difficulty int
	Time       int
}

// Challenges can be stored in memcache/redis instead of in-memory
type ChallengeManager struct {
	storage    map[string]Challenge
	difficulty int
	time       int
}

func NewChallengeManager(difficulty, time int) ChallengeManager {
	cm := ChallengeManager{make(map[string]Challenge), difficulty, time}
	go cm.cleaning()

	return cm
}

/*
+ Checking a challenge in one request
+ Checking on the server does not require a lot of resources
+ Simple complication or simplification of a challenge depending on the server load
*/
func (cm *ChallengeManager) NewChallenge() Challenge {
	h := sha256.New()
	h.Write([]byte(uuid.NewString() + salt))
	hash := fmt.Sprintf("%x", h.Sum(nil))
	now := time.Now()
	timestamp := now.Unix()

	challenge := Challenge{hash, timestamp, cm.difficulty, cm.time}
	cm.storage[challenge.Hash] = challenge

	return challenge
}

func (cm *ChallengeManager) CheckChallenge(hash string, nonce int64) bool {
	ch, ok := cm.storage[hash]
	if !ok {
		return false
	}

	now := time.Now()
	if (now.Unix() - ch.Timestamp) > int64(ch.Time) {
		return false
	}

	h := sha256.New()
	h.Write([]byte(hash + strconv.FormatInt(ch.Timestamp, 10) + strconv.FormatInt(nonce, 10)))
	hashResult := fmt.Sprintf("%x", h.Sum(nil))

	need := strings.Repeat("0", ch.Difficulty)
	if need != hashResult[:ch.Difficulty] {
		return false
	}

	delete(cm.storage, hash)

	return true
}

func (cm *ChallengeManager) cleaning() Challenge {
	for {
		time.Sleep(time.Duration(cm.time) * time.Second)
		now := time.Now()
		for key, ch := range cm.storage {
			if (now.Unix() - ch.Timestamp) > int64(cm.time) {
				delete(cm.storage, key)
			}
		}
	}
}
