package dashboard

import (
	"bytes"
	"log"
	"strings"
	"sync"

	"github.com/gdamore/tcell"
)

// The update channel is used (with the DashboardUpdateXXX() functions as
// wrappers) to update various dashboard fields. It is ready to receive records
// as soon as the application is initialized. It is kept open through the
// termination of the application to prevent panics if the application updates
// a field after the main dashboard loop has completed.

var dsh struct {
	dotStr     string               // separation for key/value fields
	diamondStr string               // overflow indictor
	blankStr   string               // overflow blank string to clear characters of previous string
	buf        *bytes.Buffer        // formatted string buffer
	fieldMap   map[int]fieldPtrType // fields registered by application
	updateChan chan updateType      // dashboard updates fields based on events arriving in this channel
	screen     tcell.Screen         // terminal screen
	active     bool                 // update channel is active
}

type itemType int

const (
	itemKeyVal itemType = iota
	itemStatic
	itemLine
	itemWalk
)

const (
	updateScreen int = iota // internal flag must be set
	updateStop
)

type updateType struct {
	internal bool   // true if id is defined internally
	id       int    // application defined field identifier
	str      string // string value
	ok       bool   // flag for walk line
}

type fieldType struct {
	id      int               // user-defined field identifier
	item    itemType          // type of field key-value, static, etc
	str     string            // static string, may include tabs for left, center, right alignment
	x, y    int               // starting position
	wd      int               // width of field, 0 for entire line
	strList []strings.Builder // series of strings for rolling logs
	pos     int               // walk position; for log series, next strList position to fill
	count   int               // number of builders assigned in strList
}

type fieldPtrType *fieldType

const (
	cnWalkWidth = 3
	cnMaxWidth  = 256
)

func init() {
	var err error

	dsh.dotStr = strings.Repeat(".", cnMaxWidth)
	dsh.diamondStr = strings.Repeat(string(tcell.RuneDiamond), cnMaxWidth)
	dsh.blankStr = strings.Repeat(" ", cnMaxWidth)
	dsh.buf = &bytes.Buffer{}
	dsh.buf.Grow(4 * cnMaxWidth)
	dsh.fieldMap = make(map[int]fieldPtrType)
	dsh.updateChan = make(chan updateType, 256)
	dsh.screen, err = tcell.NewScreen()
	dsh.active = true
	if err == nil {
		// encoding.Register() // Asian encodings, adds several megabytes to application size
	}
	if err != nil {
		panic(err)
	}
}

func listen(updateChan chan<- updateType, pollEvent func() tcell.Event, quitRunes ...rune) {
	var rn rune
	quitMap := make(map[rune]bool)
	for _, rn = range quitRunes {
		quitMap[rn] = true
	}
	loop := true
	for loop {
		ev := pollEvent()
		loop = ev != nil
		if loop {
			rn = 0
			// log.Printf("event")
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape:
					rn = 27
				case tcell.KeyCtrlL:
					updateChan <- updateType{internal: true, id: updateScreen}
				case tcell.KeyRune:
					rn = ev.Rune()
				}
			case *tcell.EventResize:
				updateChan <- updateType{internal: true, id: updateScreen}
			}
			if rn > 0 {
				quit, ok := quitMap[rn]
				if ok && quit {
					loop = false
				}
			}
		}
	}
	updateChan <- updateType{internal: true, id: updateStop}
}

func put(style tcell.Style, x, y, scrWd int, strs ...string) {
	for _, str := range strs {
		for _, r := range str {
			if x < scrWd {
				dsh.screen.SetContent(x, y, r, nil, style)
				x++
			}
		}
	}
}

func keyval(styleKey, styleVal tcell.Style, x, y, wd, scrWd int,
	keyStr string, valStr string) {
	dsh.buf.Reset()
	keyLen := len(keyStr)
	valLen := len(valStr)
	if wd > cnMaxWidth {
		wd = cnMaxWidth
	}
	if keyLen+valLen+4 <= wd {
		put(styleKey, x, y, scrWd, keyStr, " ", dsh.dotStr[:wd-2-keyLen-valLen], " ")
		put(styleVal, x+wd-valLen, y, scrWd, valStr)
	} else {
		put(styleKey, x, y, scrWd, dsh.diamondStr[:wd])
	}
}

