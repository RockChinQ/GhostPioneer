package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"time"
)

//离线安装器
//检查工作目录是否有ghostjc.zip
//有就解压到GhostJHome
func main() {
	mkGhostDir("D:\\ProgramData\\")
	Unzip("ghostjc.zip", "D:\\ProgramData\\")
	nowTime := time.Now().Unix()
	WriteFile("D:\\ProgramData\\Ghost\\"+strconv.FormatInt(nowTime, 10)+".txt", strconv.FormatInt(nowTime, 10))

	stdin:=bufio.NewReader(os.Stdin)
	fmt.Print("[提示]请设置此主机的名称(如果已有，将会覆盖，如果您不确定请不要设置，然后询问其他人员):")
	words,_,_:=stdin.ReadLine()
	err := WriteFile("D:\\ProgramData\\Ghost\\preferName.txt", string(words))
	if err != nil {
		fmt.Println("[错误]"+err.Error())
	}
	_ = launchEXE("D:\\ProgramData\\Ghost\\", "D:\\ProgramData\\Ghost\\gl.exe")
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

//在指定的工作目录启动exe
func launchEXE(workDir string, filename string) error {
	c := exec.Command(filename)
	c.Dir = workDir
	err := c.Start()
	return err
}
func mkGhostDir(dir string) {
	_ = os.MkdirAll(dir, 0777)
}
func Unzip(filename string, savePath string) {
	fmt.Println("Unzip " + filename)
	r, err := zip.OpenReader(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, k := range r.Reader.File {
		if k.FileInfo().IsDir() {
			err := os.MkdirAll(savePath+k.Name, 0644)
			if err != nil {
				fmt.Println(err)
			}
			continue
		}
		r, err := k.Open()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("unzip: ", savePath+k.Name)
		defer r.Close()
		NewFile, err := os.Create(savePath + k.Name)
		if err != nil {
			fmt.Println(err)
			continue
		}
		io.Copy(NewFile, r)
		NewFile.Close()
	}
}
