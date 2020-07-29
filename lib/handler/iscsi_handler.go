package handler

import (
	"hcc/cello/lib/config"
	"os"
	"os/exec"
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

// PrepareIscsiSetting : iscsi setting
func PrepareIscsiSetting(ServerUUID string, FileSystem string, volumeType string, volumeSize int) (bool, interface{}) {
	leaderiscsilun := iscsilun
	leaderiscsilun = strings.Replace(leaderiscsilun, "CELLO_PXE_CONF_ISCSI_LUN_ORDER", "0", -1)
	leaderiscsilun = strings.Replace(leaderiscsilun, "CELLO_PXE_CONF_ISCSI_ZVOLUME_PATH", config.VolumeConfig.VOLUMEPOOL+"/"+FileSystem+strings.ToUpper(volumeType)+"-vol-"+ServerUUID, -1)
	// leaderiscsilun = strings.Replace(leaderiscsilun, "CELLO_PXE_CONF_ISCSI_ZVOLUME_SIZE", strconv.Itoa(volumeSize)+"G", -1)
	leaderiscsilun = strings.Replace(leaderiscsilun, "CELLO_PXE_CONF_ISCSI_ZVOLUME_SIZE", "100G", -1)
	leaderiscsitarget := iscsitarget
	leaderiscsitarget = strings.Replace(leaderiscsitarget, "CELLO_PXE_CONF_ISCSI_TARGET_DOMAIN", ServerUUID, -1)
	leaderiscsitarget = strings.Replace(leaderiscsitarget, "CELLO_PXE_CONF_ISCSI_LUN", leaderiscsilun, -1)
	serverPxeDefaultDir := defaultdir + "/" + ServerUUID + "/"
	err := iscsiwriteConfigFile(serverPxeDefaultDir, ServerUUID+"-iscsi.conf", leaderiscsitarget)
	if err != nil {
		return false, "Not Write"
	}
	// err = iscsiappendFile("/etc/ctl.conf", "\ninclude \""+serverPxeDefaultDir+"\"")
	err = iscsiappendFile("/etc/ctl.conf", leaderiscsitarget)
	errstatius, result := iscsiServiceHandler()
	if !errstatius {
		return false, result
	}
	return true, "iscsi setting Complete"
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
