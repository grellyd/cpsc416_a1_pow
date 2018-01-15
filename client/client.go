package client

import (
	"fmt"
	"net"
	"a1/client/connect"
	"a1/addrparse"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
)


var SUCCESS = 0
var FAILURE = 1

func Execute(localUDPAddr net.UDPAddr, localTCPAddr net.TCPAddr, aServerAddr net.UDPAddr) (int, error) {
	fmt.Printf("udpAddr: %s,\ntcpAddr: %s,\naServerAddr: %s\n", localUDPAddr.String(), localTCPAddr.String(), aServerAddr.String())
	fmt.Println("This command has now run")

	addrString := ""
	fServerAddr, err := addrparse.TCP(addrString)
	if err != nil {
		// TODO
	}
	fortuneRequest := FortuneReqMessage{}
	fortune, err := requestFortune(localTCPAddr, fServerAddr, fortuneRequest)
	if err != nil {
		// TODO
	}
	fmt.Println(fortune.Fortune)
	return SUCCESS, nil
}


func getNonce(localUDPAddr net.UDPAddr, aServerAddr net.UDPAddr) (nonce NonceMessage, err error) {
	byteMessage, err := json.Marshal("arbitrary")
	if err != nil {
		return NonceMessage{}, fmt.Errorf("opening a connection to %s for nonce failed: %v", aServerAddr.String(), err)
	}
	response, err := connect.UDP(localUDPAddr, aServerAddr, byteMessage)
	if err != nil {
		return NonceMessage{}, fmt.Errorf("opening a connection to %s for nonce failed: %v", aServerAddr.String(), err)
	}
	nonceMsg := NonceMessage{}
	unMarshalErr := json.Unmarshal(response, &nonceMsg)
	if unMarshalErr != nil {
		errMsg := ErrMessage{}
		err := json.Unmarshal(response, &errMsg)
		if err != nil {
			return NonceMessage{}, fmt.Errorf("unable to marshal error response %v to client.ErrMessage: %v", response,  err)
		} else {
			return NonceMessage{}, fmt.Errorf("server sent back error %s: %v", errMsg.Error,  err)
		}
	}
	return nonceMsg, nil
}

// TODO: Extract UDP Comm out with getNonce. Pass return type.
func sendSecret(localUDPAddr net.UDPAddr, aServerAddr net.UDPAddr, secret string) (fortuneInfo FortuneInfoMessage, err error) {
	byteMessage, err := json.Marshal(secret)
	if err != nil {
		return FortuneInfoMessage{}, fmt.Errorf("opening a connection to %s for nonce failed: %v", aServerAddr.String(), err)
	}
	response, err := connect.UDP(localUDPAddr, aServerAddr, byteMessage)
	if err != nil {
		return FortuneInfoMessage{}, fmt.Errorf("opening a connection to %s for secret failed: %v", aServerAddr.String(), err)
	}
	unMarshalErr := json.Unmarshal(response, &fortuneInfo)
	if unMarshalErr != nil {
		errMsg := ErrMessage{}
		err := json.Unmarshal(response, &errMsg)
		if err != nil {
			return FortuneInfoMessage{}, fmt.Errorf("unable to marshal error response %v to client.ErrMessage: %v", response,  err)
		} else {
			return FortuneInfoMessage{}, fmt.Errorf("server sent back error %s: %v", errMsg.Error,  err)
		}
	}
	return fortuneInfo, nil
}

func requestFortune(localTCPAddr net.TCPAddr, fServerAddr net.TCPAddr, fortuneRequest FortuneReqMessage) (fortune FortuneMessage, err error) {
	msg, err := json.Marshal(fortuneRequest)
	response, err := connect.TCP(localTCPAddr, fServerAddr, msg)
	unMarshalErr := json.Unmarshal(response, &fortune)
	if unMarshalErr != nil {
		errMsg := ErrMessage{}
		err := json.Unmarshal(response, &errMsg)
		if err != nil {
			return FortuneMessage{}, fmt.Errorf("unable to marshal error response %v to client.ErrMessage: %v", response,  err)
		} else {
			return FortuneMessage{}, fmt.Errorf("server sent back error %s: %v", errMsg.Error,  err)
		}
	}
	return fortune, nil
}

// Returns the MD5 hash as a hex string for the (nonce + secret) value.
func computeNonceSecretHash(nonce string, secret string) string {
	h := md5.New()
	h.Write([]byte(nonce + secret))
	str := hex.EncodeToString(h.Sum(nil))
	return str
}
