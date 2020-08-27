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

func main() {
	fmt.Println("launch rescue server")
	service := ":1032"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkFatalErr(err, "解析ip地址")
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkFatalErr(err, "启动端口监听器")
	//建立键盘输入接收

	//获取连接
	fmt.Println("正在监听连接")
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("接受连接失败", err.Error())
			continue
		}
		go handleConn(conn)
	}
}
func handleConn(conn *net.TCPConn) {
	fmt.Println("新连接正在处理")
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
			}
			continue readMsg
		}
	}
}
func removeConn(conn net.Conn) {
	for key, existConn := range connMap {
		if conn == existConn {
			delete(connMap, key)
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
