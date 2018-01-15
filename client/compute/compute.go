package compute

import (
	"crypto/md5"
	"encoding/hex"
)

func ComputeSecret(nonce string, numZeros int64) (secret string, err error) {
	return "", nil
}

// Returns the MD5 hash as a hex string for the (nonce + secret) value.
func computeNonceSecretHash(nonce string, secret string) string {
	h := md5.New()
	h.Write([]byte(nonce + secret))
	str := hex.EncodeToString(h.Sum(nil))
	return str
}
