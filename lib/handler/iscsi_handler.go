package handler

import (
	"fmt"
	"hcc/cello/lib/formatter"
	"hcc/cello/lib/logger"
	"hcc/cello/model"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func iscsiServiceHandler() (bool, interface{}) {
	cmd := exec.Command("service", "ctld", "status")
	result, err := cmd.CombinedOutput()
	if strings.Contains(string(result), "ctld is not running") {
		cmd = exec.Command("service", "ctld", "start")
		result, err = cmd.CombinedOutput()
		fmt.Println("start")

	} else {
		cmd := exec.Command("service", "ctld", "reload")
		result, err = cmd.CombinedOutput()
		fmt.Println("reload")
	}
	if err != nil {
		return false, err
	}
	return true, result
}

//WriteIscsiConfigObject : For iscsi Service config writer
func WriteIscsiConfigObject(volume model.Volume) (bool, interface{}) {
	filename := "/etc/ctl.conf"
	// volume.Pool = config.volumeig.VOLUMEPOOL
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		logger.Logger.Println("[WriteIscsiConfigObject]Does not Exist File ", filename)
	}
	defer func() {
		_ = file.Close()
	}()
	input := iscsiportal
	for _, args := range formatter.VolObjectMap.Domain {
		input += configBuilder(volume, args)
	}
	_, err = file.WriteString(input)
	if err != nil {
		logger.Logger.Println("[WriteIscsiConfigObject] Can't write file")
		// strerr := "create_volume action status=>iscsistatus " + fmt.Sprintln(err)
	}

	errstatus, result := iscsiServiceHandler()
	if !errstatus {
		return false, result
	}
	return true, "iscsi setting Complete"

}
func configBuilder(volume model.Volume, domain *formatter.Clusterdomain) string {
	lunList := ""
	targetLunList := ""
	targetDomain := iscsitarget
	for i, args := range domain.Lun {
		singleLun := iscsilun
		lunname := strings.Split(formatter.VolNameBuilder(volume), "/")[1]
		// if strings.Contains(strings.ToUpper(args.UseType), "DATA") {
		lunname = lunname + "-" + strconv.Itoa(i)
		// }

		singleLun = strings.Replace(singleLun, "LUN_NAME", lunname, -1)
		singleLun = strings.Replace(singleLun, "CELLO_PXE_CONF_ISCSI_LUN_ORDER", strconv.Itoa(i), -1)
		singleLun = strings.Replace(singleLun, "CELLO_PXE_CONF_ISCSI_ZVOLUME_PATH", formatter.DevPathBuilder(volume), -1)
		singleLun = strings.Replace(singleLun, "CELLO_PXE_CONF_ISCSI_ZVOLUME_SIZE", strconv.Itoa(args.Size)+"G", -1)

		lunList += singleLun + "\n "
		targetLunList += "lun " + strconv.Itoa(i) + " " + lunname + " "
	}

	targetDomain = strings.Replace(targetDomain, "CELLO_PXE_CONF_ISCSI_TARGET_DOMAIN", domain.TargetName, -1)
	targetDomain = strings.Replace(targetDomain, "CELLO_PXE_CONF_ISCSI_LUN", targetLunList, -1)

	return lunList + "\n" + targetDomain
}

func iscsiwriteConfigFile(iscsiconfdir string, name string, contents string) error {
	// confilepath := defaultdir + "/" + ServerUUID

	confpath := iscsiconfdir + name
	err := iscsiwriteFile(confpath, contents)
	if err != nil {
		return err
	}

	return nil
}
func iscsiwriteFile(fileLocation string, input string) error {
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
func iscsiappendFile(fileLocation string, input string) error {
	file, err := os.OpenFile(fileLocation, os.O_APPEND|os.O_WRONLY, 0644)
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
