package hash

import (
	"crypto/sha256"
	"crypto/subtle"
)

type TokenHasher struct {
}

func (t *TokenHasher) Hash(data []byte) ([]byte, error) {
	hash := sha256.Sum256(data)
	return hash[:], nil
}

func (t *TokenHasher) Verify(data []byte, hash []byte) bool {
	h, err := t.Hash(data)
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare(h, hash) == 1
}

func NewTokenHasher() Hasher {
	return &TokenHasher{}
}
