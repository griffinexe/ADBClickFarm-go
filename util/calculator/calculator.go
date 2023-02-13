package calculator

import (
	"goadb/util/cmd/adb"
	"strconv"
)

func Hex2Dec(h string) (int64, error) {
	i, err := strconv.ParseInt(h, 16, 32)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func GetWH(d *adb.Device, decW, decH int64) [2]int64 {
	w := (decW - d.Info.Xmin) * d.ResW / (d.Info.Xmax - d.Info.Xmin)
	h := (decH - d.Info.Ymin) * d.ResH / (d.Info.Ymax - d.Info.Ymin)
	return [2]int64{w, h}
}
