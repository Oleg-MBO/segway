package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

var (
	counter int64
)

func main() {
	ESPAddr, err := net.ResolveUDPAddr("udp", "192.168.0.110:4210")
	CheckError(err)

	ServerAddr, err := net.ResolveUDPAddr("udp", ":")
	CheckError(err)

	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	CheckError(err)
	defer ServerConn.Close()

	sendcommand := func(command, data string) {
		msg := fmt.Sprintf(">%s:%s<", command, data)

		buf1 := []byte(msg)

		_, err := ServerConn.WriteToUDP(buf1, ESPAddr)
		if err != nil {
			fmt.Println(msg, err)
		}
		// fmt.Println(msg)
	}

	sendcommand("hello", "hello")
	go func() {
		var buf []byte = make([]byte, 50)

		for {
			n, _, err := ServerConn.ReadFromUDP(buf)
			if err != nil {
				fmt.Println(err)
				continue
			}
			// fmt.Println(string(buf[0:n]))
			// from = time.Now()

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
				submatch := ReSendAnglesData.FindAllStringSubmatch(string(data), -1)
				// fmt.Println(string(data))
				// fmt.Println(len(submatch))
				// fmt.Println(submatch)
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
				GyroData.Update(rX, rY, rZ)
				// rX, rY, rZ = GyroData.Update(rX, rY, rZ)
				// GyroData.GetValues()
				fmt.Printf("x %4.2f y %4.2f z%4.2f\n", rX, rY, rZ)
				counter++
				if counter%10 == 0 {
					sendcommand("dr1", strconv.Itoa(int(400+rY*5)))

					sendcommand("dr2", strconv.Itoa(int(400+rY*5)))
				}

			case "dr1":
				fmt.Println(command, string(data))
			case "dr2":
				fmt.Println(command, string(data))
			default:

				fmt.Printf("undefinded message: %s\n", string(buf[:n]))
			}

		}
	}()
	i := 0
	for {

		sendcommand("hello", strconv.Itoa(i))
		i++

		time.Sleep(time.Millisecond * 1000 * 10)
	}
}