func walk(screen tcell.Screen, plainStyle, blockStyle tcell.Style, y, pos, wd int) {
	offPos := pos + cnWalkWidth
	var outerRune, innerRune rune
	var style tcell.Style
	for x := 0; x < wd; x++ {
		if x >= pos && x < offPos {
			outerRune = tcell.RuneBlock
			innerRune = tcell.RuneBlock
			style = blockStyle
		} else {
			outerRune = ' '
			innerRune = tcell.RuneHLine
			style = plainStyle
		}
		screen.SetContent(x, y, outerRune, nil, style)
		screen.SetContent(x, y+1, innerRune, nil, style)
		screen.SetContent(x, y+2, outerRune, nil, style)
	}
}

func staticRender(style tcell.Style) {
	scrWd, _ := dsh.screen.Size()
	var show bool
	for _, fieldPtr := range dsh.fieldMap {
		if fieldPtr.item == itemStatic {
			var rt int
			lf := fieldPtr.x
			wd := fieldPtr.wd
			length := len(fieldPtr.str)
			if wd <= 0 {
				rt = scrWd + wd
			} else {
				rt = lf + wd
			}
			if lf+length < rt {
				put(style, fieldPtr.x, fieldPtr.y, rt-lf, fieldPtr.str, dsh.blankStr[:rt-length-lf])
			} else {
				put(style, fieldPtr.x, fieldPtr.y, rt-lf, fieldPtr.str[:rt-lf])
			}
			show = true
		}
	}
	if show {
		dsh.screen.Show()
	}
}

func run() {
	// var logList [cnLogCount]string
	// var logPos, logCount int
	// const left = 1
	plain := tcell.StyleDefault.Foreground(tcell.ColorYellow)
	white := tcell.StyleDefault.Foreground(tcell.ColorWhite)
	// red := tcell.StyleDefault.Foreground(tcell.ColorRed)
	// green := tcell.StyleDefault.Foreground(tcell.ColorGreen)
	// wd, ht := dsh.screen.Size()
	staticRender(white)
	loop := true
	// walkPos := 0
	// syncCount := 0
	for loop {
		up := <-dsh.updateChan
		if up.internal {
			// log.Printf("internal")
			switch up.id {
			case updateScreen:
				staticRender(white)
				dsh.screen.Sync()
			case updateStop:
				loop = false
			}
		} else {
			// log.Printf("external")
			// var st tcell.Style
			var fieldPtr fieldPtrType
			var ok bool
			fieldPtr, ok = dsh.fieldMap[up.id]
			if ok {
				scrWd, _ := dsh.screen.Size()
				// log.Printf("good field %d", up.id)
				switch fieldPtr.item {
				case itemKeyVal:
					// log.Printf("dsh.keyval x %d, y %d, wd %d, key %s, val %s", fieldPtr.x,
					// fieldPtr.y, fieldPtr.wd, fieldPtr.str, up.str)
					keyval(plain, white, fieldPtr.x, fieldPtr.y, fieldPtr.wd, scrWd, fieldPtr.str, up.str)
				case itemLine:
					// log.Printf("line [%s]", up.str)
					count := fieldPtr.count
					size := len(fieldPtr.strList)
					if count < size {
						count++
						fieldPtr.count = count
					}
					pos := fieldPtr.pos
					fieldPtr.strList[pos].Reset()
					fieldPtr.strList[pos].WriteString(up.str)
					pos++
					if pos >= size {
						pos = 0
					}
					fieldPtr.pos = pos
					for j := 0; j < count; j++ {
						k := pos + j
						if k >= count {
							k -= count
						}
						left := fieldPtr.x
						top := fieldPtr.y
						str := fieldPtr.strList[k].String()
						length := len(str)
						if length+left <= scrWd {
							put(white, left, top+j, scrWd, str, dsh.blankStr[:scrWd-length-left])
						} else {
							put(white, left, top+j, scrWd, str[:scrWd-left-2], "..")
						}
					}
				case itemWalk:
				}
				dsh.screen.Show()
			}
		}

		// switch up.id {
		// case updateName:
		// 	dsh.keyval(plain, white, left, 1, 40, "Name", "%s", up.str)
		// case updateTime:
		// 	dsh.keyval(plain, white, left, 9, 40, "Time", "%s", up.str)
		// case updateLog:
		// 	if logCount < cnLogCount {
		// 		logCount++
		// 	}
		// 	logList[logPos] = up.str
		// 	logPos++
		// 	if logPos >= cnLogCount {
		// 		logPos = 0
		// 	}
		// 	for j := 0; j < logCount; j++ {
		// 		k := logPos + j
		// 		if k >= logCount {
		// 			k -= logCount
		// 		}
		// 		str := logList[k]
		// 		length := len(str)
		// 		if length+left <= wd {
		// 			put(dsh.screen, white, left, 11+j, str, blankStr[:wd-length-left])
		// 		} else {
		// 			put(dsh.screen, white, left, 11+j, str[:wd-left-2], "..")
		// 		}
		// 	}
		// case updateWalk:
		// 	if up.ok {
		// 		st = green
		// 	} else {
		// 		st = red
		// 	}
		// 	walk(dsh.screen, plain, st, 5, walkPos, wd)
		// 	walkPos += cnWalkWidth
		// 	if walkPos+cnWalkWidth >= wd {
		// 		walkPos = 0
		// 	}
		// 		case updateScreen:
		// 			wd, ht = dsh.screen.Size()
		// 			syncCount++
		// 			dsh.keyval(dsh.screen, plain, white, left, 2, 40, "Width", "%d", wd)
		// 			dsh.keyval(dsh.screen, plain, white, left, 3, 40, "Height", "%d", ht)
		// 			dsh.keyval(dsh.screen, plain, white, left, 10, 40, "Sync", "%d", syncCount)
		// 			dsh.screen.Sync()
		// }
	}
	// close(dsh.quitChan)
}

