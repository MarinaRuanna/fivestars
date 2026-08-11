package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type ClaimCodeHasher interface {
	Hash(raw string) string
}

type HMACClaimCodeHasher struct {
	secret string
}

func NewHMACClaimCodeHasher(secret string) HMACClaimCodeHasher {
	return HMACClaimCodeHasher{secret: secret}
}

func (h HMACClaimCodeHasher) Hash(raw string) string {
	mac := hmac.New(sha256.New, []byte(h.secret))
	_, _ = mac.Write([]byte(raw))
	return hex.EncodeToString(mac.Sum(nil))
}
