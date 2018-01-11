package client_test

import (
	"a1/client"
	"testing"
)

func TestClient(t *testing.T) {
	var tests = []struct {
		output int
	}{
		{client.SUCCESS},
	}
	for _, test := range tests {
		result := client.Execute()
		if result != test.output {
			t.Errorf("Bad Exit: \"client.Execute()\" = %q not %q", result, test.output)
		}
	}
}
