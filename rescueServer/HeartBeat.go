package main

import (
	"bufio"
	"time"
)

func beatThread() {
	for {
		time.Sleep(time.Second * 30)
		for _, conn := range connMap {
			go beatConn(bufio.NewWriter(conn))
		}
	}
}
func beatConn(writer *bufio.Writer) {
	writer.Write([]byte("~alive" + "\n"))
	writer.Flush()
}
