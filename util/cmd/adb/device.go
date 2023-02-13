package adb

import (
	"log"
	"math"
)

const (
	// device default touch screen signal
	DEVICE_BTN_TOUCH_DOWN    = " 014a 00000001"
	DEVICE_BTN_TOUCH_UP      = " 014a 00000000"
	DEVICE_ABS_MT_POSITION_X = " 0035 "
	DEVICE_ABS_MT_POSITION_Y = " 0036 "
	// "VirtualBox USB Tablet" touch screen signal
	VBOX_BTN_TOUCH_DOWN    = " 0110 00000001"
	VBOX_BTN_TOUCH_UP      = " 0110 00000000"
	VBOX_ABS_MT_POSITION_X = "0003 0000 "
	VBOX_ABS_MT_POSITION_Y = "0003 0001 "
)

type Device struct {
	DeviceName string
	ScreenMode int
	Info       *DeviceInfo
	ResW       int64
	ResH       int64
	IsVM       bool
}

func (d *Device) PrintInfo() {
	log.Printf(`
Device name: 			%s
Device is VM: 			%t
System resulotion: 		%d*%d
Touch screen info:
	Xmin: 	%d
	Xmax: 	%d
	Ymin: 	%d
	Ymax: 	%d`, d.DeviceName, d.IsVM, d.ResW, d.ResH, d.Info.Xmin, d.Info.Xmax, d.Info.Ymin, d.Info.Ymax)
}

func (d *Device) GetInputXY(kX, kY int64) (int64, int64) {
	w := (kX - d.Info.Xmin) * d.ResW / (d.Info.Xmax - d.Info.Xmin)
	h := (kY - d.Info.Ymin) * d.ResH / (d.Info.Ymax - d.Info.Ymin)
	switch d.ScreenMode {
	case 1:
		return int64(math.Abs(float64(h))), int64(math.Abs(float64(w - d.ResW)))
	case 0:
		return w, h
	case 3:
		return int64(math.Abs(float64(h - d.ResH))), int64(math.Abs(float64(w)))
	default:
		log.Fatal("[ERROR]: device orientation not set !\n")
	}
	return 0, 0
}

func (d *Device) SetRes(w, h int64) {
	d.ResW = w
	d.ResH = h
}

/*
Set screen mode 1..3
1: landscape, USB port facing right
0: normal
3: landscape, USB port facing left
*/
func (d *Device) SetScreenMode(m int) {
	d.ScreenMode = m
}

func (d *Device) SetDeviceName(n string) {
	d.DeviceName = n
}
