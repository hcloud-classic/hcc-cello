package handler

import (
	"hcc/cello/lib/config"
	"hcc/cello/lib/logger"
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
func CreateVolume(FileSystem string, ServerUUID string, VolType string, Size int) (bool, interface{}) {
	volumePoolCheck()
	volcheck, err := QuotaCheck(ServerUUID)
	logger.Logger.Println("CreateVolume :QuotaCheck")
	if !volcheck {
		logger.Logger.Println("CreateVolume : check Faild", err)
		return volcheck, err
	}
	if VolType == "os" {
		logger.Logger.Println("Create ZFS(OS) : Faild")
		createcheck, err := clonezvol(FileSystem, ServerUUID, strings.ToUpper(VolType), strconv.Itoa(Size))
		if !createcheck {
			logger.Logger.Println("Create ZFS(OS) : Faild")
			return false, err
		}
	} else {
		createcheck, err := createzvol(FileSystem, ServerUUID, strings.ToUpper(VolType), strconv.Itoa(Size))
		if !createcheck {
			logger.Logger.Println("Create ZFS(DATA) : Faild")
			return false, err
		}
	}
	logger.Logger.Println("CreateVolume :After VolType")

	setquota(ServerUUID, Size)

	return true, err

}
func createzvol(FileSystem string, ServerUUID string, VolType string, Size string) (bool, interface{}) {

	volname := FileSystem + VolType + "-vol-" + ServerUUID
	zsysteminfo.ZfsName = zsysteminfo.PoolName + "/" + volname
	convertSize := Size + "G"
	volblocksize := "volblocksize=" + "4096"
	cmd := exec.Command("zfs", "create", "-V", convertSize, "-o", volblocksize, zsysteminfo.ZfsName)
	result, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	return true, result
}
func clonezvol(FileSystem string, ServerUUID string, VolType string, Size string) (bool, interface{}) {

	volname := FileSystem + VolType + "-vol-" + ServerUUID
	zsysteminfo.ZfsName = zsysteminfo.PoolName + "/" + volname
	cmd := exec.Command("zfs", "clone", config.VolumeConfig.ORIGINVOL, zsysteminfo.ZfsName)
	logger.Logger.Println("clonezvol : [", config.VolumeConfig.ORIGINVOL, "<      >", zsysteminfo.ZfsName, "]")
	result, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	return true, result
}

func createzfs(FileSystem string, ServerUUID string, VolType string) (bool, interface{}) {
	volname := FileSystem + VolType + "-vol-" + ServerUUID
	mountpath := "mountpoint=" + defaultdir + "/" + ServerUUID + "/" + FileSystem + "/" + VolType + "/"
	zsysteminfo.ZfsName = zsysteminfo.PoolName + "/" + volname
	cmd := exec.Command("zfs", "create", "-o", mountpath, zsysteminfo.ZfsName)
	logger.Logger.Println("createzfs : [", config.VolumeConfig.ORIGINVOL, "<      >", zsysteminfo.ZfsName, "]")

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
		logger.Logger.Println(result, err)
	}
	zsysteminfo.ZfsName = ""
	return true, err
}
func volumePoolCheck() {
	zsysteminfo.PoolName = config.VolumeConfig.VOLUMEPOOL
	// cmd := exec.Command("hostname")
	// result, err := cmd.CombinedOutput()
	// zsysteminfo.PoolName = strings.TrimSpace(string(result))
	if len(zsysteminfo.PoolName) == 0 {
		logger.Logger.Println("Please configuration volume management pool (/etc/hcc/cello/cello.conf)")
	}
}

//QuotaCheck : Zfs Available Quota check
func QuotaCheck(ServerUUID string) (bool, interface{}) {
	cmd := exec.Command("zfs", "get", "available", zsysteminfo.PoolName)
	result, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}
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
