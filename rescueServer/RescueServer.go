package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

var connMap = make(map[string]net.Conn)

//已聚焦的连接
var focused *bufio.Writer
var focusedName string

/**
键盘输入读取
*/
var stdin *bufio.Reader

func main() {
	focusedName = "null"
	fmt.Println("launch rescue server")
	service := ":1032"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkFatalErr(err, "resolving ip addr")
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkFatalErr(err, "launching listener")
	//建立键盘输入接收
	stdin = bufio.NewReader(os.Stdin)
	go readStdin()
	//心跳数据
	go beatThread()
	//online daemon

	//获取连接
	fmt.Println("listening..")
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("cannot accept conn", err.Error())
			continue
		}
		go handleConn(conn)
	}
}
func handleConn(conn *net.TCPConn) {
	fmt.Println("handling new conn")
	ipStr := conn.RemoteAddr().String()
	defer func() {
		fmt.Println(" Disconnected : " + ipStr)
		conn.Close()
	}()

	name := ""
	reader := bufio.NewReader(conn)
readMsg:
	for {
		msg, _, err := reader.ReadLine()
		if err != nil || err == io.EOF {
			removeConn(conn)
			return
		}
		//fmt.Println(string(msg))
		msgSpt := strings.Split(string(msg), " ")
		switch msgSpt[0] {
		case "info":
			if len(msgSpt) >= 2 {
				name = msgSpt[1]
				fmt.Println("name:", name)
				connMap[name] = conn
			}
			continue readMsg
		case "~alives":
			continue readMsg
		}
		fmt.Println(string(msg))
		//fmt.Print(strings.ReplaceAll(string(msg),"#ln","\n"))
	}
}

func removeConn(conn net.Conn) {
	for key, existConn := range connMap {
		if conn == existConn {
			delete(connMap, key)
			if focusedName == key {
				focusedName = "null"
			}
			fmt.Println("kill " + key)
			return
		}
	}
}

func checkFatalErr(err error, motion string) {
	if err != nil {
		fmt.Println("fatal error:\n", err.Error(), "\n"+motion+"失败")
		os.Exit(-1)
	}
}
