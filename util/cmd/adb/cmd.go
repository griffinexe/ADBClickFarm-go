package adb

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	ADB         = "adb"
	ADBS        = "./bin/adb.exe"
	SHELL       = "shell"
	INPUT       = "input"
	GETEVENT    = "getevent"
	LIST_DEVICE = "devices"
	STOP        = "kill-server"
)

func Mkdir() {
	err := os.Mkdir("scripts", os.ModeAppend)
	if err != nil {
		log.Println(err)
	}
}

func StartADB() {
	c1 := exec.Command(ADB, STOP)
	c1.Run()
	c2 := exec.Command(ADB, LIST_DEVICE)
	c2.Run()
}

func ADBSimTap(d *Device) {
	c := exec.Command(ADB, "-s", d.DeviceName, INPUT, "tap", "0", "0")
	c.Run()
}

func ADBConnect(ip string) {
	c := exec.Command(ADB, "connect", ip)
	c.Run()
}

func ADBDisconnect(ip string) {
	c := exec.Command(ADB, "disconnect", ip)
	c.Run()
}

// 此功能复制自demo
func Play(dev *Device, fname string) error {
	f, err := os.Open("scripts/" + fname + ".rec")
	if err != nil {
		return err
	}
	defer f.Close()
	tInfo := new(TouchInfo)
	fsc := bufio.NewScanner(f)
	for fsc.Scan() {
		// log.Println(1)
		if fsc.Err() == io.EOF {
			log.Println("[PLAY]: EOF")
			return err
		}
		if fsc.Err() != nil {
			log.Println(err)
			return err
		}
		sX := strings.Fields(fsc.Text())
		tInfo.DelayBefore = Atoi(strings.Split(sX[1], ":")[1])
		tInfo.X1, tInfo.Y1 = dev.GetInputXY(Atoi(strings.Split(sX[2], ":")[1]), Atoi(strings.Split(sX[3], ":")[1]))
		tInfo.X2, tInfo.Y2 = dev.GetInputXY(Atoi(strings.Split(sX[4], ":")[1]), Atoi(strings.Split(sX[5], ":")[1]))
		tInfo.TouchDuration = Atoi(strings.Split(sX[6], ":")[1])
		doTouch(dev, tInfo)
	}
	return nil
}

func doTouch(d *Device, t *TouchInfo) error {
	// log.Printf("[DoTouch]: Wait %d ms then swipe from %d:%d to %d:%d in %d ms.\n", t.DelayBefore, t.X1, t.Y1, t.X2, t.Y2, t.TouchDuration)
	DelayMs(t.DelayBefore)
	err := SwipeTouch(d.DeviceName, t.X1, t.Y1, t.X2, t.Y2, t.TouchDuration)
	if err != nil {
		return err
	}
	// log.Println("[DoTouch]: command completed")
	return nil
}

func Record(dev *Device, fname string, ctx context.Context) error {
	// log.Println("[REC]")
	dev.PrintInfo()
	f, err := os.Create("scripts/" + fname + ".rec")
	if err != nil {
		return err
	}
	defer f.Close()
	c := exec.Command(ADB, "-s", dev.DeviceName, SHELL, GETEVENT)
	cmdStdout, _ := c.StdoutPipe()
	c.Start()

	var cX int64
	var cY int64
	pressing := false
	tInfo := new(TouchInfo)
	x1Set := false
	y1Set := false
	x2Set := false
	y2Set := false
	prePressSet := false
	pressDurSet := false

	cmdReader := bufio.NewReader(cmdStdout)
	tIdle := time.Now()
	tPress := time.Now()
	for {
		select {
		case <-ctx.Done():
			// log.Println("[REC] END")
			c.Process.Kill()
			return nil
		default:
			str, _, err := cmdReader.ReadLine()
			if err != nil {
				log.Fatalln("record error")
			}

			if strings.Contains(string(str), " 014a 00000001") {
				elapsed := time.Since(tIdle)
				tPress = time.Now()
				tIdle = time.Now()
				pressing = true
				tInfo.DelayBefore = elapsed.Milliseconds()
				prePressSet = true
				// println("> DOWN")
				continue
			}
			if strings.Contains(string(str), "0003 0035 ") {
				ssc := strings.Split(string(str), " ")
				i, err := strconv.ParseInt(ssc[len(ssc)-1], 16, 64)
				if err != nil {
					log.Println(err)
				}
				cX = i
				continue
			}
			if strings.Contains(string(str), "0003 0036 ") {
				ssc := strings.Split(string(str), " ")
				i, err := strconv.ParseInt(ssc[len(ssc)-1], 16, 64)
				if err != nil {
					log.Println(err)
				}
				cY = i
				continue
			}
			if strings.Contains(string(str), " 014a 00000000") {
				Telapsed := time.Since(tPress)
				tPress = time.Now()
				pressing = false
				tInfo.TouchDuration = Telapsed.Milliseconds()
				pressDurSet = true
			}
			if !x1Set && pressing && cX != 0 {
				tInfo.X1 = cX
				x1Set = true
			}
			if !y1Set && pressing && cY != 0 {
				tInfo.Y1 = cY
				y1Set = true
			}
			if !x2Set && !pressing && cX != 0 {
				tInfo.X2 = cX
				x2Set = true
			}
			if !y2Set && !pressing && cY != 0 {
				tInfo.Y2 = cY
				y2Set = true
			}
			if x1Set && x2Set && y1Set && y2Set && prePressSet && pressDurSet {
				// fmt.Printf("> pre:%d x1:%d y1:%d x2:%d y2:%d dur:%d\n", tInfo.DelayBefore, tInfo.X1, tInfo.Y1, tInfo.X2, tInfo.Y2, tInfo.TouchDuration)
				fmt.Fprintf(f, "> pre:%d x1:%d y1:%d x2:%d y2:%d dur:%d\n", tInfo.DelayBefore, tInfo.X1, tInfo.Y1, tInfo.X2, tInfo.Y2, tInfo.TouchDuration)
				x1Set = false
				x2Set = false
				y1Set = false
				y2Set = false
				prePressSet = false
				pressDurSet = false
				pressing = false
				cX, cY = 0, 0
			}
		}
	}
}

