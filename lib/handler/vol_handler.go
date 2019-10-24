package handler

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// ZSystem : Struct of ZSystem
type ZSystem struct {
	PoolName     string
	PoolCapacity string
	ZfsName      string //making zfs name
}

var zsysteminfo ZSystem

// CreateVolume : Creatte Volume
func CreateVolume(FileSystem string, ServerUUID string, OS string, Size int) (bool, interface{}) {
	// check hostname =>  because zpool follow hostname
	//
	hostCheck()
	volcheck, err := QuotaCheck(ServerUUID)
	if !volcheck {
		fmt.Println("CreateVolume : check Faild", err)
		return volcheck, err
	}
	createcheck, err := createzfs(FileSystem, ServerUUID, OS)
	if !createcheck {
		fmt.Println("Create ZFS : Faild")
		return createcheck, err
	}
	setquota(ServerUUID, Size)

	return true, err

}

func createzfs(FileSystem string, ServerUUID string, OS string) (bool, interface{}) {
	mountpath := "mountpoint=" + defaultdir + "/" + ServerUUID
	volname := FileSystem + OS + "-vol-" + ServerUUID
	zsysteminfo.ZfsName = zsysteminfo.PoolName + "/" + volname
	cmd := exec.Command("zfs", "create", "-o", "mountpoint", mountpath, zsysteminfo.ZfsName)
	result, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	return true, result
}

// zfs set quota=20G refquota=20G master/UUID-TEST
func setquota(ServerUUID string, Size int) (bool, interface{}) {
	qutoa := "quota=" + strconv.Itoa(Size) + "G"
	refquota := "refquota=" + strconv.Itoa(Size) + "G"

	cmd := exec.Command("zfs", "set", qutoa, refquota, zsysteminfo.ZfsName)
	result, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(result, err)
	}
	zsysteminfo.ZfsName = ""
	return true, err
}
func hostCheck() {
	cmd := exec.Command("hostname")
	result, err := cmd.CombinedOutput()
	zsysteminfo.PoolName = string(result)
	if err != nil {
		fmt.Println(result, err)
	}
}

//QuotaCheck : Zfs Available Quota check
func QuotaCheck(ServerUUID string) (bool, interface{}) {
	cmd := exec.Command("zfs", "get", "available", zsysteminfo.PoolName)
	result, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	// fmt.Println("=> ", strings.Fields(string(testqq)))
	tmpstr := strings.Fields(string(result))
	var posofvalue int
	for i, words := range tmpstr {
		if words == "VALUE" {
			posofvalue = (len(tmpstr) / 2) + i
		}
	}
	zsysteminfo.PoolCapacity = tmpstr[posofvalue]
	return true, tmpstr[posofvalue]
}

// DeleteVolume :
// TODO : Implement delete volume
func DeleteVolume() {

}

// UpdateVolume :
// TODO : Implement update volume
func UpdateVolume() {

}
