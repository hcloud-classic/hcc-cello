package handler

import (
	"errors"
	"fmt"
	"hcc/cello/lib/logger"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var once sync.Once
var serverPxeDefaultDir string

//PreparePxeSetting : Prepare Pxe Setting, that pxelinux.cfg/default context update
// create pxe
func PreparePxeSetting(ServerUUID string, OS string, networkIP string) (bool, interface{}) {

	// if !createdir(ServerUUID) {
	// 	fmt.Println("Error")
	// }
	// copyresult, err := copydefaultsetting(defaultdir+"/defaultLeader", defaultdir+"/"+ServerUUID+"/"+"Leader")
	// if !copyresult {
	// 	fmt.Println(err)
	// 	str := fmt.Sprintf("%v", err)
	// 	return false, errors.New("Leader Pxe Setting Failed : " + str)
	// }
	// copyresult, err = copydefaultsetting(defaultdir+"/defaultCompute", defaultdir+"/"+ServerUUID+"/"+"Compute")
	// if !copyresult {
	// 	fmt.Println(err)
	// 	str := fmt.Sprintf("%v", err)
	// 	return false, errors.New("Leader Pxe Setting Failed : " + str)

	// }
	// serverPxeDefaultDir := defaultdir + "/" + ServerUUID + "/"
	// fmt.Println("serverPxeDefaultDir=> ", serverPxeDefaultDir)
	// if !rebuildPxeSetting(serverPxeDefaultDir, networkIP) {
	// 	return false, errors.New("RebuildPxeSetting Failed")
	// }
	// return true, "Complete Pxe Setting"
	return CreateDir(ServerUUID), "qwe"
}
func rebuildPxeSetting(pxeDir string, networkIP string) bool {
	leaderpxecfg := grubdefault + leaderoption + commonoption
	leaderpxecfg = strings.Replace(leaderpxecfg, "CELLO_PXE_CONF_LEADER_INITRAMFS", "initrd.img-2.6.30-krg", -1)
	leaderpxecfg = strings.Replace(leaderpxecfg, "CELLO_PXE_CONF_LEADER_ROOT", "/dev/mapper/krg-root", -1)
	fmt.Println("leaderpxecfg => ", leaderpxecfg)
	err := writeConfigFile(pxeDir, "Leader", leaderpxecfg)
	if err != nil {
		return false
	}
	computepxecfg := grubdefault + computeoption + commonoption
	computepxecfg = strings.Replace(computepxecfg, "CELLO_PXE_CONF_COMPUTE_INITRAMFS", "initrd.img-2.6.30-krg-nfs", -1)

	computepxecfg = strings.Replace(computepxecfg, "CELLO_PXE_CONF_COMPUTE_ROOT", "/dev/nfs", -1)
	computepxecfg = strings.Replace(computepxecfg, "CELLO_PXE_CONF_COMPUTE_NFS_IP", networkIP, -1)
	fmt.Println("computepxecfg => ", computepxecfg)

	err = writeConfigFile(pxeDir, "Compute", computepxecfg)
	if err != nil {
		return false
	}
	return true
}

func writeConfigFile(pxeDir string, name string, contents string) error {
	// confilepath := defaultdir + "/" + ServerUUID
	confpath := pxeDir + name
	fmt.Println("confpath => ", confpath)
	err := logger.CreateDirIfNotExist(confpath)
	if err != nil {
		return err
	}
	confpath = confpath + pxecfgpath
	err = writeFile(confpath, contents)
	if err != nil {
		return err
	}

	return nil
}
func writeFile(fileLocation string, input string) error {
	file, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	_, err = file.WriteString(input)
	if err != nil {
		return err
	}

	return nil
}
func copydefaultsetting(src string, dst string) (bool, interface{}) {

	cmd := exec.Command("cp", "-r", src, dst)
	result, err := cmd.CombinedOutput()
	if err != nil {
		return false, errors.New("Pxe Config can't write  " + string(result) + "=>  " + src + "  =>  " + dst)
	}
	return true, result
}

func CreateDir(ServerUUID string) bool {
	var err error
	returnValue := false
	once.Do(func() {
		// Create directory if not exist
		if _, err = os.Stat(defaultdir + "/" + ServerUUID); os.IsNotExist(err) {
			err = CreateDirIfNotExist(defaultdir + "/" + ServerUUID)
			if err != nil {
				panic(err)
			}
		}
		returnValue = true
	})

	return returnValue
}

// CreateDirIfNotExist : Make directory if not exist
func CreateDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}
