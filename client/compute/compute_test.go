package compute_test

import (
	"testing"
	"strings"
	"as1_c6y8/client/compute"
)

func TestSecret(t *testing.T) {
	var tests = []struct{
		nonce string
		numZeros int64
	}{
		{
			"test",
			3,
		},
		{
			"thisisastring",
			5,
		},
   		{
   			"nonce-ahoy",
   			7,
   		},
// 		{
// 			"nonce-hola",
// 			10,
// 		},
	}
	for _, test := range tests {
		secret, err := compute.Secret(test.nonce, test.numZeros)
		if err != nil {
			t.Errorf("Bad Exit: \"TestSecret(%v)\" produced err: %v", test, err)
		}
		numPresentZeros := int64(strings.Count(compute.ComputeNonceSecretHash(test.nonce, secret), "0")) 
		if !compute.ValidHash(test.nonce, secret, test.numZeros) {
			t.Errorf("Bad Exit: Not enough zeros with %s! %d instead of %d", secret, numPresentZeros, test.numZeros)
		}
	}
}
