package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"runtime/pprof"
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

	// flag.BoolVar(&SendAcc, "acc", false, "is used to enable sending angle from gyroscope")
	// flag.BoolVar(&SendGyroAngle, "gyroA", false, "is used to enable sending angle from gyroscope")
	// flag.BoolVar(&SendAccAngle, "accA", false, "is used to enable sending angle form accelerometer")
	// flag.BoolVar(&SendAllData, "all", false, "is used to enable sending ALL data from esp")

	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")

	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	// nothing to do with data
	writeToCsv := func(data client.EspData, otherData map[string]float64) {}

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

	segway, err := client.InitEspClient("192.168.173.150:4210")
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
		log.Printf("segway err: %v\n", err)
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
	pid := pidR.NewPIDSimple(pidR.ConfNewPIDSimple{
		P:    1024 / 5,
		D:    0.75,
		DLim: 500.0,
		I:    5.0,
		Ilim: 10.0,
	})

	printTicker := time.NewTicker(time.Millisecond * 100)
	defer printTicker.Stop()

	var rY float64

	otherDataToCsv := make(map[string]float64)

	segwayDataChan := segway.GetDataChan()

	var prevMilis uint64

	for segway.IsWorking() {
		select {
		case data := <-segwayDataChan:

			if prevMilis == 0 {
				prevMilis = data.Milis
				continue
			}
			dt := (float64(data.Milis) - float64(prevMilis)) / 1000
			prevMilis = data.Milis

			rY1 := data.AngleX
			rY = rY1 - 1.3

			apLink1.Update(dt, rY)
			apLink2.Update(dt, apLink1.Output())
			rY = apLink2.Output()

			if math.Abs(rY) > 60 {
				segway.SetDriveRef(0, 0)
				continue
			}

			apLinkOut.Update(dt, (rY))
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
			otherDataToCsv["p"] = pid.P.Output()
			otherDataToCsv["i"] = pid.I.Output()
			otherDataToCsv["d"] = pid.D.Output()
			otherDataToCsv["pid"] = pid.Output()
			writeToCsv(data, otherDataToCsv)

		case <-printTicker.C:

			if !segway.IsConnected() {
				fmt.Println("segway not connected")
				continue
			}

			if math.Abs(rY) > 60 {
				fmt.Println("|angle| > 60, angle ==", rY)
				continue
			}

			fmt.Printf("rY: %+7.3f | p:%+9.2f, d:%+9.2f, i:%+9.2f |pid:%+9.1f\n", rY, pid.P.Output(), pid.D.Output(), pid.I.Output(), pid.Output())
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
