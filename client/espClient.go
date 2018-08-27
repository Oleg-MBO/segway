package client

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

// EspClient struct that represent client to segway on the esp
type EspClient struct {
	espAddr *net.UDPAddr
	udpConn *net.UDPConn

	// rotate pos
	// rX, rY, rZ float64
	// drive reference
	dr1, dr2 int

	isDone          bool
	isStarted       bool
	mutex           sync.Mutex
	lastMessageTime time.Time
	buf             []byte

	dataChan chan EspData

	prepareData EspData
}

// InitEspClient initialase esp client and return structure with client
// if address is empty default adress will be "192.168.0.110:4210"
func InitEspClient(address string) (*EspClient, error) {
	if address == "" {
		address = "192.168.0.110:4210"
	}
	client := new(EspClient)
	ESPAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return client, err
	}
	client.espAddr = ESPAddr

	conn, err := net.DialUDP("udp", nil, ESPAddr)
	if err != nil {
		return client, err
	}
	client.udpConn = conn

	client.buf = make([]byte, 120)

	client.isStarted = false
	client.isDone = false

	client.dataChan = make(chan EspData, 1)

	return client, nil
}

func (esp *EspClient) Start() {
	if !esp.isStarted {
		go esp.handleIncomingCommand()
		// send hello command every 1.5 s
		go func() {
			counter := 0
			for esp.IsWorking() {
				esp.SendCommand("hello", strconv.Itoa(counter))
				time.Sleep(time.Millisecond * 1500)
			}

		}()
	}

	go esp.handleDriveRef()
	esp.isStarted = false
}

func (esp *EspClient) SendCommand(command, data string) {
	msg := fmt.Sprintf(">%s:%s<", command, data)

	buf1 := []byte(msg)

	_, err := esp.udpConn.Write(buf1)
	if err != nil {
		log.Println(msg, err)
	}
}

func (esp *EspClient) handleIncomingCommand() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("handleIncomingCommand Recovered in f", r)
		}
	}()
	buf := esp.buf

	for !esp.IsDone() {
		n, _, err := esp.udpConn.ReadFromUDP(buf)

		if err != nil {
			fmt.Println(err)
			continue
		}
		// if n == 0 {
		// 	continue
		// }

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

		dataStr := string(data)

		switch command {
		case "SendAngles":
			X, Y, Z, err := parceXYZ(dataStr)
			if err != nil {
				log.Printf(`Error during parse command "%s" err:%#v`, command, err)
				continue
			}
			esp.prepareData.AngleX = X
			esp.prepareData.AngleY = Y
			esp.prepareData.AngleZ = Z

			// fmt.Printf("x %4.2f y %4.2f z%4.2f\n", rX, rY, rZ)

		case "SendAcc":
			X, Y, Z, err := parceXYZ(dataStr)
			if err != nil {
				log.Printf(`Error during parse command "%s" err:%#v`, command, err)
				continue
			}
			esp.prepareData.AccX = X
			esp.prepareData.AccY = Y
			esp.prepareData.AccZ = Z
		case "SendAccAngle":
			X, Y, Z, err := parceXYZ(dataStr)
			if err != nil {
				log.Printf(`Error during parse command "%s" err:%#v`, command, err)
				continue
			}
			esp.prepareData.AAngleX = X
			esp.prepareData.AAngleX = Y
			esp.prepareData.AAngleX = Z
		case "SendGyroAngle":
			X, Y, Z, err := parceXYZ(dataStr)
			if err != nil {
				log.Printf(`Error during parse command "%s" err:%#v`, command, err)
				continue
			}
			esp.prepareData.GyroX = X
			esp.prepareData.GyroX = Y
			esp.prepareData.GyroX = Z

		// case "dr1":
		// 	fmt.Println(command, string(data))
		// case "dr2":
		// 	fmt.Println(command, string(data))
		case "SendDataDone":
			// send to chan if can
			// and not wait if can`t
			select {
			case esp.dataChan <- esp.prepareData:
				esp.prepareData = EspData{}
			default:
			}
		default:
			fmt.Printf("undefinded message: %s\n", string(buf[:n]))
		}
	}
	log.Println("esp is done.")

}

// GetDataChan return chan with data from esp
func (esp *EspClient) GetDataChan() <-chan EspData {
	return esp.dataChan
}

func (esp *EspClient) SetDriveRef(dr1, dr2 int) {
	esp.mutex.Lock()
	esp.dr1 = dr1
	esp.dr2 = dr2
	esp.mutex.Unlock()
}

func (esp *EspClient) handleDriveRef() {

	for !esp.IsDone() {
		esp.mutex.Lock()
		dr1 := esp.dr1
		dr2 := esp.dr2
		esp.mutex.Unlock()

		esp.SendCommand("dr1", strconv.Itoa(dr1))
		esp.SendCommand("dr2", strconv.Itoa(dr2))
		time.Sleep(time.Millisecond * 5)
	}
}

func (esp *EspClient) IsWorking() bool {
	esp.mutex.Lock()
	defer esp.mutex.Unlock()
	return esp.isStarted && !esp.isDone
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
	esp.mutex.Lock()
	defer esp.mutex.Unlock()
	// send every comand several times
	// and wait to send
	esp.SendCommand("dr1", "0")
	esp.SendCommand("dr2", "0")
	esp.SendCommand("dr1", "0")
	esp.SendCommand("dr2", "0")
	esp.SendCommand("dr1", "0")
	esp.SendCommand("dr2", "0")
	esp.SendCommand("dr1", "0")
	esp.SendCommand("dr2", "0")
	time.Sleep(100 * time.Millisecond)
	esp.udpConn.Close()
	esp.isDone = true
}
