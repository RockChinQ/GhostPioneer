package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

func main() {
	//是否有在命令行指定target
	if len(os.Args) < 2 {
		fmt.Println("usage d.exe ip:port workerCount")
		os.Exit(-1)
	}
	worker, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("usage d.exe ip:port workerCount")
		os.Exit(-1)
	}
	ddos, err := New(os.Args[1], worker)
	if err != nil {
		fmt.Println("err:" + err.Error())
		os.Exit(1)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	ddos.Run()
	wg.Wait()
}
