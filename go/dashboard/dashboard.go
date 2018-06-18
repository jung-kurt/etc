package dashboard

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/gdamore/tcell"
)

// The update channel is used (with the UpdateXXX() functions as wrappers) to
// update various dashboard fields. It is ready to receive records as soon as
// the application is initialized. It is kept open through the termination of
// the application to prevent panics if the application updates a field after
// the main dashboard loop has completed.

var dsh struct {
	dotStr     string               // separation for key/value fields
	diamondStr string               // overflow indictor
	blankStr   string               // overflow blank string to clear characters of previous string
	buf        *strings.Builder     // formatted string buffer; accessed only by 'dashboard show' goroutine
	fieldMap   map[int]fieldPtrType // fields registered by application
	fieldMtx   sync.Mutex           // mutex for accessing fieldMap
	updateChan chan updateType      // dashboard updates fields based on events arriving in this channel
	screen     tcell.Screen         // terminal screen
	active     bool                 // update channel is active
	activeMtx  sync.Mutex           // mutex for accessing the active flag
}

type itemType int

const (
	itemKeyVal itemType = iota
	itemHeader
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
	item    itemType          // type of field key-value, header, etc
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
	dsh.buf = &strings.Builder{}
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
				_, ok := quitMap[rn]
				if ok {
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

func headerPut(style tcell.Style, scrWd int, fld fieldPtrType) {
	list := strings.Split(fld.str, "\\t")
	var gapA, gapB, totalLen int
	for j := len(list); j < 3; j++ {
		list = append(list, "")
	}
	for j, str := range list {
		log.Printf("Field %d, string %d: [%s]", fld.id, j, str)
	}
	totalLen = len(list[0]) + len(list[1]) + len(list[2])
	wd := fld.wd
	if wd <= 0 {
		wd = scrWd + wd - fld.x
	}
	gap := wd - totalLen
	if gap < 2 {
		gapA = 1
		gapB = 1
	} else {
		gapA = gap / 2
		gapB = gap - gapA
	}
	dsh.buf.Reset()
	fmt.Fprintf(dsh.buf, "%s%s%s%s%s", list[0], dsh.blankStr[:gapA], list[1], dsh.blankStr[:gapB], list[2])
	str := dsh.buf.String()
	if len(str) > wd {
		str = str[:wd]
	}
	put(style, fld.x, fld.y, scrWd, str)

	// var rt int
	// lf := fld.x
	// wd := fld.wd
	// length := len(fld.str)
	// if wd <= 0 {
	// 	rt = scrWd + wd
	// } else {
	// 	rt = lf + wd
	// }
	// if lf+length < rt {
	// 	put(style, fld.x, fld.y, rt-lf, fld.str, dsh.blankStr[:rt-length-lf])
	// } else {
	// 	put(style, fld.x, fld.y, rt-lf, fld.str[:rt-lf])
	// }
}

func headerRender(style tcell.Style) {
	var list []fieldPtrType
	scrWd, _ := dsh.screen.Size()
	dsh.fieldMtx.Lock()
	for _, fieldPtr := range dsh.fieldMap {
		if fieldPtr.item == itemHeader {
			list = append(list, fieldPtr)
		}
	}
	dsh.fieldMtx.Unlock()
	for _, fieldPtr := range list {
		headerPut(style, scrWd, fieldPtr)
	}
	if len(list) > 0 {
		dsh.screen.Show()
	}
}

func run() {
	// var logList [cnLogCount]string
	// var logPos, logCount int
	// const left = 1
	plain := tcell.StyleDefault.Foreground(tcell.ColorYellow)
	white := tcell.StyleDefault.Foreground(tcell.ColorWhite)
	banner := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorYellow)
	// red := tcell.StyleDefault.Foreground(tcell.ColorRed)
	// green := tcell.StyleDefault.Foreground(tcell.ColorGreen)
	// wd, ht := dsh.screen.Size()
	// headerRender(banner)
	loop := true
	// walkPos := 0
	// syncCount := 0
	for loop {
		up := <-dsh.updateChan
		if up.internal {
			// log.Printf("internal")
			switch up.id {
			case updateScreen:
				headerRender(banner)
				dsh.screen.Sync()
			case updateStop:
				loop = false
			}
		} else {
			// log.Printf("external")
			// var st tcell.Style
			var fieldPtr fieldPtrType
			var ok bool
			dsh.fieldMtx.Lock()
			fieldPtr, ok = dsh.fieldMap[up.id]
			dsh.fieldMtx.Unlock()
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

func fieldRegister(id int, fldPtr fieldPtrType) {
	dsh.fieldMtx.Lock()
	dsh.fieldMap[id] = fldPtr
	dsh.fieldMtx.Unlock()
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
	fieldRegister(id, &fld)
}

// UpdateLine updates the rolling line field specified by id with the str.
func UpdateLine(id int, str string) {
	dsh.updateChan <- updateType{id: id, str: str}
}

// RegisterKeyVal registers a dashboard key/value pair with the identifier
// specified by id. Its coordinates are specified by x and y, and the total
// field's width is specified by wd. The static key is specified by keyStr.
func RegisterKeyVal(id, x, y, wd int, keyStr string) {
	fieldRegister(id, &fieldType{id: id, item: itemKeyVal, x: x, y: y, wd: wd, str: keyStr})
}

// UpdateKeyVal updates the key/value pair specified by id with the value
// specified by str.
func UpdateKeyVal(id int, str string) {
	dsh.updateChan <- updateType{id: id, str: str}
}

// RegisterHeader registers a dashboard static line with the identifier
// specified by id. Its coordinates are specified by x and y. The total field's
// width is specified by wd. A zero value for wd indicates the full width of
// the screen, and a negative value indicates the position from the right. The
// static key is specified by keyStr.
func RegisterHeader(id, x, y, wd int, keyStr string) {
	fieldRegister(id, &fieldType{id: id, item: itemHeader, x: x, y: y, wd: wd, str: keyStr})
}

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
		log.Printf("stop\n")
		dsh.screen.Fini()
	}
	return
}

// Active returns true if the dashboard is currently active. It may be called
// safely from other goroutines. It is typically used in application loops.
func Active() (active bool) {
	dsh.activeMtx.Lock()
	active = dsh.active
	dsh.activeMtx.Unlock()
	return
}
