package client

import (
	"fmt"
	"net"
	"a1/client/connect"
	"a1/addrparse"
	// "crypto/md5"
	// "encoding/hex"
	// "encoding/json"
)


var SUCCESS = 0
var FAILURE = 1

func Execute(localUDPAddr net.UDPAddr, localTCPAddr net.TCPAddr, aServerAddr net.UDPAddr) (int, error) {
	fmt.Printf("udpAddr: %s,\ntcpAddr: %s,\naServerAddr: %s\n", udpAddr.String(), tcpAddr.String(), aServerAddr.String())
	fmt.Println("This command has now run")
	return SUCCESS, nil
}


func getNonce(aServerAddr net.UDPAddr) (nonce NonceMessage, err ErrMessage) {
	response, err := connect.UDP(aServerAddr, "arbitrary")
	if err != nil {
		return nil, fmt.Errorf("opening a connection to %s for nonce failed: %v", aServerAddr.String(), err)
	}
	nonceMsg := NonceMessage{}
	err := json.Unmarshal(response, &nonceMsg)
	if err != nil {
		errMsg := ErrMessage{}
		err := json.Unmarshal(response, &errMsg)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal response %v to client.ErrMessage: %v", response,  err)
		} else {
			return nil, fmt.Errorf("server sent back error %v: %v", errMsg,  err)
		}
	}
	return nonceMsg, nil
}

// TODO: Extract UDP Comm out with getNonce. Pass return type.
func sendSecret(aServerAddr net.UDPAddr) (fortuneInfo FortuneInfoMessage, err ErrMessage) {
	response, err := connect.UDP(aServerAddr, "arbitrary")
	if err != nil {
		return nil, fmt.Errorf("opening a connection to %s for secret failed: %v", aServerAddr.String(), err)
	}
	fortuneInfo := FortuneInfoMessage{}
	err := json.Unmarshal(response, &fortuneInfo)
	if err != nil {
		errMsg := ErrMessage{}
		err := json.Unmarshal(response, &errMsg)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal response %v to client.ErrMessage: %v", response,  err)
		} else {
			return nil, fmt.Errorf("server sent back error %v: %v", errMsg,  err)
		}
	}
	return fortuneInfo, nil
}

func requestFortune(fServerAddr net.TCPAddr, fortunerRequest FortuneRequestMessage) (fortune FortuneMessage, err ErrMessage) {
	msg := json.Marshal(fortuneRequest)
	response, err := connect.TCP(fServerAddr, msg)
	return nil, nil
}

// Returns the MD5 hash as a hex string for the (nonce + secret) value.
func computeNonceSecretHash(nonce string, secret string) string {
	h := md5.New()
	h.Write([]byte(nonce + secret))
	str := hex.EncodeToString(h.Sum(nil))
	return str
}
