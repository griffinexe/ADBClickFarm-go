package cli

import (
	"bufio"
	"context"
	"fmt"
	"goadb/util"
	"goadb/util/cmd/adb"
	"log"
	"os"
	"strings"
	"time"
)

var pre string = "\n> "
var title string = `
==========GoClickFarm================
Tip:
	> List device with "ls"
	> Connect ADB device over LAN "connect [ip:port]"
	> Disonnect ADB device over LAN "disconnect [ip:port]"
	> Select device with number and USB port facing "select 0 L"
	> Record input "record example"
	> Play record with file name and repeat times "play example 100"

`
var devList []string
var currentDev adb.Device = adb.Device{}
var jobs map[string]job = make(map[string]job)

type job struct {
	DevName string
	Fname   string
	Repeat  int
}

func cliLoop(s string) {
	if s == "ls" {
		devs, err := adb.ListDevice()
		if err != nil {
			log.Fatal(err)
		}
		devList = devs
		fmt.Printf("Listing %d device(s):\n", len(devs))
		for k, dev := range devs {
			fmt.Printf("%d\t%s\n", k, dev)
		}
		return
	}
	if strings.Contains(s, "connect") {
		cmds := strings.Split(s, " ")
		if len(cmds) < 2 {
			fmt.Println("invalid args")
			fmt.Println("connect [ip:port]")
			return
		}
		adb.ADBConnect(cmds[1])
	}
	if strings.Contains(s, "disconnect") {
		cmds := strings.Split(s, " ")
		if len(cmds) < 2 {
			fmt.Println("invalid args")
			fmt.Println("disconnect [ip:port]")
			return
		}
		adb.ADBDisconnect(cmds[1])
	}
	if strings.Contains(s, "select") {
		cmds := strings.Split(s, " ")
		if len(cmds) < 3 {
			fmt.Println("invalid args")
			fmt.Println("select [#dev] L/R/D")
			return
		}
		currentDev.DeviceName = devList[int(adb.Atoi(cmds[1]))]
		switch cmds[2] {
		case "L":
			currentDev.ScreenMode = 3
		case "R":
			currentDev.ScreenMode = 1
		case "D":
			currentDev.ScreenMode = 0
		default:
			fmt.Println("invalid args")
			fmt.Println("device facing must be one of USB port facing L(eft)/R(ight)/D(own)")
		}
		info, err := adb.GetScreenInfo(currentDev.DeviceName)
		if err != nil {
			log.Println("get device info error")
			log.Println(err)
			return
		}
		currentDev.Info = info
		res, err := adb.GetSysResulotion(currentDev.DeviceName)
		if err != nil {
			log.Println("get system resolution error")
			log.Println(err)
			return
		}
		currentDev.ResW = res[0]
		currentDev.ResH = res[1]
		currentDev.IsVM = false
		fmt.Println("device", currentDev.DeviceName, "selected")
		pre = "\n> "
		pre += currentDev.DeviceName + "> "
		return
	}
	if s == "info" {
		fmt.Printf(`
		Device name: 			%s
		Device is VM: 			%t
		System resulotion: 		%d*%d
		Touch screen info:
			Xmin: 	%d
			Xmax: 	%d
			Ymin: 	%d
			Ymax: 	%d
		`, currentDev.DeviceName, currentDev.IsVM, currentDev.ResW, currentDev.ResH, currentDev.Info.Xmin, currentDev.Info.Xmax, currentDev.Info.Ymin, currentDev.Info.Ymax)
	}
	if strings.Contains(s, "rec") {
		cmds := strings.Split(s, " ")
		if len(cmds) < 2 {
			fmt.Println("invalid args")
			fmt.Println("rec [filename]")
			return
		}
		ctx, cancel := context.WithCancel(context.Background())
		go adb.Record(&currentDev, cmds[1], ctx)
		fmt.Println("[REC] press 'q' to stop recording.")
		for {
			t := readIn()
			if t == "q" {
				cancel()
				time.Sleep(500 * time.Millisecond)
				adb.ADBSimTap(&currentDev)
				fmt.Printf("[REC] record end; saved as %s.rec", cmds[1])
				return
			} else {
				fmt.Println("[REC] press 'q' to stop recording.")
			}
		}
	}
	if strings.Contains(s, "play") {
		cmds := strings.Split(s, " ")
		if len(cmds) < 2 {
			fmt.Println("invalid args")
			fmt.Println("play [filename] [repeat]")
			return
		}
		jobId := util.RandomString(8)
		jobs[jobId] = job{
			DevName: currentDev.DeviceName,
			Fname:   cmds[1],
			Repeat:  int(adb.Atoi(cmds[2])),
		}
		go func(d *adb.Device) {
			for rep := 0; rep < int(adb.Atoi(cmds[2])); rep++ {
				fmt.Printf("[PLAY]: playing %s.rec on device %s repeat %d\n", cmds[1], d.DeviceName, rep)
				adb.Play(d, cmds[1])
			}
			delete(jobs, jobId)
			fmt.Printf("[PLAY END]: playing %s.rec on device %s\n", cmds[1], d.DeviceName)
		}(&currentDev)
	}
}

func cliRunFirst() {
	fmt.Println(title)
	fmt.Println("[INIT]: starting ADB")
	adb.StartADB()
	fmt.Println("[INIT]: ADB started")
}

func CLIStart() {
	adb.Mkdir()
	cliRunFirst()
	for {
		fmt.Print(pre)
		txt := readIn()
		if txt == "exit" {
			fmt.Println(">>> exit.")
			return
		}
		cliLoop(txt)
	}
}

func readIn() string {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	line = strings.ReplaceAll(strings.ReplaceAll(line, "\n", ""), "\r", "")
	return line
}
