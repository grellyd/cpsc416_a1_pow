package compute

import (
	"crypto/md5"
	"fmt"
	"encoding/hex"
	"strings"
	"math/rand"
	"time"
)

var characters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func Secret(nonce string, numZeros int64) (secret string, err error) {
	secret = ""
	rand.Seed(time.Now().UnixNano())
	for {
		for i := 3; i < 10; i++ {
			secret = generateRandomString(i)
//			fmt.Printf("Trying: %s\n", secret)
			if ValidHash(nonce, secret, numZeros) {
				fmt.Println(ComputeNonceSecretHash(nonce, secret))
				return secret, nil
			}
		}
	}
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

func ValidHash(nonce string, secret string, numZeros int64) bool {
	valid := false
	hash := ComputeNonceSecretHash(nonce, secret)
	last_index := strings.LastIndex(hash, zeroString(numZeros))
	if int64(last_index) == int64(len(hash)) - numZeros {
		valid = true
	}
	return valid
}

func zeroString(num int64) string {
	str := ""
	for i := 0; i < int(num); i++{
		str += "0"
	}
	return str
}
