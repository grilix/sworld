package sworld

import (
	"math/rand"
	"time"
)

var randomIDSrc = rand.NewSource(time.Now().UnixNano())

const randomIDLetterBytes = "0123456789" +
	"abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandomID(size int) string {
	b := make([]byte, size)

	for i, cache, remain := size-1, randomIDSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randomIDSrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(randomIDLetterBytes) {
			b[i] = randomIDLetterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
