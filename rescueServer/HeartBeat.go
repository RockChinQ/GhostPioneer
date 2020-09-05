package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

func beatThread() {
	for {
		time.Sleep(time.Second * 30)
		for _, conn := range connMap {
			go beatConn(bufio.NewWriter(conn))
		}
		go checkDiscClient()
	}
}
func checkDiscClient() {
	//fmt.Println("check clients")
	//检查没有启动的客户端
	clientsStr, _ := ReadFile("onlineClients.txt")
	//fmt.Println("clientStr:"+clientsStr)
	clients := strings.Split(clientsStr, " ")
	//遍历rescue，找出存在rescue但是没有client的机器
	for key, conn := range connMap {
		//如果rescue叫rtest则跳过
		if key == "rtest" {
			continue
		}
		//clients列表里是否存在与rescue同名的连接
		//不存在则启动
		fmt.Println("check rescue:" + key)
		if !IsContains(clients, key) {
			write := bufio.NewWriter(conn)
			fmt.Println("launching disc client:" + key)
			write.Write([]byte("launch\n"))
			write.Flush()
		}
	}
}
func beatConn(writer *bufio.Writer) {
	writer.Write([]byte("~alive" + "\n"))
	writer.Flush()
}

func ReadFile(filename string) (string, error) {

	f, err := ioutil.ReadFile(filename)

	if err != nil {
		return "", err
	}

	return string(f), nil

}
func IsContains(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}
