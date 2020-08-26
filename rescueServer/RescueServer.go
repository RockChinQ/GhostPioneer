package main

import (
	"bufio"
	"fmt"
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
	//获取连接
	fmt.Println("正在监听连接")
getConn:
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("接受连接失败", err.Error())
			continue getConn
		}
		go handleConn(conn)
	}
}
func checkFatalErr(err error, motion string) {
	if err != nil {
		fmt.Println("fatal error:\n", err.Error(), "\n"+motion+"失败")
		os.Exit(-1)
	}
}
func handleConn(conn *net.TCPConn) {
	fmt.Println("新连接正在处理")
	name := ""
	reader := bufio.NewReader(conn)
readMsg:
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			removeConn(conn)
			return
		}
		fmt.Println(msg)
		msgSpt := strings.Split(msg, " ")
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
