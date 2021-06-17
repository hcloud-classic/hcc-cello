package handler

import (
	"errors"
	"fmt"
	"hcc/cello/dao"
	"hcc/cello/lib/config"
	"hcc/cello/lib/formatter"
	"hcc/cello/lib/logger"
	"hcc/cello/model"
	"math"
	"os/exec"
	"sort"
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
			convIntPoolFree, _ := strconv.ParseFloat(strings.Trim(strings.TrimSpace(poolobj.Free), "GTBM"), 64)
			if strings.Contains(poolobj.Free, "T") {
				convIntPoolFree *= 1024
			}
			poolobj.Free = strconv.Itoa(int(math.Floor(convIntPoolFree)))

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

			var intSize int
			for _, domainArgs := range formatter.VolObjectMap.Domain {
				for _, singleLun := range domainArgs.Lun {
					intSize += singleLun.Size
				}
			}
			convIntPoolSize, _ := strconv.ParseFloat(strings.Trim(strings.TrimSpace(poolobj.Size), "GTBM"), 64)

			if strings.Contains(poolobj.Size, "T") {
				convIntPoolSize *= 1024
			}
			poolobj.Size = strconv.Itoa(int(math.Floor(convIntPoolSize)))
			poolobj.AvailableSize = strconv.Itoa(int(math.Floor(convIntPoolSize - float64(intSize))))

			poolobj.Used = strconv.Itoa(intSize)
			fmt.Println("convIntPoolSize : ", convIntPoolSize, " intSize: ", poolobj.Used)

			act = "zpool list -H -o health " + args
			cmd = exec.Command("/bin/bash", "-c", act)
			result, err = cmd.CombinedOutput()
			if err != nil {
				fmt.Println("Can't exec zpool health ", args, err)
				return err

			}
			poolobj.Health = string(result)

			poolobj.Name = args
			// act = "zfs get -H -o value available master/volpool-1"
			fmt.Println("pool => ", poolobj)
			formatter.PoolObjectMap.PutPool(poolobj)
		}
	}
	return nil

}
func ReloadAllOfVolInfo() error {
	var recvErr error
	// var recvStr string
	var ErrStr string
	var ErrCode uint64
	celloParams := make(map[string]interface{})
	celloParams["row"] = 254
	celloParams["page"] = 1
	dbVol, ErrCode, ErrStr := dao.ReadVolumeAll(celloParams)
	if recvErr != nil {
		logger.Logger.Println("ReloadAllOfVolInfo(): Failed to read volumes", ErrCode, ErrStr)
		return recvErr
	}
	formatter.GlobalVolumesDB = dbVol.([]model.Volume)
	fmt.Println("ReloadAllOfVolInfo", formatter.GlobalVolumesDB)
	sort.Slice(formatter.GlobalVolumesDB, func(i, j int) bool {
		return formatter.GlobalVolumesDB[i].LunNum < formatter.GlobalVolumesDB[j].LunNum
	})
	PreLoad()
	fmt.Println("ReloadAllOfVolInfo : \n", formatter.VolObjectMap.GetIscsiMap())

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
			strerr := " ZFS(Lun Numbering) : Faild )"
			logger.Logger.Println(strerr)
			removeVolObejct(args)
		}
	}
}

func addVolObejct(volume model.Volume) string {
	var lunNum string
	if formatter.VolObjectMap.Domain[volume.ServerUUID] == nil {
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
		strerr := " ZFS(Lun Numbering) : Faild )"
		logger.Logger.Println(strerr)
		// lunInt, _ := strconv.Atoi(lunNum)
		removeVolObejct(volume)
		return false, errors.New("[Cello]  : " + strerr)
	}
	volume.LunNum, _ = strconv.Atoi(lunNum)
	if volume.UseType == "os" {
		check, err := clonezvol(volume)
		if !check {
			// lunInt, _ := strconv.Atoi(lunNum)
			removeVolObejct(volume)
			logger.Logger.Println(" ZFS(OS) : Faild")
			return false, err
		}
	} else {
		check, err := createzvol(volume)
		if !check {
			// lunInt, _ := strconv.Atoi(lunNum)
			removeVolObejct(volume)
			logger.Logger.Println(" ZFS(DATA) : Faild")
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
		// fmt.Println("[Debug] : ", ejectDomain)
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
	refquota := "refreservation=none"
	compression := "compression=on"
	reservation := "reservation=" + convertSize

	cmd := exec.Command("zfs", "create", "-V", convertSize, "-o", volblocksize, "-o", refquota, "-o", reservation, "-o", compression, volname)
	result, err := cmd.CombinedOutput()
	if err != nil {
		logger.Logger.Println("createzvol Failed : ", string(result), " ", err.Error(), " : ", cmd)
		return false, err
	}

	return true, result
}

func clonezvol(volume model.Volume) (bool, interface{}) {
	volname := formatter.VolNameBuilder(volume)
	var originVolName = ""
	for _, args := range config.VolumeConfig.ORIGINVOL {
		if strings.Contains(strings.ToLower(args), strings.ToLower(volume.Filesystem)) {
			originVolName = volume.Pool + "/" + args
			break
		}
	}
	logger.Logger.Println("clonezvol : [", originVolName, " To ", volname, "]")

	cmd := exec.Command("zfs", "clone", originVolName, volname)
	result, err := cmd.CombinedOutput()
	if err != nil {
		logger.Logger.Println("Clone Failed ", string(result), " err: ", err.Error())
		return false, err
	}

	return true, result
}

// func createzfs(volume model.Volume) (bool, interface{}) {
// 	volname := formatter.VolNameBuilder(volume)
// 	mountpath := "mountpoint=" + defaultdir + "/" + volname
// 	cmd := exec.Command("zfs", "", "-o", mountpath, volname)

// 	result, err := cmd.CombinedOutput()
// 	if err != nil {
// 		return false, err
// 	}
// 	return true, result
// }

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
		fmt.Println(args.Name)
	}
	if poolname == "" {
		logger.Logger.Println("There Is No Available Pool")
	}
	fmt.Println("############## ", poolname)
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

func round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

// convertRealPoolSize : Calculate Real Pool size, 2GB => 1.862645149230957 ==1.85
func convertRealPoolSize(size int) float64 {
	convFloatSize := float64(size)
	floatResult := float64(convFloatSize * 1000 * 1000 * 1000 / 1024 / 1024 / 1024)
	return round(floatResult, 0.05)

}

// convertRealVolumeSize : For  Volume quota
func convertRealVolumeSize(size int) float64 {
	convFloatSize := float64(size)
	floatResult := float64(convFloatSize * 1000 * 1000 / 1024 / 1024)
	return round(floatResult, 0.05)

}
