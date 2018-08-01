package main

import (
	"fmt"
	"log"
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

	ticker := time.NewTicker(time.Millisecond * 500)
	for segway.IsWorking() {
		<-ticker.C
		if segway.IsConnected() {
			fmt.Println("connected")
		} else {
			fmt.Println("disconnected")
		}
	}
}
