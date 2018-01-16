package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

/////////// Msgs used by both auth and fortune servers:

// An error message from the server.
type ErrMessage struct {
	Error string
}

/////////// Auth server msgs:

// Message containing a nonce from auth-server.
type NonceMessage struct {
	Nonce string
	N     int64 // PoW difficulty: number of zeroes expected at end of md5(nonce+secret)
}

// Message containing an the secret value from client to auth-server.
type SecretMessage struct {
	Secret string
}

// Message with details for contacting the fortune-server.
type FortuneInfoMessage struct {
	FortuneServer string // TCP ip:port for contacting the fserver
	FortuneNonce  int64
}

/////////// Fortune server msgs:

// Message requesting a fortune from the fortune-server.
type FortuneReqMessage struct {
	FortuneNonce int64
}

// Response from the fortune-server containing the fortune.
type FortuneMessage struct {
	Fortune string
	Rank    int64 // Rank of this client solution
}

var SUCCESS = 0
var FAILURE = 1
var nullByte = "\x00"
var characters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func main() {
	fmt.Printf("Starting at %s", time.Now())
	args := os.Args[1:]
	udpAddr, err := parseUDPAddr(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	tcpAddr, err := parseTCPAddr(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	aServerAddr, err := parseUDPAddr(args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	code, err := Execute(udpAddr, tcpAddr, aServerAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	os.Exit(code)
}

func Execute(localUDPAddr net.UDPAddr, localTCPAddr net.TCPAddr, aServerAddr net.UDPAddr) (int, error) {
	fmt.Printf("udpAddr: %s,\ntcpAddr: %s,\naServerAddr: %s\n", localUDPAddr.String(), localTCPAddr.String(), aServerAddr.String())

	arbitraryStr := "arbitrary"
	// fetch the nonce
	nonceMsg, err := getNonce(localUDPAddr, aServerAddr, arbitraryStr)
	if err != nil {
		return FAILURE, fmt.Errorf("fetching nonce from %s with %s failed: %v", aServerAddr.String(), arbitraryStr, err)
	}
	// compute the secret
	secret, err := computeSecret(nonceMsg.Nonce, nonceMsg.N)
	if err != nil {
		return FAILURE, fmt.Errorf("computing secret for %s with %d zeros failed: %v", nonceMsg.Nonce, nonceMsg.N, err)
	}
	fmt.Printf("Found secret: \"%s\"", secret)
	// fetch the fserver info
	fortuneInfo, err := sendSecret(localUDPAddr, aServerAddr, secret)
	if err != nil {
		return FAILURE, fmt.Errorf("sending secret to %s with %s failed: %v", aServerAddr.String(), secret, err)
	}
	fServerAddr, err := parseTCPAddr(fortuneInfo.FortuneServer)
	if err != nil {
		return FAILURE, fmt.Errorf("parsing fServerAddr at %s failed: %v", fServerAddr.String(), err)
	}
	fortuneRequest := FortuneReqMessage{FortuneNonce: fortuneInfo.FortuneNonce}
	fortune, err := requestFortune(localTCPAddr, fServerAddr, fortuneRequest)
	if err != nil {
		return FAILURE, fmt.Errorf("requesting fortune at %v with %v failed: %v", fServerAddr, fortuneRequest, err)
	}
	fmt.Println(fortune.Fortune)
	fmt.Println(fortune.Rank)
	return SUCCESS, nil
}

func getNonce(localUDPAddr net.UDPAddr, aServerAddr net.UDPAddr, msg string) (nonceMsg NonceMessage, err error) {
	byteMessage, err := json.Marshal(msg)
	if err != nil {
		return NonceMessage{}, fmt.Errorf("opening a connection to %s for nonce failed: %v", aServerAddr.String(), err)
	}
	response, err := connectUDP(localUDPAddr, aServerAddr, byteMessage)
	if err != nil {
		return NonceMessage{}, fmt.Errorf("opening a connection to %s for nonce failed: %v", aServerAddr.String(), err)
	}
	unMarshalErr := json.Unmarshal(bytes.Trim(response, nullByte), &nonceMsg)
	fmt.Println(string(response))
	if unMarshalErr != nil {
		if len(bytes.Trim(response, nullByte)) == 0 {
			return NonceMessage{}, fmt.Errorf("server sent back empty response: %v", err)
		}
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
	byteMessage, err := json.Marshal(SecretMessage{Secret: secret})
	if err != nil {
		return FortuneInfoMessage{}, fmt.Errorf("opening a connection to %s for nonce failed: %v", aServerAddr.String(), err)
	}
	response, err := connectUDP(localUDPAddr, aServerAddr, byteMessage)
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
	if fortuneInfo.FortuneServer == "" {
		return FortuneInfoMessage{}, fmt.Errorf("Invalid secret %s: %v", secret, err)
	}
	return fortuneInfo, nil
}

func requestFortune(localTCPAddr net.TCPAddr, fServerAddr net.TCPAddr, fortuneRequest FortuneReqMessage) (fortune FortuneMessage, err error) {
	msg, err := json.Marshal(fortuneRequest)
	response, err := connectTCP(localTCPAddr, fServerAddr, msg)
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

// =======================
// Connection Functions
// =======================
func connectUDP(localAddr net.UDPAddr, remoteAddr net.UDPAddr, msg []byte) ([]byte, error) {
	newLocalAddr, err := net.ResolveUDPAddr("udp", localAddr.String())
	if err != nil {
		return nil, fmt.Errorf("resolving addr %s as UDPs failed: %v", localAddr.String(), err)
	}
	conn, err := net.DialUDP("udp", newLocalAddr, &remoteAddr)
	if err != nil {
		return nil, fmt.Errorf("opening a connection to %s from %s as UDPs failed: %v", remoteAddr.String(), localAddr.String(), err)
	}
	defer conn.Close()
	conn.Write(msg)

	buffer := make([]byte, 1024)
	conn.Read(buffer)
	return buffer, nil
}

func connectTCP(localAddr net.TCPAddr, remoteAddr net.TCPAddr, msg []byte) ([]byte, error) {
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

// =======================
// Compute Functions
// =======================
func computeSecret(nonce string, numZeros int64) (secret string, err error) {
	secret = ""
	rand.Seed(time.Now().UnixNano())
	count := 0
	for {
		for i := 3; i < 10; i++ {
			secret = generateRandomString(i)
			count += 1
			//			fmt.Printf("Trying: %s\n", secret)
			if validHash(nonce, secret, numZeros) {
				fmt.Println(computeNonceSecretHash(nonce, secret))
				fmt.Printf("Count: %d", count)
				return secret, nil
			}
		}
	}
}

// Returns the MD5 hash as a hex string for the (nonce + secret) value.
func computeNonceSecretHash(nonce string, secret string) string {
	h := md5.New()
	h.Write([]byte(nonce + secret))
	str := hex.EncodeToString(h.Sum(nil))
	return str
}

func generateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, length)
	for i := range b {
		b[i] = characters[rand.Intn(len(characters))]
	}
	return string(b)
}

func validHash(nonce string, secret string, numZeros int64) bool {
	valid := false
	hash := computeNonceSecretHash(nonce, secret)
	last_index := strings.LastIndex(hash, zeroString(numZeros))
	if int64(last_index) == int64(len(hash))-numZeros {
		valid = true
	}
	return valid
}

func zeroString(num int64) string {
	str := ""
	for i := 0; i < int(num); i++ {
		str += "0"
	}
	return str
}

// =======================
// Arg Parsing Functions
// =======================
func parseTCPAddr(addr string) (net.TCPAddr, error) {
	splitAddr := strings.SplitN(addr, ":", 2)
	ip, err := parseIP(splitAddr[0])
	if err != nil {
		return net.TCPAddr{}, fmt.Errorf("address parsing failed: %s as TCP: %v", addr, err)
	}
	port, err := parsePort(splitAddr[1])
	if err != nil {
		return net.TCPAddr{}, fmt.Errorf("address parsing failed: %s as TCP: %v", addr, err)
	}
	return net.TCPAddr{IP: ip, Port: port, Zone: ""}, nil
}

func parseUDPAddr(addr string) (net.UDPAddr, error) {
	splitAddr := strings.SplitN(addr, ":", 2)
	ip, err := parseIP(splitAddr[0])
	if err != nil {
		return net.UDPAddr{}, fmt.Errorf("address parsing failed: %s as UDP: %v", addr, err)
	}
	port, err := parsePort(splitAddr[1])
	if err != nil {
		return net.UDPAddr{}, fmt.Errorf("address parsing failed: %s as UDP: %v", addr, err)
	}
	return net.UDPAddr{IP: ip, Port: port, Zone: ""}, nil
}

// validates and converts a string containing a Port
func parsePort(port string) (int, error) {
	portValue, err := strconv.Atoi(port)
	if err != nil {
		return -1, fmt.Errorf("port parsing failed: %s as port: %v", port, err)
	}
	return portValue, nil
}

// validates and converts a string containing an IP
func parseIP(ip string) (net.IP, error) {
	addr := []byte{}
	for _, addrPart := range strings.Split(ip, ".") {
		parsedValue, err := strconv.Atoi(addrPart)
		if err != nil {
			return nil, fmt.Errorf("ip parsing failed: %s as ip: %v", ip, err)
		}
		addr = append(addr, byte(parsedValue))
	}
	if len(addr) != 4 {
		return nil, fmt.Errorf("ip parsing failed: %s as ip", ip)
	}
	return net.IPv4(addr[0], addr[1], addr[2], addr[3]), nil
}
