package handler

import (
	"errors"
	"fmt"
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

//ReloadPoolObject : Reload Pool status
func ReloadPoolObject() error {
	act := "zpool list -H -o name"
	cmd := exec.Command("/bin/bash", "-c", act)
	result, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Can't exec zpool name ", err)
		return err
	}

	var poolobj formatter.Pool
	for _, args := range strings.Split(string(result), "\n") {
		if args != "" {
			act = "zpool list -H -o free " + args
			cmd = exec.Command("/bin/bash", "-c", act)
			result, err = cmd.CombinedOutput()
			if err != nil {
				fmt.Println("Can't exec zpool free : ", args, err)
				return err

			}
			poolobj.Free = string(result)

			act = "zpool list -H -o capacity " + args
			cmd = exec.Command("/bin/bash", "-c", act)
			result, err = cmd.CombinedOutput()
			if err != nil {
				fmt.Println("Can't exec zpool capacity : ", args, err)
				return err

			}
			poolobj.Capacity = string(result)

			act = "zpool list -H -o size " + args
			cmd = exec.Command("/bin/bash", "-c", act)
			result, err = cmd.CombinedOutput()
			if err != nil {
				fmt.Println("Can't exec zpool size ", args, err)
				return err

			}
			poolobj.Size = string(result)

			act = "zpool list -H -o health " + args
			cmd = exec.Command("/bin/bash", "-c", act)
			result, err = cmd.CombinedOutput()
			if err != nil {
				fmt.Println("Can't exec zpool health ", args, err)
				return err

			}
			poolobj.Health = string(result)

			poolobj.Name = args
			fmt.Println("pool => ", poolobj)
			formatter.PoolObjectMap.PutPool(poolobj)
		}
	}
	return nil

}

func addVolObejct(volume model.Volume) string {
	var lunNum string
	if formatter.VolObjectMap.Domain[volume.ServerUUID] == nil && strings.ToUpper(volume.UseType) != "DATA" {
		formatter.VolObjectMap.PutDomain(volume.ServerUUID)
	}
	if formatter.VolObjectMap.Domain[volume.ServerUUID] != nil {
		lunNum = formatter.VolObjectMap.SetIscsiLun(volume)
	}
	return lunNum
}

//To do
func removeVolObejct(volume model.Volume, lunNum int) {

	if formatter.VolObjectMap.Domain[volume.ServerUUID] != nil {
		formatter.VolObjectMap.RemoveIscsiLun(volume, lunNum)
	}
}

// CreateVolume : Creatte Volume
func CreateVolume(volume model.Volume) (bool, interface{}) {
	lunNum := addVolObejct(volume)
	if lunNum == "" {
		strerr := "Create ZFS(Lun Numbering) : Faild )"
		logger.Logger.Println(strerr)
		lunInt, _ := strconv.Atoi(lunNum)
		removeVolObejct(volume, lunInt)
		return false, errors.New("[Cello]  : " + strerr)
	}
	volume.LunNum, _ = strconv.Atoi(lunNum)
	if volume.UseType == "os" {
		createcheck, err := clonezvol(volume)
		if !createcheck {
			lunInt, _ := strconv.Atoi(lunNum)
			removeVolObejct(volume, lunInt)
			logger.Logger.Println("Create ZFS(OS) : Faild")
			return false, err
		}
	} else {
		createcheck, err := createzvol(volume)
		if !createcheck {
			lunInt, _ := strconv.Atoi(lunNum)
			removeVolObejct(volume, lunInt)
			logger.Logger.Println("Create ZFS(DATA) : Faild")
			return false, err
		}
	}

	return true, lunNum

}

func createzvol(volume model.Volume) (bool, interface{}) {

	volname := formatter.VolNameBuilder(volume) + "-" + strconv.Itoa(volume.LunNum)
	convertSize := strconv.Itoa(volume.Size) + "G"
	volblocksize := "volblocksize=" + "4096"
	cmd := exec.Command("zfs", "create", "-V", convertSize, "-o", volblocksize, volname)
	result, err := cmd.CombinedOutput()
	if err != nil {
		logger.Logger.Println(result, " ", err, " : ", cmd)
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

func AvailablePoolCheck() string {
	min := "0"
	var poolname string
	for _, args := range formatter.PoolObjectMap.PoolMap {
		tmpa, _ := strconv.Atoi(min)
		tmpstr := strings.Trim(strings.TrimSpace(args.Free), "G")
		tmpb, _ := strconv.Atoi(tmpstr)
		if tmpa <= tmpb {
			min = args.Free
			poolname = args.Name
		}
	}
	if poolname == "" {
		logger.Logger.Println("There Is No Available Pool")
	}

	return poolname
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
