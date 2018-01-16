package compute

import (
	"testing"
)

func TestSecret(t *testing.T) {
	var tests []struct{
		nonce string
		numZeros int64
	}{
		{
			"test",
			3,
		},
	}
	for _, test := range tests {
		result, err := Secret(test.nonce, test.numZeros)
		if 
