package compute

import (
	"crypto/md5"
	"fmt"
	"encoding/hex"
	"strings"
	"math/rand"
	"time"
)

var characters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func Secret(nonce string, numZeros int64) (secret string, err error) {
	secret = ""
	rand.Seed(time.Now().UnixNano())
	for i := 1; i < 1024; i++ {
		for j := 0; j < 1024; j++ {
			secret = generateRandomString(i)
			fmt.Printf("Trying: %s\n", secret)
			numPresentZeros := int64(strings.Count(ComputeNonceSecretHash(nonce, secret), "0")) 
			if numPresentZeros == numZeros {
				return secret, nil
			}
		}
	}
	return "", fmt.Errorf("unable to find secret")
}

// Returns the MD5 hash as a hex string for the (nonce + secret) value.
func ComputeNonceSecretHash(nonce string, secret string) string {
	h := md5.New()
	h.Write([]byte(nonce + secret))
	str := hex.EncodeToString(h.Sum(nil))
	return str
}

func generateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, length)
	for i := range b {
        b[i] = characters[rand.Intn(len(characters))]
    }
	return string(b)
}


