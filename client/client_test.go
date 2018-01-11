package client_test

import (
	"a1/client"
	"net"
	"testing"
)

func TestClient(t *testing.T) {
	var tests = []struct {
		output int
	}{
		{client.SUCCESS},
	}
	for _, test := range tests {
		result, err := client.Execute(net.UDPAddr{}, net.TCPAddr{}, net.UDPAddr{})
		if err != nil {
			t.Errorf("Bad Exit: \"client.Execute()\" produced err: %v", err)
		}
		if result != test.output {
			t.Errorf("Bad Exit: \"client.Execute()\" = %q not %q", result, test.output)
		}
	}
}
