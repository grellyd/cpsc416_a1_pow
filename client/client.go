package client

import (
	"as1_c6y8/addrparse"
	"as1_c6y8/client/compute"
	"as1_c6y8/client/connect"
	"encoding/json"
	"fmt"
	"net"
    "bytes"
)

var SUCCESS = 0
var FAILURE = 1
var nullByte = "\x00"

func Execute(localUDPAddr net.UDPAddr, localTCPAddr net.TCPAddr, aServerAddr net.UDPAddr) (int, error) {
	fmt.Printf("udpAddr: %s,\ntcpAddr: %s,\naServerAddr: %s\n", localUDPAddr.String(), localTCPAddr.String(), aServerAddr.String())
	fmt.Println("This command has now run")

	arbitraryStr := "arbitrary"
	// fetch the nonce
	nonceMsg, err := getNonce(localUDPAddr, aServerAddr, arbitraryStr)
	if err != nil {
		return FAILURE, fmt.Errorf("fetching nonce from %s with %s failed: %v", aServerAddr.String(), arbitraryStr, err)
	}
	// compute the secret
	secret, err := compute.Secret(nonceMsg.Nonce, nonceMsg.N)
	if err != nil {
		return FAILURE, fmt.Errorf("computing secret for %s with %d zeros failed: %v", nonceMsg.Nonce, nonceMsg.N, err)
	}
	// fetch the fserver info
	fortuneInfo, err := sendSecret(localUDPAddr, aServerAddr, secret)
	if err != nil {
		return FAILURE, fmt.Errorf("sending secret to %s with %s failed: %v", aServerAddr.String(), secret, err)
	}
	fServerAddr, err := addrparse.TCP(fortuneInfo.FortuneServer)
	if err != nil {
		return FAILURE, fmt.Errorf("parsing fServerAddr at %s failed: %v", fServerAddr.String(), err)
	}
	fortuneRequest := FortuneReqMessage{FortuneNonce: fortuneInfo.FortuneNonce}
	fortune, err := requestFortune(localTCPAddr, fServerAddr, fortuneRequest)
	if err != nil {
		return FAILURE, fmt.Errorf("requesting fortune at %v with %v failed: %v", fServerAddr, fortuneRequest, err)
	}
	fmt.Println(fortune.Fortune)
	return SUCCESS, nil
}

func getNonce(localUDPAddr net.UDPAddr, aServerAddr net.UDPAddr, msg string) (nonceMsg NonceMessage, err error) {
	byteMessage, err := json.Marshal(msg)
	if err != nil {
		return NonceMessage{}, fmt.Errorf("opening a connection to %s for nonce failed: %v", aServerAddr.String(), err)
	}
	response, err := connect.UDP(localUDPAddr, aServerAddr, byteMessage)
	if err != nil {
		return NonceMessage{}, fmt.Errorf("opening a connection to %s for nonce failed: %v", aServerAddr.String(), err)
	}
	unMarshalErr := json.Unmarshal(bytes.Trim(response, nullByte), &nonceMsg)
	if unMarshalErr != nil {
		errMsg := ErrMessage{}
		err := json.Unmarshal(bytes.Trim(response, nullByte), &errMsg)
		if err != nil {
			return NonceMessage{}, fmt.Errorf("unable to marshal error response %v to client.ErrMessage: %v", response, err)
		} else {
			return NonceMessage{}, fmt.Errorf("server sent back error %s: %v", errMsg.Error, err)
		}
	}
	return nonceMsg, nil
}

// TODO: Extract UDP Comm out with getNonce. Pass return type.
// TODO: Extract out Unmarshall
func sendSecret(localUDPAddr net.UDPAddr, aServerAddr net.UDPAddr, secret string) (fortuneInfo FortuneInfoMessage, err error) {
	byteMessage, err := json.Marshal(secret)
	if err != nil {
		return FortuneInfoMessage{}, fmt.Errorf("opening a connection to %s for nonce failed: %v", aServerAddr.String(), err)
	}
	response, err := connect.UDP(localUDPAddr, aServerAddr, byteMessage)
	if err != nil {
		return FortuneInfoMessage{}, fmt.Errorf("opening a connection to %s for secret failed: %v", aServerAddr.String(), err)
	}
    fmt.Println(string(response))
	unMarshalErr := json.Unmarshal(bytes.Trim(response, nullByte), &fortuneInfo)
    fmt.Println(fortuneInfo.FortuneServer)
	if unMarshalErr != nil {
		errMsg := ErrMessage{}
		err := json.Unmarshal(bytes.Trim(response, nullByte), &errMsg)
		if err != nil {
			return FortuneInfoMessage{}, fmt.Errorf("unable to marshal error response %v to client.ErrMessage: %v", response, err)
		} else {
			return FortuneInfoMessage{}, fmt.Errorf("server sent back error %s: %v", errMsg.Error, err)
		}
	}
	return fortuneInfo, nil
}

func requestFortune(localTCPAddr net.TCPAddr, fServerAddr net.TCPAddr, fortuneRequest FortuneReqMessage) (fortune FortuneMessage, err error) {
	msg, err := json.Marshal(fortuneRequest)
	response, err := connect.TCP(localTCPAddr, fServerAddr, msg)
	unMarshalErr := json.Unmarshal(bytes.Trim(response, nullByte), &fortune)
	if unMarshalErr != nil {
		errMsg := ErrMessage{}
		err := json.Unmarshal(bytes.Trim(response, nullByte), &errMsg)
		if err != nil {
			return FortuneMessage{}, fmt.Errorf("unable to marshal error response %v to client.ErrMessage: %v", response, err)
		} else {
			return FortuneMessage{}, fmt.Errorf("server sent back error %s: %v", errMsg.Error, err)
		}
	}
	return fortune, nil
}
