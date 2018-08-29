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

	dataChan     chan EspData
	errorHandler func(error)

	prepareData EspData

	sendAcc, sendGyroAngle, sendAccAngle bool
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

	client.buf = make([]byte, 120)

	client.isStarted = false
	client.isDone = false

	client.dataChan = make(chan EspData, 1)
	client.errorHandler = func(err error) {
		log.Printf("%#v", err)
	}

	return client, nil
}

// SetErrCalback set func for handle errors
func (esp *EspClient) SetErrCalback(f func(err error)) {
	esp.errorHandler = f
}

// HandleErr set func for handle errors
func (esp *EspClient) HandleErr(err error) {
	esp.errorHandler(err)
}

// Start is user for start gorutines to correct work
func (esp *EspClient) Start() error {
	if !esp.isStarted {
		conn, err := net.DialUDP("udp", nil, esp.espAddr)
		if err != nil {
			return err
		}
		esp.udpConn = conn

		go esp.handleIncomingCommand()
		// send hello command every 1.5 s
		go func() {
			counter := 0
			for esp.IsWorking() {
				esp.SendCommand("hello", strconv.Itoa(counter))
				time.Sleep(time.Millisecond * 1500)
			}

		}()

		go esp.handleDriveRef()

	}

	esp.isStarted = true
	return nil
}

// EnableConf is used to enable sending additional data
// like send axeleration or giro angle
func (esp *EspClient) EnableConf(confConst espConf) {
	switch confConst {
	case SendAcc:
		esp.sendAcc = true
	case SendAccAngle:
		esp.sendAccAngle = true
	case SendGyroAngle:
		esp.sendGyroAngle = true
	}
}

// DisableConf is used to disable sending additional data
// like send axeleration or giro angle
func (esp *EspClient) DisableConf(confConst espConf) {
	switch confConst {
	case SendAcc:
		esp.sendAcc = false
	case SendAccAngle:
		esp.sendAccAngle = false
	case SendGyroAngle:
		esp.sendGyroAngle = false
	}
}

func (esp *EspClient) sendConf() {
	esp.SendCommand("SendAcc", boolTo1or0(esp.sendAcc))
	esp.SendCommand("SendGyroAngle", boolTo1or0(esp.sendGyroAngle))
	esp.SendCommand("SendAccAngle", boolTo1or0(esp.sendAccAngle))

}

// SendCommand is used for send command and data
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
			esp.HandleErr(err)
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
				esp.HandleErr(&ErrParceData{command, err})
				continue
			}
			esp.prepareData.AngleX = X
			esp.prepareData.AngleY = Y
			esp.prepareData.AngleZ = Z

			// fmt.Printf("x %4.2f y %4.2f z%4.2f\n", rX, rY, rZ)

		case "SendAcc":
			X, Y, Z, err := parceXYZ(dataStr)
			if err != nil {
				esp.HandleErr(&ErrParceData{command, err})
				continue
			}
			esp.prepareData.AccX = X
			esp.prepareData.AccY = Y
			esp.prepareData.AccZ = Z
		case "SendAccAngle":
			X, Y, Z, err := parceXYZ(dataStr)
			if err != nil {
				esp.HandleErr(&ErrParceData{command, err})
				continue
			}
			esp.prepareData.AAngleX = X
			esp.prepareData.AAngleX = Y
			esp.prepareData.AAngleX = Z
		case "SendGyroAngle":
			X, Y, Z, err := parceXYZ(dataStr)
			if err != nil {
				esp.HandleErr(&ErrParceData{command, err})
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
			esp.HandleErr(&ErrUndefinnedCommand{string(buf[:n])})
		}
	}
	esp.HandleErr(&ErrEspIsDone{})
}

// GetDataChan return chan with data from esp
func (esp *EspClient) GetDataChan() <-chan EspData {
	return esp.dataChan
}

// SetDriveRef is used for set reference on drive
func (esp *EspClient) SetDriveRef(dr1, dr2 int) {
	esp.mutex.Lock()
	esp.dr1 = dr1
	esp.dr2 = dr2
	esp.mutex.Unlock()
}

func (esp *EspClient) handleDriveRef() {
	tickerDriveRef := time.NewTicker(time.Millisecond * 10)
	defer tickerDriveRef.Stop()

	tickerSetOtherSetup := time.NewTicker(time.Millisecond * 1000)
	defer tickerSetOtherSetup.Stop()

	for esp.IsWorking() {
		select {
		case <-tickerDriveRef.C:
			esp.mutex.Lock()
			dr1 := esp.dr1
			dr2 := esp.dr2
			esp.mutex.Unlock()

			esp.SendCommand("dr1", strconv.Itoa(dr1))
			esp.SendCommand("dr2", strconv.Itoa(dr2))
		case <-tickerSetOtherSetup.C:
			esp.sendConf()
		}

	}
}

// IsWorking return true if driver on work
func (esp *EspClient) IsWorking() bool {
	esp.mutex.Lock()
	defer esp.mutex.Unlock()
	return esp.isStarted && !esp.isDone
}

// IsDone return true if Stop was called
func (esp *EspClient) IsDone() bool {
	esp.mutex.Lock()
	defer esp.mutex.Unlock()
	return esp.isDone
}

// IsConnected return false if more than 0.5s hasn`t got messa
func (esp *EspClient) IsConnected() bool {
	esp.mutex.Lock()
	defer esp.mutex.Unlock()
	return time.Now().Sub(esp.lastMessageTime).Seconds() <= 0.5
}

// Stop is used for stop sending reference data
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
	time.Sleep(10 * time.Millisecond)
	esp.udpConn.Close()
	esp.isDone = true
}
