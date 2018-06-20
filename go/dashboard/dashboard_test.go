package dashboard_test

import (
	"fmt"
	stdlog "log"
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

func log(str string) {
	stdlog.Println(str)
	if dashboard.Active() {
		dashboard.UpdateLine(cnLog, str)
	}
}

func logf(format string, args ...interface{}) {
	log(fmt.Sprintf(format, args...))
}

func updateName() {
	names := []string{"Shiawassee", "Prairie", "Natasha", "Tess", "Grant"}
	length := len(names)
	for dashboard.Updateable() {
		dashboard.UpdateKeyVal(cnName, names[rand.Intn(length)])
		sleep(750)
	}
}

func updateCount() {
	var buf strings.Builder
	var format string
	var j int
	buf.Grow(1024)
	for dashboard.Updateable() {
		if j&1 == 0 {
			format = "[%d]"
		} else {
			format = "%d"
		}
		fmt.Fprintf(&buf, format, j)
		// logf("update %d %s", j, buf.String())
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
	for dashboard.Updateable() {
		fmt.Fprintf(&buf, "This is line %s %d", fillStr[:rand.Intn(32)+1], j)
		log(buf.String())
		buf.Reset()
		sleep(1250)
		j++
	}
}

func sleep(waitMs int) {
	time.Sleep(time.Duration(waitMs) * time.Millisecond)
}

func Example() {
	var fl *os.File
	var err error

	fl, err = os.Create("log")
	if err == nil {
		stdlog.SetFlags(stdlog.LstdFlags | stdlog.Lmicroseconds)
		stdlog.SetOutput(fl)
		dashboard.RegisterHeader(cnBannerA, 0, 0, 0, "\\tDashboard (v 0.2)")
		dashboard.RegisterKeyVal(cnCount, 1, 1, 40, "Count")
		dashboard.RegisterLine(cnLog, 1, 2, 5, "2006-Jan-02 15:04:05 ")
		dashboard.RegisterHeaderLine(cnBannerC, 1, 7, 40, "\\t Dog ")
		dashboard.RegisterKeyVal(cnName, 1, 8, 40, "Name")
		dashboard.RegisterHeader(cnBannerB, 0, 9, 0, " This is a banner \\t more work to do \\t Press Q to quit ")
		go updateCount()
		go updateName()
		go updateLog()
		err = dashboard.Run('Q', 'q', 27)
		stdlog.SetOutput(os.Stdout)
		fl.Close()
		fmt.Println("Success")
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
	// Output:
	// Success
}
