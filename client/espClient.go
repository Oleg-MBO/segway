package client

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type EspClient struct {
	servAddr, espAddr *net.UDPAddr
	udpConn           *net.UDPConn
	rX, rY, rZ        float64
	isDone            bool
	isInitializated   bool
	mutex             sync.Mutex
	lastMessageTime   time.Time
	buf               []byte
}

// "192.168.0.110:4210"
func NewEspClient(address string) (*EspClient, error) {
	client := new(EspClient)
	ESPAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return client, err
	}
	client.espAddr = ESPAddr

	ServerAddr, err := net.ResolveUDPAddr("udp", ":")
	if err != nil {
		return client, err
	}
	client.servAddr = ServerAddr

	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	if err != nil {
		return client, err
	}
	client.udpConn = ServerConn

	client.buf = make([]byte, 120)

	client.isInitializated = true
	client.isDone = false
	go client.handleIncomingCommand()
	// send hello command every 1.5 s
	go func() {
		counter := 0
		for client.IsWorking() {
			client.SendCommand("hello", strconv.Itoa(counter))
			time.Sleep(time.Millisecond * 1500)
		}
	}()
	return client, nil
}
func (esp *EspClient) IsWorking() bool {
	esp.mutex.Lock()
	defer esp.mutex.Unlock()
	return esp.isInitializated && !esp.isDone
}

func (esp *EspClient) IsDone() bool {
	esp.mutex.Lock()
	defer esp.mutex.Unlock()
	return esp.isDone
}

func (esp *EspClient) IsConnected() bool {
	esp.mutex.Lock()
	defer esp.mutex.Unlock()
	return time.Now().Sub(esp.lastMessageTime).Seconds() <= 0.5
}

func (esp *EspClient) Stop() {
	fmt.Println(1)
	esp.mutex.Lock()
	defer esp.mutex.Unlock()
	esp.udpConn.Close()
	fmt.Println(2)
	esp.isDone = true
}

func (esp *EspClient) SendCommand(command, data string) {
	// fmt.Println("SendCommand")
	msg := fmt.Sprintf(">%s:%s<", command, data)

	buf1 := []byte(msg)

	n, err := esp.udpConn.WriteToUDP(buf1, esp.espAddr)
	if err != nil {
		fmt.Println(msg, err)
	}
	fmt.Println(n)
}

var reSendAnglesData = regexp.MustCompile(`\w{2} (.?\d+.\d+)`)

func (esp *EspClient) handleIncomingCommand() {
	buf := esp.buf

	for esp.IsDone() {
		n, _, err := esp.udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// if n == 0 {
		// 	continue
		// }
		fmt.Println("handleIncomingCommand")
		fmt.Println(string(buf[:n]))

		esp.mutex.Lock()
		esp.lastMessageTime = time.Now()
		esp.mutex.Unlock()

		startTagPos := bytes.Index(buf, []byte(">"))
		dataPos := bytes.Index(buf, []byte(":"))

		endTagPos := bytes.Index(buf, []byte("<"))
		if startTagPos == -1 || endTagPos == -1 || endTagPos > n || dataPos == -1 {
			// this is not command
			continue
		}
		command := string(buf[startTagPos+1 : dataPos])
		data := make([]byte, endTagPos-dataPos, endTagPos-dataPos)
		copy(data[:], buf[dataPos+1:endTagPos])

		switch command {
		case "SendAngles":
			submatch := reSendAnglesData.FindAllStringSubmatch(string(data), -1)

			if len(submatch) != 3 {
				log.Printf("error parse command %s with data %s\n", command, string(data))
			}
			rX, err := strconv.ParseFloat(submatch[0][1], 64)
			if err != nil {
				log.Println(err)
				continue
			}
			rY, err := strconv.ParseFloat(submatch[1][1], 64)
			if err != nil {
				log.Println(err)
				continue
			}
			rZ, err := strconv.ParseFloat(submatch[2][1], 64)
			if err != nil {
				log.Println(err)
				continue
			}
			esp.rX = rX
			esp.rY = rY
			esp.rZ = rZ

			fmt.Printf("x %4.2f y %4.2f z%4.2f\n", rX, rY, rZ)

		case "dr1":
			fmt.Println(command, string(data))
		case "dr2":
			fmt.Println(command, string(data))
		default:

			fmt.Printf("undefinded message: %s\n", string(buf[:n]))
		}
	}
}

func (esp *EspClient) GetRotatePos() (rX, rY, rZ float64) {
	esp.mutex.Lock()
	defer esp.mutex.Unlock()
	return esp.rX, esp.rY, esp.rZ
}
