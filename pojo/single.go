package pojo

import (
	"os"
	"os/exec"
	"xiazhicai.top/deskTool/util"
)

type Single struct {
	No int
	Id string
	Name string
	Size string
	ByteSize int
	DeskName string
	Status bool
}

// 通过shell获取多个磁盘，并返回磁盘切片
func GetList(plist Res) (list []*Single) {
	var single *Single
	for idx1,item := range plist.AllDisksAndPartitions {
		for _, partitions := range item.Partitions{
			single = newSingle(partitions.DeviceIdentifier, partitions.VolumeName, util.ChangeSize(partitions.Size),partitions.Size, idx1 + 1)
			list = append(list, single)
		}
	}
	return list
}

func newSingle(id,name,size string,byteSize,no int) *Single {
	if name == "" {
		name = "Untitled"
	}
	return &Single{
		No: no,
		Id:   id,
		Name: name,
		Size: size,
		Status: false,
	}
}

// 推出
func (s *Single) Umount(cmd *exec.Cmd) (err error) {
	cmd = exec.Command("sudo","umount","/dev/" + s.Id)
	_,err = cmd.CombinedOutput()
	if err != nil {
		return
	}
	s.Status = false
	deskName := "./Desktop/" + s.Id
	exists := util.IsExist(deskName)
	if exists {
		os.Remove(deskName)
	}
	return nil
}

// 挂载
func (s *Single) Mount(cmd *exec.Cmd) (err error) {
	// 判断文件夹是否存在
	deskName := "./Desktop/" + s.Id
	exists := util.IsExist(deskName)
	if !exists {
		err = os.MkdirAll(deskName, os.ModePerm)
		if err != nil{
			return
		}
	}
	cmd = exec.Command("sudo","mount_ntfs","-o","rw,nobrowse","/dev/" + s.Id,deskName)
	_,err = cmd.CombinedOutput()
	if err != nil {
		return
	}
	s.Status = true
	return
}
