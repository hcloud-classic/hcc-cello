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
func PreparePxeSetting(ServerUUID string, OS string, networkIP string) (bool, interface{}) {

	err := logger.CreateDirIfNotExist(defaultdir + "/" + ServerUUID)
	if err != nil {
		return false, "Can't Create Directory at " + serverPxeDefaultDir
	}

	err = logger.CreateDirIfNotExist("/root/boottp/HCC/" + ServerUUID)
	if err != nil {
		logger.Logger.Fatal(err)
	}

	copyresult, test := copydefaultsetting(defaultdir+"/defaultLeader", defaultdir+"/"+ServerUUID+"/"+"Leader")
	if !copyresult {
		str := fmt.Sprintf("%v", test)
		return false, errors.New("Leader Pxe Setting Failed : " + str)
	}
	copyresult, test = copydefaultsetting(defaultdir+"/defaultCompute", defaultdir+"/"+ServerUUID+"/"+"Compute")
	if !copyresult {
		str := fmt.Sprintf("%v", test)
		return false, errors.New("Compute Pxe Setting Failed : " + str)

	}
	serverPxeDefaultDir := defaultdir + "/" + ServerUUID + "/"
	logger.Logger.Println("PxeDefaultDir=> ", serverPxeDefaultDir)
	if !rebuildPxeSetting(serverPxeDefaultDir, networkIP) {
		return false, errors.New("RebuildPxeSetting Failed")
	}

	return true, "Complete Pxe Setting"

}
func rebuildPxeSetting(pxeDir string, networkIP string) bool {
	leaderpxecfg := grubdefault + leaderoption + commonoption
	leaderpxecfg = strings.Replace(leaderpxecfg, "CELLO_PXE_CONF_LEADER_INITRAMFS", "Please Vaild ramfs", -1)
	leaderpxecfg = strings.Replace(leaderpxecfg, "CELLO_PXE_CONF_LEADER_ROOT", "Please Vaild root partition", -1)
	err := writeConfigFile(pxeDir, "Leader", leaderpxecfg)
	if err != nil {
		return false
	}
	computepxecfg := grubdefault + computeoption + commonoption
	computepxecfg = strings.Replace(computepxecfg, "CELLO_PXE_CONF_COMPUTE_INITRAMFS", "Please Vaild ramfs", -1)

	computepxecfg = strings.Replace(computepxecfg, "CELLO_PXE_CONF_COMPUTE_ROOT", "/dev/nfs", -1)
	computepxecfg = strings.Replace(computepxecfg, "CELLO_PXE_CONF_COMPUTE_NFS_IP", networkIP, -1)

	err = writeConfigFile(pxeDir, "Compute", computepxecfg)
	if err != nil {
		return false
	}
	return true
}

func writeConfigFile(pxeDir string, name string, contents string) error {
	confpath := pxeDir + name
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
	tmpstr := "cp -R " + src + " " + dst
	cmd := exec.Command("/bin/bash", "-c", tmpstr)
	result, err := cmd.CombinedOutput()
	if err != nil {
		return false, errors.New("Pxe Config can't write  " + string(result) + "=>  " + src + "  =>  " + dst)
	}
	return true, result
}
