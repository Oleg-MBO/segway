package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/Oleg-MBO/segway/client"
	pidR "github.com/Oleg-MBO/segway/control/pid"
	"github.com/Oleg-MBO/segway/control/tool"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var (
		csvFileName string

		SendAcc       bool
		SendGyroAngle bool
		SendAccAngle  bool
		SendAllData   bool
	)

	flag.StringVar(&csvFileName, "csv", "", "save data to csv file if specified")

	flag.BoolVar(&SendAcc, "acc", false, "is used to enable sending angle from gyroscope")
	flag.BoolVar(&SendGyroAngle, "gyroA", false, "is used to enable sending angle from gyroscope")
	flag.BoolVar(&SendAccAngle, "accA", false, "is used to enable sending angle form accelerometer")
	flag.BoolVar(&SendAllData, "all", false, "is used to enable sending ALL data from esp")

	flag.Parse()

	// nothing to do with data
	writeToCsv := func(data client.EspData, otherData map[string]float64) {

	}

	if csvFileName != "" {
		if !strings.HasSuffix(csvFileName, ".csv") {
			csvFileName += ".csv"
		}
		file, err := os.Create(csvFileName)
		checkErr(err)
		defer func() {
			log.Println("Closing file..")
			file.Close()
		}()

		csvWriter := NewEspDataCSVWriter(file)
		defer func() {
			log.Println("Flushing csv data..")
			csvWriter.Flush()
		}()

		writeToCsv = func(data client.EspData, otherData map[string]float64) {
			csvWriter.WriteEspData(data, otherData)
		}

	}

	segway, err := client.InitEspClient("192.168.0.110:4210")
	checkErr(err)

	// 	SendAcc       bool
	// SendGyroAngle bool
	// SendAccAngle  bool
	// SendAllData   bool
	if SendAcc {
		segway.EnableConf(client.SendAcc)
	}

	if SendGyroAngle {
		segway.EnableConf(client.SendGyroAngle)
	}

	if SendAccAngle {
		segway.EnableConf(client.SendAccAngle)
	}

	if SendAllData {
		segway.EnableConf(client.SendAcc)
		segway.EnableConf(client.SendGyroAngle)
		segway.EnableConf(client.SendAccAngle)
	}

	err = segway.Start()
	checkErr(err)

	segwayErrChan := make(chan error, 1)
	segway.SetErrCalback(func(err error) {
		log.Println("segway err:", err)
		segwayErrChan <- err
	})

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Println("\nReceived an interrupt, stopping services...")
			segway.Stop()
			fmt.Println("segway.Stop() executed")
			// os.Exit(0)
		}
	}()

	// apLinkT := 100.0
	// apLinkT := 1 / 0.01
	apLinkT := 1 / 0.05

	fmt.Println(apLinkT)
	apLink1 := tool.NewAperiodicLink(apLinkT)
	apLink2 := tool.NewAperiodicLink(apLinkT)

	apLinkOut := tool.NewAperiodicLink(apLinkT)

	// rYOffset := 2.0

	pkoef := 75.0
	pkoef = 1024 / 5

	dkoef := 0.75
	// dkoef = 0.0
	dlim := 500.0

	ikoef := 5.0
	ilim := 10.0

	ikoef = 0

	p := pidR.NewPRegul(pkoef)
	d := pidR.NewDRegul(dkoef, dlim)
	i := pidR.NewIRegul(ikoef, ilim)
	pid := pidR.NewPID(p, d, i)

	now := time.Now()
	prev := time.Now()

	// ticker := time.NewTicker(time.Millisecond * 10)
	printTicker := time.NewTicker(time.Millisecond * 100)

	var rY float64

	otherDataToCsv := make(map[string]float64)

	segwayDataChan := segway.GetDataChan()
	for segway.IsWorking() {
		select {
		case data := <-segwayDataChan:
			dt := now.Sub(prev).Seconds()
			prev = now
			now = time.Now()

			// _, rY, _ = segway.GetRotatePos()
			rY = data.AngleY
			// rY = 95.0 + rYOffset - rY
			rY = 92.1 - rY

			apLink1.Update(dt, rY)
			apLink2.Update(dt, apLink1.Output())
			rY = apLink2.Output()

			// blind := 1.0
			// if (rY > 0 && rY < blind) || (rY < 0 && rY > -blind) {
			// 	rY = 0
			// }

			if math.Abs(rY) > 60 {
				segway.SetDriveRef(0, 0)
				continue
			}

			apLinkOut.Update(dt, -(rY))
			pid.Update(dt, apLinkOut.Output())

			drMax := 2000

			dr1 := int(pid.Output())
			// dr1 = 0
			if dr1 > drMax {
				dr1 = drMax
			}
			if dr1 < -drMax {
				dr1 = -drMax
			}

			segway.SetDriveRef(dr1, dr1)

			// write data to csv
			otherDataToCsv["p"] = p.Output()
			otherDataToCsv["i"] = i.Output()
			otherDataToCsv["d"] = d.Output()
			otherDataToCsv["pid"] = pid.Output()
			writeToCsv(data, otherDataToCsv)

		case <-printTicker.C:

			if !segway.IsConnected() {
				fmt.Println("segway not connected")
				continue
			}

			if math.Abs(rY) > 60 {
				segway.SetDriveRef(0, 0)
				fmt.Println("|angle| > 60, angle ==", rY)
				continue
			}

			fmt.Printf("rY: %+7.3f | p:%+9.2f, d:%+9.2f, i:%+9.2f |pid:%+9.1f\n", rY, p.Output(), d.Output(), i.Output(), pid.Output())
		case err := <-segwayErrChan:
			switch err.(type) {
			case *client.ErrEspIsDone:
				// if esp is done
				// 	// break the loop
				break
			}

		}

	}
}
