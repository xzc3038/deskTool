package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"xiazhicai.top/deskTool/pojo"
	"xiazhicai.top/deskTool/util"
)

var info pojo.Info
var cmd *exec.Cmd
var list []*pojo.Single
var choose int

func main() {
	var str,ip []byte
	var err error

	// 1、检测当前环境
	mac := runtime.GOOS
	if mac != "darwin" {
		util.Exit("不支持当前系统!",2)
	}
	//getMode()

	// 2、获取获取当前Ip
	cmd = exec.Command("/bin/sh", "-c", `/sbin/ifconfig en0 | grep -E 'inet ' | awk '{print $2}'`)
	ip,err = cmd.CombinedOutput()
	if err == nil {
		info.Ip = string(ip)
	}

	// 3、获取挂载列表
	cmd = exec.Command("diskutil","list","-plist","external","physical")
	str,err = cmd.CombinedOutput()
	if err != nil {
		util.Exit("获取磁盘列表失败:" + err.Error(),2)
	}
	jsonStr := plist2json(str)

	plist := pojo.NewPlist(jsonStr)
	deskName := "./plist.plist"
	exists := util.IsExist(deskName)
	if exists {
		os.Remove(deskName)
	}
	list = pojo.GetList(plist)
	info.Desk = list
	mounted()
}

// plist转json
func plist2json(str []byte) []byte {
	fileObj, err := os.OpenFile("./plist.plist", os.O_WRONLY|os.O_CREATE, 0777)
	defer fileObj.Close()
	if err != nil {
		util.Exit("创建临时文件失败:" + err.Error(),2)
	}
	_, err = fileObj.Write(str)
	if err != nil {
		util.Exit("写入临时文件失败:" + err.Error(),2)
	}
	// 转化
	cmd = exec.Command("plutil","-convert","json","./plist.plist")
	_, err = cmd.CombinedOutput()
	if err != nil {
		util.Exit("转化临时文件失败:" + err.Error(),2)
	}
	pjson, err := ioutil.ReadFile("./plist.plist")
	if err != nil {
		util.Exit("读取临时文件失败:" + err.Error(),2)
	}
	return pjson
}

// 选择磁盘挂载
func mounted()  {
	var no string
	var number int
	var length = len(info.Desk)
	var buffer bytes.Buffer
	var strs,state string
	for {
		buffer.Reset()
		for _,value := range info.Desk {
			if value.Status {
				state = "已挂载"
			} else {
				state = "未挂载"
			}
			strs = fmt.Sprintf("                        %d    %s     %s       %s    %s \n",value.No,value.Name,value.Id,value.Size,state)
			buffer.WriteString(strs)
		}
		fmt.Println(`
			=======  外部磁盘列表  =======
			no   名称       ID           真实容量     本次状态
		`)
		fmt.Println(buffer.String())
		fileObj := bufio.NewReader(os.Stdin)
		fmt.Printf("请输入挂载磁盘的NO：")
		no,_ = fileObj.ReadString('\n')
		number,_ = strconv.Atoi(strings.Trim(no,"\n"))

		if number <= 0 || number > length{
			fmt.Println(fmt.Sprintf("请选择正确的NO,1~%d",length))
		} else {
			choose = number
			err := info.Desk[choose - 1].Umount(cmd)
			if err != nil {
				util.Exit("step1失败:" + err.Error(),0)
			} else {
				err = info.Desk[choose - 1].Mount(cmd)
				if err != nil {
					util.Exit("step2失败:" + err.Error(),0)
				} else {
					fmt.Println("挂载成功")
				}
			}
		}
	}
}

// 获取当前模式
func getMode() {
	var mode string
	for {
		fmt.Println(`
		====== 磁盘挂载(多分区磁盘默认挂载第一个分区，有需求可以自己改;将会挂载到当前路径) ======
		模式一、退出脚本后，自动取消挂载
		模式二、退出脚本后，不取消挂载，自行推出磁盘并删除桌面新增文件夹，名称如desk1
		模式三、挂载到别的地方去，先不写……
	`)
		fileObj := bufio.NewReader(os.Stdin)
		fmt.Printf("请输入模式：")
		mode,_ = fileObj.ReadString('\n')
		mode = strings.Trim(mode,"\n")
		if mode != "1" && mode != "2"{
			fmt.Println("请选择正确的模式 1 or 2")
		} else {
			info.Mode,_ = strconv.Atoi(mode)
			break
		}
	}
}