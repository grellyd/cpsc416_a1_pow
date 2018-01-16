package compute

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	"math/rand"
	"time"
)

var characters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func ComputeSecret(nonce string, numZeros int64) (secret string, err error) {
	secret = generateRandomString(5)
	numPresentZeros := int64(strings.Count(computeNonceSecretHash(nonce, secret), "0")) 
	if numPresentZeros >= numZeros {
		return secret, nil
	}
	return "", nil
}

// Returns the MD5 hash as a hex string for the (nonce + secret) value.
func computeNonceSecretHash(nonce string, secret string) string {
	h := md5.New()
	h.Write([]byte(nonce + secret))
	str := hex.EncodeToString(h.Sum(nil))
	return str
}

func generateRandomString(length int) string {
	rand.Seed(time.Now().Unix())
	b := make([]rune, length)
	for i := range b {
        b[i] = characters[rand.Intn(len(characters))]
    }
	return string(b)
}


