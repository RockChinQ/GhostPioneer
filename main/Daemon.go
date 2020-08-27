package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
)

const TIME_OUT = 2 * 60

/**
GhostJ的守护进程
*/
func initDaemon() {
	fmt.Println("Daemon launch.")
	for {
		checkAliveThread()
		time.Sleep(TIME_OUT * time.Second)
	}
	wg.Done()
}
func checkAliveThread() {
	nowTime := time.Now().Unix()
	fmt.Println("check alive ", nowTime)
	if !Exists("alive") {
		WriteFile("alive", strconv.FormatInt(nowTime, 10))
		err := launchClient()
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	}
	clientTimeStr, err := ioutil.ReadFile("alive")
	clientTime, _ := strconv.ParseInt(string(clientTimeStr), 10, 64)
	if err != nil {
		return
	}
	//当前时间与客户端上次打出的时间差距超过timeout的秒数则启动客户端
	if (nowTime - clientTime) > int64(TIME_OUT) {
		go launchClient()
	}
}
