// This file wraps around an http interface dependancy.
package connect

import (
	"fmt"
	"net"
)

func UDP(addr net.UDPAddr, msg string) ([]byte, error) {
	return communicate("udp", addr.String(), msg)
}

func TCP(addr net.TCPAddr, msg string) ([]byte, error) {
	return communicate("tcp", addr.String(), msg)
}

func communicate(network string, addr string, msg string) ([]byte, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, fmt.Errorf("opening a connection to %s as %s failed: %v", addr, network, err)
	}
	defer conn.Close()
	conn.Write([]byte(msg))

	buffer := make([]byte, 1024)
	conn.Read(buffer)

	return buffer, nil
}
