package cli

import (
	"fmt"
	"a1/client"
)

func Run() (int, error) {
	fmt.Println("Parsing command line args")
	udpAddr, err := ParseUDPAddr("192.168.0.1:3000")
	if err != nil {
		return 1, fmt.Errorf("cli failed: %v", err)
	}
	tcpAddr, err := ParseTCPAddr("192.168.0.1:9090")
	if err != nil {
		return 1, fmt.Errorf("cli failed: %v", err)
	}
	aServerAddr, err := ParseUDPAddr("127.0.0.1:3030")
	if err != nil {
		return 1, fmt.Errorf("cli failed: %v", err)
	}
	return client.Execute(udpAddr, tcpAddr, aServerAddr)
}
