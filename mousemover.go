package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procSetCursorPos     = user32.NewProc("SetCursorPos")
	procGetCursorPos     = user32.NewProc("GetCursorPos")
	procGetSystemMetrics = user32.NewProc("GetSystemMetrics")
)

type POINT struct {
	X, Y int32
}

func setCursorPos(x, y int) {
	procSetCursorPos.Call(uintptr(x), uintptr(y))
}

func getCursorPos() (int, int) {
	var p POINT
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&p)))
	return int(p.X), int(p.Y)
}

func getScreenSize() (int, int) {
	smCXSCREEN := 0
	smCYSCREEN := 1
	x, _, _ := procGetSystemMetrics.Call(uintptr(smCXSCREEN))
	y, _, _ := procGetSystemMetrics.Call(uintptr(smCYSCREEN))
	return int(x), int(y)
}

func main() {
	screenW, screenH := getScreenSize()
	fmt.Printf("Mouse walking on screen %dx%d\n", screenW, screenH)

	x, y := getCursorPos()
	prevX, prevY := x, y

	moveX := 1
	moveY := 1

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		fmt.Println("\nInterrupted, exiting...")
		os.Exit(0)
	}()

	for {
		x += moveX
		y += moveY

		setCursorPos(x, y)

		if x >= screenW-1 {
			moveX = -1
		} else if x <= 0 {
			moveX = 1
		}
		if y >= screenH-1 {
			moveY = -1
		} else if y <= 0 {
			moveY = 1
		}

		curX, curY := getCursorPos()
		if abs(curX-prevX) != abs(curY-prevY) {
			fmt.Println("User moved mouse — exiting.")
			break
		}

		prevX, prevY = curX, curY
		time.Sleep(time.Millisecond)
	}
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
