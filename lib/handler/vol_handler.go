package handler

import (
	"hcc/cello/lib/config"
	"hcc/cello/lib/formatter"
	"hcc/cello/lib/logger"
	"hcc/cello/model"
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

func addVolObejct(volume model.Volume) {
	if formatter.VolObjectMap.Domain[volume.ServerUUID] == nil {
		formatter.VolObjectMap.PutDomain(volume.ServerUUID)
	}
	if formatter.VolObjectMap.Domain[volume.ServerUUID] != nil {
		formatter.VolObjectMap.SetIscsiLun(volume, formatter.DevPathBuilder(volume))
	}
}

//To do
func removeVolObejct(volume model.Volume) {
	// if formatter.VolObjectMap.Domain[volume.ServerUUID] == nil {
	// 	formatter.VolObjectMap.PutDomain(volume.ServerUUID)
	// }
	// if formatter.VolObjectMap.Domain[volume.ServerUUID] != nil {
	// 	formatter.VolObjectMap.SetIscsiLun(volume, formatter.DevPathBuilder(volume))
	// }
}

// CreateVolume : Creatte Volume
func CreateVolume(volume model.Volume) (bool, interface{}) {
	volumePoolCheck()
	volume.Pool = config.VolumeConfig.VOLUMEPOOL
	volcheck, err := QuotaCheck(volume.ServerUUID)
	logger.Logger.Println("CreateVolume :QuotaCheck")
	if !volcheck {
		logger.Logger.Println("CreateVolume : check Faild", err)
		return volcheck, err
	}
	if volume.UseType == "os" {
		createcheck, err := clonezvol(volume)
		if !createcheck {
			logger.Logger.Println("Create ZFS(OS) : Faild")
			return false, err
		}
	} else {
		createcheck, err := createzvol(volume)
		if !createcheck {
			logger.Logger.Println("Create ZFS(DATA) : Faild")
			return false, err
		}
	}
	addVolObejct(volume)
	logger.Logger.Println("CreateVolume :After VolType")

	// setquota(ServerUUID, Size)

	return true, err

}

func createzvol(volume model.Volume) (bool, interface{}) {

	// volname := FileSystem + "-" + VolType + "-" + ServerUUID

	volname := formatter.VolNameBuilder(volume)
	convertSize := strconv.Itoa(volume.Size) + "G"
	volblocksize := "volblocksize=" + "4096"
	cmd := exec.Command("zfs", "create", "-V", convertSize, "-o", volblocksize, volname)
	result, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	return true, result
}

func clonezvol(volume model.Volume) (bool, interface{}) {
	volname := formatter.VolNameBuilder(volume)
	var originVolName = ""
	for _, args := range config.VolumeConfig.ORIGINVOL {
		if strings.Contains(strings.ToLower(args), strings.ToLower(volume.Filesystem)) {
			originVolName = args
			break
		}
	}
	cmd := exec.Command("zfs", "clone", originVolName, volname)
	logger.Logger.Println("clonezvol : [", originVolName, " To ", volname, "]")
	result, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	return true, result
}

func createzfs(volume model.Volume) (bool, interface{}) {
	volname := formatter.VolNameBuilder(volume)
	mountpath := "mountpoint=" + defaultdir + "/" + volname
	cmd := exec.Command("zfs", "create", "-o", mountpath, volname)

	result, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	return true, result
}

//Deprecate
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

//Deprecate
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
