package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

func readStdin() {
	for {
		word, _, _ := stdin.ReadLine()
		spt := strings.Split(string(word), " ")
		switch spt[0] {
		case "~focus": //聚焦
			if len(spt) < 2 {
				fmt.Println("~focus <key>")
				continue
			}
			conn, ok := connMap[spt[1]]
			if !ok {
				fmt.Println("no such rescue client")
				continue
			}
			focused = bufio.NewWriter(conn)
			focusedName = spt[1]
			fmt.Println("focus " + spt[1])
			continue
		case "~ls":
			fmt.Println("list all rescue client.total:" + strconv.Itoa(len(connMap)))
			for key, conn := range connMap {
				fmt.Print(key + "\t" + conn.RemoteAddr().String())
				if focusedName == key {
					fmt.Println("\tfocused")
				} else {
					fmt.Println("\tdaemon")
				}
			}
			fmt.Println("done.")
			continue
		}
		//不是控制指令，发送到rescue client
		if focused != nil && focusedName != "null" {
			focused.Write([]byte(string(word) + "\n"))
			focused.Flush()
		} else {
			fmt.Println("no client focused")
		}
	}
}
