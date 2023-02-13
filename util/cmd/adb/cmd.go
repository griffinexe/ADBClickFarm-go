package adb

import (
	"bufio"
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

// func Tap(x, y int) error

// 此功能待验证
func Play(dev *Device, fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()
	tInfo := new(TouchInfo)
	// lastLine := make(chan string, 1)
	var currentX int64 = 0
	var currentY int64 = 0
	// isFirstXset := false
	// isFirstYset := false
	// isLastXset := false
	// isLastYset := false
	// pressBegin := false
	// pressEnd := false
	pressing := false
	isLastPosSet := false
	isFirstPosSet := false
	isPreDelaySet := false
	isDurationSet := false
	fsc := bufio.NewScanner(f)
	for fsc.Scan() {
		if fsc.Err() == io.EOF {
			log.Println("[PLAY]: EOF")
			return err
		}
		if fsc.Err() != nil {
			log.Println(err)
			return err
		}

		if strings.Contains(fsc.Text(), "T ") {
			tInfo.DelayBefore = atoi(strings.Split(fsc.Text(), " ")[1])
			isPreDelaySet = true
			continue
		}

		if strings.Contains(fsc.Text(), "@ DOWN") {
			pressing = true
			continue
			// 	pressBegin = true
			// } else {
			// 	pressBegin = false
		}
		if strings.Contains(fsc.Text(), "@ UP") {
			pressing = false
			continue
			// 	pressEnd = true
			// } else {
			// 	pressEnd = false
		}

		if strings.Contains(fsc.Text(), "X ") {
			currentX = atoi(strings.Split(fsc.Text(), " ")[1])
			continue
		}
		if strings.Contains(fsc.Text(), "Y ") {
			currentY = atoi(strings.Split(fsc.Text(), " ")[1])
			continue
		}

		if pressing && !isFirstPosSet {
			tInfo.X1 = currentX
			tInfo.Y1 = currentY
			isFirstPosSet = true
		}

		if !pressing && !isLastPosSet && isFirstPosSet {
			tInfo.X2 = currentX
			tInfo.Y2 = currentY
			isLastPosSet = true
		}
		// if strings.Contains(fsc.Text(), "# ") {
		// 	if pressBegin {
		// 		corr := strings.Split(strings.Split(fsc.Text(), " ")[1], ";")
		// 		tInfo.X1, tInfo.Y1 = dev.GetInputXY(atoi(corr[0]), atoi(corr[1]))
		// 	} else if pressEnd {
		// 		corr := strings.Split(strings.Split(<-lastLine, " ")[1], ";")
		// 		tInfo.X2, tInfo.Y2 = dev.GetInputXY(atoi(corr[0]), atoi(corr[1]))
		// 	} else {
		// 		continue
		// 	}
		// }
		if strings.Contains(fsc.Text(), "P ") {
			tInfo.TouchDuration = atoi(strings.Split(fsc.Text(), " ")[1])
			isDurationSet = true
		}

		if isDurationSet && isFirstPosSet && isLastPosSet && isPreDelaySet {
			tInfo.X1, tInfo.Y1 = dev.GetInputXY(tInfo.X1, tInfo.Y1)
			tInfo.X2, tInfo.Y2 = dev.GetInputXY(tInfo.X2, tInfo.Y2)
			doTouch(dev, tInfo)
			pressing = false
			isLastPosSet = false
			isFirstPosSet = false
			isPreDelaySet = false
			isDurationSet = false
		}

		// if len(lastLine) != 0 {
		// 	<-lastLine
		// }
		// lastLine <- fsc.Text()
	}
	return nil
}

func doTouch(d *Device, t *TouchInfo) error {
	log.Printf("[DoTouch]: Wait %d ms then swipe from %d:%d to %d:%d in %d ms.\n", t.DelayBefore, t.X1, t.Y1, t.X2, t.Y2, t.TouchDuration)
	DelayMs(t.DelayBefore)
	err := SwipeTouch(d.DeviceName, t.X1, t.Y1, t.X2, t.Y2, t.TouchDuration)
	if err != nil {
		return err
	}
	log.Println("[DoTouch]: command completed")
	return nil
}

func Record(dev *Device, fname string, stopCh chan bool) error {
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()
	c := exec.Command(ADB, "-s", dev.DeviceName, SHELL, GETEVENT)
	cmdStdout, _ := c.StdoutPipe()
	c.Start()
	cmdReader := bufio.NewReader(cmdStdout)
	lastLine := make(chan string, 1)
	T := time.Now()
	Tp := time.Now()
	for {
		select {
		case <-stopCh:
			c.Process.Kill()
		default:
			str, _, err := cmdReader.ReadLine()
			if err != nil {
				log.Fatalln("record error")
			}
			// if !(strings.Contains(string(str), " 014a 00000001") && strings.Contains(string(str), " 014a 00000000") && strings.Contains(string(str), "0003 0035 ") && strings.Contains(string(str), "0003 0036 ")) {
			// 	continue
			// }
			if strings.Contains(string(str), " 014a 00000001") {
				elapsed := time.Since(T)
				Tp = time.Now()
				T = time.Now()
				fmt.Println("T ", elapsed.Milliseconds())
				f.Write([]byte(fmt.Sprintln("T", elapsed.Milliseconds())))
				fmt.Print("@ DOWN\n")
				f.Write([]byte("@ DOWN\n"))
			}
			if strings.Contains(string(str), " 014a 00000000") {
				Telapsed := time.Since(Tp)
				Tp = time.Now()
				fmt.Print("@ UP\n")
				f.Write([]byte("@ UP\n"))
				fmt.Println("P", Telapsed.Milliseconds())
				f.Write([]byte(fmt.Sprintln("P", Telapsed.Milliseconds())))

			}
			if strings.Contains(string(str), "0003 0035 ") {
				ssc := strings.Split(string(str), " ")
				i, err := strconv.ParseInt(ssc[len(ssc)-1], 16, 64)
				if err != nil {
					log.Println(err)
				}
				fmt.Printf("X %d\n", i)
				f.Write([]byte(fmt.Sprintf("X %d\n", i)))

			}
			if strings.Contains(string(str), "0003 0036 ") {
				ssc := strings.Split(string(str), " ")
				i, err := strconv.ParseInt(ssc[len(ssc)-1], 16, 64)
				if err != nil {
					log.Println(err)
				}
				fmt.Printf("Y %d\n", i)
				f.Write([]byte(fmt.Sprintf("Y %d\n", i)))
			}
			if len(lastLine) != 0 {
				<-lastLine
			}
			lastLine <- string(str)
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

func atoi(a string) int64 {
	i, err := strconv.Atoi(a)
	if err != nil {
		log.Println(err)
		return 0
	}
	return int64(i)
}
