package hmachash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
)

type HMAC interface {
	Hash(input string) string
}

type hm struct {
	hmac hash.Hash
}

func NewHMAC(key string) hm {
	h := hmac.New(sha256.New, []byte(key))
	return hm{
		hmac: h,
	}
}

func (h hm) Hash(input string) string {
	h.hmac.Reset()
	h.hmac.Write([]byte(input))
	hashedData := h.hmac.Sum(nil)
	return base64.URLEncoding.EncodeToString(hashedData)
}
