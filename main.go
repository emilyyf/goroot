package main

import (
	"github.com/daviseidel/xlib"
	"time"
)

func main() {
	disp := xlib.XOpenDisplay(0)
	root := xlib.XDefaultRootWindow(disp)
	for {
		date := time.Now().Format("Mon Jan 2 2006 15:04:05.000000")
		xlib.XStoreName(disp, root, date)
		xlib.XFlush(disp)
		time.Sleep(time.Millisecond * (1000 / 60))
	}
}
