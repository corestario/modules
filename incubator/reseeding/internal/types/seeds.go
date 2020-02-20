package types

import (
	"crypto/md5"
	"encoding/hex"
)

type SeedVotes map[string]struct{}

type Seeds map[string]SeedVotes

func getKey(seed []byte) string {
	seedHash := md5.Sum(seed)
	return hex.EncodeToString(seedHash[:])
}

func (s Seeds) Add(seed []byte, sender string) {
	key := getKey(seed)
	votes, ok := s[key]
	if !ok {
		votes = make(map[string]struct{})
	}
	votes[sender] = struct{}{}
	s[key] = votes
}

func (s Seeds) GetVotesForSeed(seed []byte) int {
	return len(s[getKey(seed)])
}
