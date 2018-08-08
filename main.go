package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
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

	// apLinkT := 100.0
	apLinkT := 1 / 0.01
	fmt.Println(apLinkT)
	apLink1 := tool.NewAperiodicLink(apLinkT)
	apLink2 := tool.NewAperiodicLink(apLinkT)

	apLinkOut := tool.NewAperiodicLink(apLinkT)

	rYOffset := 2.0

	pkoef := 75.0
	pkoef = 199.7 / 2

	dkoef := 1.0
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

	ticker := time.NewTicker(time.Millisecond * 10)
	printTicker := time.NewTicker(time.Millisecond * 25)

	var rY float64
	for segway.IsWorking() {
		select {
		case <-ticker.C:
			dt := now.Sub(prev).Seconds()
			prev = now
			now = time.Now()
			if segway.IsConnected() {
				_, rY, _ = segway.GetRotatePos()
				rY = 95.0 + rYOffset - rY

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

				drMax := 1000

				dr1 := int(pid.Output())
				// dr1 = 0
				if dr1 > drMax {
					dr1 = drMax
				}
				if dr1 < -drMax {
					dr1 = -drMax
				}

				segway.SetDriveRef(dr1, dr1)

			}
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

		}
	}

}
