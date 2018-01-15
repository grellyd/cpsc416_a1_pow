package cli

import (
	"as1_c6y8/addrparse"
	"as1_c6y8/client"
	"fmt"
)

func Run(args []string) (int, error) {
	udpAddr, err := addrparse.UDP(args[0])
	if err != nil {
		return 1, fmt.Errorf("cli failed: %v", err)
	}
	tcpAddr, err := addrparse.TCP(args[1])
	if err != nil {
		return 1, fmt.Errorf("cli failed: %v", err)
	}
	aServerAddr, err := addrparse.UDP(args[2])
	if err != nil {
		return 1, fmt.Errorf("cli failed: %v", err)
	}
	return client.Execute(udpAddr, tcpAddr, aServerAddr)
}
