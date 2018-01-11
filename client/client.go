package client

import (
	"fmt"
	"net"
)

var SUCCESS = 0
var FAILURE = 1

func Execute(udpAddr net.UDPAddr, tcpAddr net.TCPAddr, aServerAddr net.UDPAddr) (int, error) {
	fmt.Printf("udpAddr: %s,\ntcpAddr: %s,\naServerAddr: %s\n", udpAddr.String(), tcpAddr.String(), aServerAddr.String())
	fmt.Println("This command has now run")
	return SUCCESS, nil
}
