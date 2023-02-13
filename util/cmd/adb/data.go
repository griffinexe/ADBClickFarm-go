package adb

type DeviceInfo struct {
	Xmin int64
	Xmax int64
	Ymin int64
	Ymax int64
}

type TouchInfo struct {
	DelayBefore   int64
	X1            int64
	Y1            int64
	X2            int64
	Y2            int64
	TouchDuration int64
}
