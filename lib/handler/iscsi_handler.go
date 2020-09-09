package handler

import (
	"hcc/cello/lib/formatter"
	"hcc/cello/lib/logger"
	"hcc/cello/model"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func iscsiServiceHandler() (bool, interface{}) {
	cmd := exec.Command("service", "ctld", "reload")
	result, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	return true, result
}

func WriteIscsiConfigObject(volume model.Volume) (bool, interface{}) {
	filename := "/etc/ctl.conf"
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
	iscsitarget := iscsitarget
	for i, args := range domain.Lun {
		singleLun := iscsilun
		lunname := formatter.VolNameBuilder(volume)
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

	iscsitarget = strings.Replace(iscsitarget, "CELLO_PXE_CONF_ISCSI_TARGET_DOMAIN", domain.TargetName, -1)
	iscsitarget = strings.Replace(iscsitarget, "CELLO_PXE_CONF_ISCSI_LUN", targetLunList, -1)

	return lunList + "\n" + iscsitarget
}

// // PrepareIscsiConfigObject : iscsi setting
// func PrepareIscsiConfigObject(volume model.Volume) (bool, interface{}) {
// 	leaderiscsilun := iscsilun
// 	leaderiscsilun = strings.Replace(leaderiscsilun, "CELLO_PXE_CONF_ISCSI_LUN_ORDER", "0", -1)
// 	leaderiscsilun = strings.Replace(leaderiscsilun, "CELLO_PXE_CONF_ISCSI_ZVOLUME_PATH", config.VolumeConfig.VOLUMEPOOL+"/"+FileSystem+strings.ToUpper(volumeType)+"-vol-"+ServerUUID, -1)
// 	// leaderiscsilun = strings.Replace(leaderiscsilun, "CELLO_PXE_CONF_ISCSI_ZVOLUME_SIZE", strconv.Itoa(volumeSize)+"G", -1)
// 	leaderiscsilun = strings.Replace(leaderiscsilun, "CELLO_PXE_CONF_ISCSI_ZVOLUME_SIZE", "100G", -1)
// 	leaderiscsitarget := iscsitarget
// 	leaderiscsitarget = strings.Replace(leaderiscsitarget, "CELLO_PXE_CONF_ISCSI_TARGET_DOMAIN", volume.ServerUUID, -1)
// 	leaderiscsitarget = strings.Replace(leaderiscsitarget, "CELLO_PXE_CONF_ISCSI_LUN", leaderiscsilun, -1)
// 	serverPxeDefaultDir := defaultdir + "/" + volume.ServerUUID + "/"
// 	err := iscsiwriteConfigFile(serverPxeDefaultDir, volume.ServerUUID+"-iscsi.conf", leaderiscsitarget)
// 	if err != nil {
// 		return false, "Not Write"
// 	}
// 	// err = iscsiappendFile("/etc/ctl.conf", "\ninclude \""+serverPxeDefaultDir+"\"")
// 	err = iscsiappendFile("/etc/ctl.conf", leaderiscsitarget)
// 	errstatus, result := iscsiServiceHandler()
// 	if !errstatus {
// 		return false, result
// 	}
// 	return true, "iscsi setting Complete"
// }

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
