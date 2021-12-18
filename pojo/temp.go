package pojo

import (
	"encoding/json"
)

type Res struct {
	AllDisksAndPartitions []Item `json:"AllDisksAndPartitions"`
}

type Item struct {
	Partitions []Partitions `json:"Partitions"`
	Content string `json:Content`
	Size int `json:Size`
	DeviceIdentifier string `json:DeviceIdentifier`
}

type Partitions struct {
	MountPoint string `json:MountPoint`
	VolumeName string `json:VolumeName`
	Content string `json:Content`
	Size int `json:Size`
	DeviceIdentifier string `json:DeviceIdentifier`
}

func NewPlist(str []byte) Res {
	res := &Res{}
	json.Unmarshal(str,res)
	return *res
}
