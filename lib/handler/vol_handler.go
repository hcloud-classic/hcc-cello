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

func findVolObejct(volume model.Volume) bool {
	if formatter.VolObjectMap.Domain[volume.ServerUUID] != nil && strings.ToUpper(volume.UseType) != "" && volume.UUID != "" {
		return true
	}
	return false
}

func PreLoad() {
	for _, args := range formatter.GlobalVolumesDB {
		lunNum := addVolObejct(args)
		if lunNum == "" {
			strerr := "Create ZFS(Lun Numbering) : Faild )"
			logger.Logger.Println(strerr)
			removeVolObejct(args)
		}
	}
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
func removeVolObejct(volume model.Volume) {
	if findVolObejct(volume) {
		_, lunNum := formatter.VolObjectMap.GetIscsiLun(volume)
		lunOrderInList, _ := strconv.Atoi(lunNum)
		formatter.VolObjectMap.RemoveIscsiLun(volume, lunOrderInList)
	}
}

// CreateVolume : Creatte Volume
func CreateVolume(volume model.Volume) (bool, interface{}) {
	lunNum := addVolObejct(volume)
	logger.Logger.Println("Codex :\n", volume, "\n Lunnum: ", lunNum)
	if lunNum == "" {
		strerr := "Create ZFS(Lun Numbering) : Faild )"
		logger.Logger.Println(strerr)
		// lunInt, _ := strconv.Atoi(lunNum)
		removeVolObejct(volume)
		return false, errors.New("[Cello]  : " + strerr)
	}
	volume.LunNum, _ = strconv.Atoi(lunNum)
	if volume.UseType == "os" {
		createcheck, err := clonezvol(volume)
		if !createcheck {
			// lunInt, _ := strconv.Atoi(lunNum)
			removeVolObejct(volume)
			logger.Logger.Println("Create ZFS(OS) : Faild")
			return false, err
		}
	} else {
		createcheck, err := createzvol(volume)
		if !createcheck {
			// lunInt, _ := strconv.Atoi(lunNum)
			removeVolObejct(volume)
			logger.Logger.Println("Create ZFS(DATA) : Faild")
			return false, err
		}
	}

	return true, lunNum

}

// DeleteVolumeObj : only remove volume obj
// TODO : Implement delete volume
func DeleteVolumeObj(volume model.Volume) (bool, interface{}) {
	if strings.ToLower(volume.UseType) == "os" {
		ejectDomain := formatter.VolObjectMap.GetDomain(volume.ServerUUID)
		checkResult := formatter.VolObjectMap.RemoveDomain(volume.ServerUUID)
		if !checkResult {
			logger.Logger.Println("Delete OS Volume Failed")
			return false, nil
		}
		fmt.Println("[Debug] : ", ejectDomain)
		return true, ejectDomain
	} else {
		lunStructure, _ := formatter.VolObjectMap.GetIscsiLun(volume)
		removeVolObejct(volume)
		return true, lunStructure
	}

}

//DeleteVolumeZFS : Delete Physical zVolume
func DeleteVolumeZFS(volName string) (bool, interface{}) {

	checkResult := destroyzvol(volName)
	if !checkResult {
		logger.Logger.Println("Delete Data Volume Failed")
		return checkResult, "Delete Data Volume Failed"

	}
	return checkResult, nil
}

func destroyzvol(volumeName string) bool {
	fmt.Println("destroyzvol : ", volumeName)
	cmd := exec.Command("zfs", "destroy", volumeName)
	result, err := cmd.CombinedOutput()
	if err != nil {
		logger.Logger.Println("destroyzvol : ", result, " ", err, " : ", cmd)
		return false
	}

	return true
}

func createzvol(volume model.Volume) (bool, interface{}) {

	volname := formatter.VolNameBuilder(volume) + "-" + strconv.Itoa(volume.LunNum)
	convertSize := strconv.Itoa(volume.Size) + "G"
	volblocksize := "volblocksize=" + "4096"
	cmd := exec.Command("zfs", "create", "-V", convertSize, "-o", volblocksize, volname)

	result, err := cmd.CombinedOutput()
	if err != nil {
		logger.Logger.Println("createzvol : ", result, " ", err, " : ", cmd)
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
	result, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	logger.Logger.Println("clonezvol : [", originVolName, " To ", volname, "]")

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

// UpdateVolume :
// TODO : Implement update volume
func UpdateVolume() {

}
