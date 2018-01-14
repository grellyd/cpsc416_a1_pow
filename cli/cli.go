package cli

import (
	"a1/client"
	"fmt"
)

func Run(args []string) (int, error) {
	udpAddr, err := ParseUDPAddr(args[0])
	if err != nil {
		return 1, fmt.Errorf("cli failed: %v", err)
	}
	tcpAddr, err := ParseTCPAddr(args[1])
	if err != nil {
		return 1, fmt.Errorf("cli failed: %v", err)
	}
	aServerAddr, err := ParseUDPAddr(args[2])
	if err != nil {
		return 1, fmt.Errorf("cli failed: %v", err)
	}
	return client.Execute(udpAddr, tcpAddr, aServerAddr)
}