func GetScreenInfo(devname string) (*DeviceInfo, error) {
	devInfo := new(DeviceInfo)
	c := exec.Command(ADB, "-s", devname, SHELL, GETEVENT, "-p")
	out, err := c.Output()
	if err != nil {
		return nil, err
	}
	sout := strings.Split(string(out), "\n")
	for _, s := range sout {
		if strings.Contains(s, " 0035 ") {
			s2 := strings.Split(s, ", ")
			for _, s3 := range s2 {
				if strings.Contains(s3, "min ") {
					i, err := strconv.ParseInt(strings.Split(s3, " ")[1], 10, 64)
					if err != nil {
						return nil, err
					}
					devInfo.Xmin = i
				}
				if strings.Contains(s3, "max ") {
					i, err := strconv.ParseInt(strings.Split(s3, " ")[1], 10, 64)
					if err != nil {
						return nil, err
					}
					devInfo.Xmax = i
				}
			}
		}
		if strings.Contains(s, " 0036 ") {
			s2 := strings.Split(s, ", ")
			for _, s3 := range s2 {
				if strings.Contains(s3, "min ") {
					i, err := strconv.ParseInt(strings.Split(s3, " ")[1], 10, 64)
					if err != nil {
						return nil, err
					}
					devInfo.Ymin = i
				}
				if strings.Contains(s3, "max ") {
					i, err := strconv.ParseInt(strings.Split(s3, " ")[1], 10, 64)
					if err != nil {
						return nil, err
					}
					devInfo.Ymax = i
				}
			}
		}
	}
	return devInfo, nil
}

func GetSysResulotion(devname string) ([]int64, error) {
	res := []int64{}
	c := exec.Command(ADB, "-s", devname, SHELL, "wm", "size")
	out, err := c.Output()
	if err != nil {
		return nil, err
	}
	s := strings.Split(strings.Split(string(out), ": ")[1], "x")
	for _, v := range s {
		r, err := strconv.ParseInt(strings.ReplaceAll(strings.ReplaceAll(v, "\r", ""), "\n", ""), 10, 64)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil
}

func ListDevice() ([]string, error) {
	var devlist []string
	c := exec.Command(ADB, LIST_DEVICE)
	out, err := c.Output()
	if err != nil {
		return nil, err
	}
	sout := strings.Split(string(out), "\n")
	for _, s := range sout {
		if strings.Contains(s, "\tdevice") {
			devlist = append(devlist, strings.Split(s, "\t")[0])
		}
	}
	return devlist, nil
}

func SwipeTouch(devName string, X1, Y1, X2, Y2, Du int64) error {
	c := exec.Command(ADB, "-s", devName, SHELL, INPUT, "swipe", fmt.Sprint(X1), fmt.Sprint(Y1), fmt.Sprint(X2), fmt.Sprint(Y2), fmt.Sprint(Du))
	return c.Run()
}

func DelayMs(ms int64) {
	time.Sleep(time.Millisecond * time.Duration(ms))
}

func Atoi(a string) int64 {
	i, err := strconv.Atoi(a)
	if err != nil {
		log.Println(err)
		return 0
	}
	return int64(i)
}
