package cli

import (
	"net"
	"testing"
)

func TestParseTCPAddr(t *testing.T) {
	var tests = []struct {
		input string
		output net.TCPAddr
	}{
		{
			"192.168.0.1:5000", 
			net.TCPAddr{
				IP: net.IPv4(byte(192), byte(168), byte(0), byte(1)),
				Port: 5000,
				Zone: "",
			},
		},
		{
			"127.134.0.1:3030", 
			net.TCPAddr{
				IP: net.IPv4(byte(127), byte(134), byte(0), byte(1)),
				Port: 3030,
				Zone: "",
			},
		},
		{
			"198.162.33.54:5555", 
			net.TCPAddr{
				IP: net.IPv4(byte(198), byte(162), byte(33), byte(54)),
				Port: 5555,
				Zone: "",
			},
		},

	}
	for _, test := range tests {
		result, err := ParseTCPAddr(test.input)
		if err != nil {
			t.Errorf("Bad Exit: \"ParseTCPAddr(%s)\" produced err: %v", test.input, err)
		}
		if !result.IP.Equal(test.output.IP) || result.Port != test.output.Port || result.Zone != test.output.Zone {
			t.Errorf("Bad Exit: \"ParseTCPAddr(%s)\" = %q not %q", test.input, result, test.output)
		}
	}
}

func TestParseUDPAddr(t *testing.T) {
	var tests = []struct {
		input string
		output net.UDPAddr
	}{
		{
			"192.168.0.1:5000", 
			net.UDPAddr{
				IP: net.IPv4(byte(192), byte(168), byte(0), byte(1)),
				Port: 5000,
				Zone: "",
			},
		},
		{
			"127.134.0.1:3030", 
			net.UDPAddr{
				IP: net.IPv4(byte(127), byte(134), byte(0), byte(1)),
				Port: 3030,
				Zone: "",
			},
		},
		{
			"198.162.33.54:5555", 
			net.UDPAddr{
				IP: net.IPv4(byte(198), byte(162), byte(33), byte(54)),
				Port: 5555,
				Zone: "",
			},
		},

	}
	for _, test := range tests {
		result, err := ParseUDPAddr(test.input)
		if err != nil {
			t.Errorf("Bad Exit: \"ParseUDPAddr(%s)\" produced err: %v", test.input, err)
		}
		if !result.IP.Equal(test.output.IP) || result.Port != test.output.Port || result.Zone != test.output.Zone {
			t.Errorf("Bad Exit: \"ParseUDPAddr(%s)\" = %q not %q", test.input, result, test.output)
		}
	}
}