func register(id int, fldPtr fieldPtrType) {
	var mtx sync.Mutex
	mtx.Lock()
	dsh.fieldMap[id] = fldPtr
	mtx.Unlock()
}

// RegisterLine registers a dashboard rolling line field with the identifier
// specified by id. Its coordinates are specified by x and y, and the total
// number of rows used is specified by lineCount.
func RegisterLine(id, x, y, lineCount int) {
	var fld fieldType

	fld.strList = make([]strings.Builder, lineCount)
	for j := 0; j < lineCount; j++ {
		fld.strList[j].Grow(cnMaxWidth)
	}
	fld.item = itemLine
	fld.id = id
	fld.x = x
	fld.y = y
	register(id, &fld)
}

// UpdateLine updates the rolling line field specified by id with the str.
func UpdateLine(id int, str string) {
	dsh.updateChan <- updateType{id: id, str: str}
}

// RegisterKeyVal registers a dashboard key/value pair with the identifier
// specified by id. Its coordinates are specified by x and y, and the total
// field's width is specified by wd. The static key is specified by keyStr.
func RegisterKeyVal(id, x, y, wd int, keyStr string) {
	register(id, &fieldType{id: id, item: itemKeyVal, x: x, y: y, wd: wd, str: keyStr})
}

// UpdateKeyVal updates the key/value pair specified by id with the value
// specified by str.
func UpdateKeyVal(id int, str string) {
	dsh.updateChan <- updateType{id: id, str: str}
}

// RegisterStatic registers a dashboard static line with the identifier
// specified by id. Its coordinates are specified by x and y. The total field's
// width is specified by wd. A zero value for wd indicates the full width of
// the screen, and a negative value indicates the position from the right. The
// static key is specified by keyStr.
func RegisterStatic(id, x, y, wd int, keyStr string) {
	register(id, &fieldType{id: id, item: itemStatic, x: x, y: y, wd: wd, str: keyStr})
}

// UpdateStatic updates the key/value pair specified by id with the value
// specified by str.
// func UpdateStatic(id int, str string) {
//	dsh.updateChan <- updateType{id: id, str: str}
// }

// Run changes the screen to a dashboard. This function does not return until
// one of the keys included in the list of quitRunes is pressed. Up until that
// time, all application logic should be handled in other goroutines that call
// one or more of the UpdateXXX() functions to update the dashboard.
func Run(quitRunes ...rune) (err error) {
	err = dsh.screen.Init()
	// log.Printf("start, err == nil: %v", err == nil)
	if err == nil {
		// log.Printf("hide cursor")
		dsh.screen.HideCursor()
		go listen(dsh.updateChan, dsh.screen.PollEvent, quitRunes...)
		run()
		active(true, false)
		log.Printf("stop\n")
		dsh.screen.Fini()
	}
	return
}

func active(assign, active bool) (ret bool) {
	var mtx sync.Mutex
	mtx.Lock()
	ret = dsh.active
	if assign {
		dsh.active = active
	}
	mtx.Unlock()
	return
}

// Active returns true if the dashboard is currently active. It may be called
// safely from other goroutines. It is typically used in application loops.
func Active() bool {
	return active(false, false)
}
