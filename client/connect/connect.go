// This file wraps around an http interface dependancy.
package connect

import (
	"fmt"
	"net"
)

func UDP(localAddr net.UDPAddr, remoteAddr net.UDPAddr, msg []byte) ([]byte, error) {
	conn, err := net.DialUDP("udp", &localAddr, &remoteAddr)
	if err != nil {
		return nil, fmt.Errorf("opening a connection to %s from %s as UDPs failed: %v", remoteAddr.String(), localAddr.String(), err)
	}
	defer conn.Close()
	conn.Write(msg)

	buffer := make([]byte, 1024)
	conn.Read(buffer)
	return buffer, nil
}

func TCP(localAddr net.TCPAddr, remoteAddr net.TCPAddr, msg []byte) ([]byte, error) {
	conn, err := net.DialTCP("tcp", &localAddr, &remoteAddr)
	if err != nil {
		return nil, fmt.Errorf("opening a connection to %s from %s as TCP failed: %v", remoteAddr.String(), localAddr.String(), err)
	}
	defer conn.Close()
	conn.Write(msg)

	buffer := make([]byte, 1024)
	conn.Read(buffer)

	return buffer, nil
}
