package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"time"

	"github.com/Oleg-MBO/segway/client"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	segway, error := client.NewEspClient("192.168.0.110:4210")
	checkErr(error)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		fmt.Println("\nReceived an interrupt, stopping services...")
		segway.Stop()
		fmt.Println("segway.Stop() executed")
		os.Exit(0)
	}()

	ticker := time.NewTicker(time.Millisecond * 20)
	for segway.IsWorking() {
		<-ticker.C
		if segway.IsConnected() {
			rX, _, _ := segway.GetRotatePos()
			fmt.Println(rX)
			if math.Abs(rX) > 21 {
				segway.SetDriveRef(0, 0)
				continue
			}
			offset := float64(3.2)
			dr1 := int((rX - offset) * ((1024) / 5))
			// dr1 := 0
			// dr1 := int(math.Sinh((rX - offset)) * 1023)
			segway.SetDriveRef(dr1, dr1)
		}
	}
}
