package handler

import (
	"errors"
	"fmt"
	"hcc/cello/lib/logger"
	"os"
	"os/exec"
	"strings"
)

//PreparePxeSetting : Prepare Pxe Setting, that pxelinux.cfg/default context update
// create pxe
func PreparePxeSetting(ServerUUID string, OS string, networkIP string, gateway string) (bool, interface{}) {

	err := logger.CreateDirIfNotExist(defaultdir + "/" + ServerUUID)
	if err != nil {
		return false, "Can't Create Directory at " + defaultdir + "/" + ServerUUID
	}

	// if _, err := os.Stat("/root/boottp/HCC/" + ServerUUID); os.IsNotExist(err) {
	// 	err = os.MkdirAll("/root/boottp/HCC/"+ServerUUID, 0755)
	// 	if err != nil {
	// 		return false, err
	// 	}
	// }

	copyresult, test := copydefaultsetting(defaultdir+"/defaultLeader", defaultdir+"/"+ServerUUID+"/"+"Leader")
	if !copyresult {
		// logger.Logger.Println(test)
		str := fmt.Sprintf("%v", test)
		return false, errors.New("Leader Pxe Setting Failed : " + str)
	}
	copyresult, test = copydefaultsetting(defaultdir+"/defaultCompute", defaultdir+"/"+ServerUUID+"/"+"Compute")
	if !copyresult {
		// logger.Logger.Println(test)
		str := fmt.Sprintf("%v", test)
		return false, errors.New("Compute Pxe Setting Failed : " + str)

	}
	serverPxeDefaultDir := defaultdir + "/" + ServerUUID + "/"
	logger.Logger.Println("PxeDefaultDir=> ", serverPxeDefaultDir)
	if !rebuildPxeSetting(ServerUUID, serverPxeDefaultDir, networkIP, gateway) {
		return false, errors.New("RebuildPxeSetting Failed")
	}

	return true, "Complete Pxe Setting"

}
func rebuildPxeSetting(ServerUUID string, pxeDir string, networkIP string, gateway string) bool {
	leaderpxecfg := grubdefault + leaderoption + iscsioption + commonoption
	leaderpxecfg = strings.Replace(leaderpxecfg, "CELLO_PXE_CONF_KERNEL", "vmlinuz-hcc", -1)
	leaderpxecfg = strings.Replace(leaderpxecfg, "CELLO_PXE_CONF_LEADER_INITRAMFS", "initrd.img-hcc", -1)
	leaderpxecfg = strings.Replace(leaderpxecfg, "CELLO_PXE_CONF_LEADER_ROOT", "/dev/sda1", -1)
	splitip := strings.Split(networkIP, ".")
	leaderpxecfg = strings.Replace(leaderpxecfg, "CELLO_PXE_CONF_COMPUTE_SESSION_ID", splitip[2], -1)
	leaderpxecfg = strings.Replace(leaderpxecfg, "CELLO_PXE_CONF_ISCSI_SERVER_IP", gateway, -1)
	leaderpxecfg = strings.Replace(leaderpxecfg, "CELLO_PXE_CONF_ISCSI_TARGET_DOMAIN", ServerUUID, -1)

	// logger.Logger.Println("leaderpxecfg => ", leaderpxecfg)
	err := writeConfigFile(pxeDir, "Leader", leaderpxecfg)
	if err != nil {
		return false
	}
	computepxecfg := grubdefault + computeoption + commonoption
	computepxecfg = strings.Replace(computepxecfg, "CELLO_PXE_CONF_COMPUTE_INITRAMFS", "initrd.img-hcc-nfs", -1)

	computepxecfg = strings.Replace(computepxecfg, "CELLO_PXE_CONF_COMPUTE_ROOT", "/dev/nfs", -1)
	computepxecfg = strings.Replace(computepxecfg, "CELLO_PXE_CONF_Leader_IP", networkIP, -1)
	// logger.Logger.Println("computepxecfg => ", computepxecfg)

	err = writeConfigFile(pxeDir, "Compute", computepxecfg)
	if err != nil {
		return false
	}
	return true
}

func writeConfigFile(pxeDir string, name string, contents string) error {
	// confilepath := defaultdir + "/" + ServerUUID
	confpath := pxeDir + name
	// logger.Logger.Println("confpath => ", confpath)
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
