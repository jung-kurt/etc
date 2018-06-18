package dashboard_test

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/jung-kurt/etc/go/dashboard"
)

const (
	cnCount int = iota
	cnName
	cnLog
	cnBannerA
	cnBannerB
	cnBannerC
)

func updateName() {
	names := []string{"Shiawassee", "Prairie", "Natasha", "Tess", "Grant"}
	length := len(names)
	for dashboard.Active() {
		dashboard.UpdateKeyVal(cnName, names[rand.Intn(length)])
		sleep(750)
	}
}

func updateCount() {
	var buf strings.Builder
	var format string
	var j int
	buf.Grow(1024)
	for dashboard.Active() {
		if j&1 == 0 {
			format = "[%d]"
		} else {
			format = "%d"
		}
		fmt.Fprintf(&buf, format, j)
		// log.Printf("update %d %s", j, buf.String())
		dashboard.UpdateKeyVal(cnCount, buf.String())
		buf.Reset()
		sleep(1000)
		j++
	}
}

func updateLog() {
	var buf strings.Builder
	var j int
	buf.Grow(256)
	fillStr := strings.Repeat("x", 256)
	for dashboard.Active() {
		fmt.Fprintf(&buf, "This is log line %s %d", fillStr[:rand.Intn(32)+1], j)
		dashboard.UpdateLine(cnLog, buf.String())
		buf.Reset()
		sleep(1250)
		j++
	}
}

func Example() {
	var fl *os.File
	var err error

	fl, err = os.Create("log")
	if err == nil {
		log.SetOutput(fl)
		dashboard.RegisterHeader(cnBannerA, 0, 0, 0, "\\tDashboard (v 0.2)")
		dashboard.RegisterKeyVal(cnCount, 1, 1, 40, "Count")
		dashboard.RegisterLine(cnLog, 1, 2, 5)
		dashboard.RegisterHeaderLine(cnBannerC, 1, 7, 40, "\\t Dog ")
		dashboard.RegisterKeyVal(cnName, 1, 8, 40, "Name")
		dashboard.RegisterHeader(cnBannerB, 0, 9, 0, " This is a banner \\t more work to do \\t Press Q to quit ")
		go updateCount()
		go updateName()
		go updateLog()
		// sleep(3000) // test
		err = dashboard.Run('Q', 'q', 27)
		// sleep(3000) // test
		log.SetOutput(os.Stdout)
		fmt.Printf("Hello, world\n")
		fl.Close()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
	// Output:
	// Hello, world
}

func sleep(waitMs int) {
	time.Sleep(time.Duration(waitMs) * time.Millisecond)
}

//
//
// func timeChange(updateChan chan<- updateType) {
// 	tm := time.Now().Add(-time.Hour * 24 * 365 * 50).Unix()
// 	for {
// 		updateChan <- updateType{id: updateTime, str: time.Unix(tm, 0).Format("2 Jan 2006, 15:04:03")}
// 		tm += rand.Int63n(60 * 60 * 24)
// 		sleep(1250)
// 	}
// }
//
// func walkChange(updateChan chan<- updateType) {
// 	for {
// 		updateChan <- updateType{id: updateWalk, ok: rand.Intn(5) > 0}
// 		sleep(1500)
// 	}
// }
//
// func logWrite(updateChan chan<- updateType) {
// 	var buf bytes.Buffer
// 	var count int
//
// 	buf.Grow(4 * cnMaxWidth)
// 	for {
// 		count++
// 		fmt.Fprintf(&buf, "This is log line %d", count)
// 		updateChan <- updateType{id: updateLog, str: buf.String()}
// 		buf.Reset()
// 		sleep(1000)
// 	}
// }
//
// func main() {
// 	var screen tcell.Screen
// 	var err error
//
// 	// prf := profile.Start(profile.MemProfile)
// 	screen, err = tcell.NewScreen()
// 	if err == nil {
// 		encoding.Register()
// 		err = screen.Init()
// 		if err == nil {
// 			screen.HideCursor()
// 			quit := make(chan struct{})
// 			update := make(chan updateType, 128)
// 			go listen(update, quit, screen) // listen for terminate event
// 			go nameChange(update)           // generate names
// 			go timeChange(update)           // generate times
// 			go walkChange(update)           // generate walk index events
// 			go logWrite(update)             // Generate log lines
// 			go show(screen, update)         // display summary of events as they are received
// 			<-quit                          // terminate event occurred
// 			screen.Fini()
// 		}
// 	}
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "%s", err)
// 	}
// 	// prf.Stop()
// }
