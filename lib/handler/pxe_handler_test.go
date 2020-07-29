package handler

import (
	"fmt"
	"strings"
)

// Pxetest: asd
func Pxetest(pxeDir string, networkIP string) bool {
	splitip := strings.Split(networkIP, ",")
	for i, words := range splitip {
		fmt.Println("[", i, "] ", words)
	}
	leaderpxecfg := grubdefault + leaderoption + commonoption
	leaderpxecfg = strings.Replace(leaderpxecfg, "CELLO_PXE_CONF_LEADER_INITRAMFS", "initrd.img-hcc", -1)
	leaderpxecfg = strings.Replace(leaderpxecfg, "CELLO_PXE_CONF_LEADER_ROOT", "/dev/sda1", -1)
	// logger.Logger.Println("leaderpxecfg => ", leaderpxecfg)
	// err := writeConfigFile(pxeDir, "Leader", leaderpxecfg)
	// if err != nil {
	// 	return false
	// }

	fmt.Println("leaderpxecfg [" + leaderpxecfg + "]")
	computepxecfg := grubdefault + computeoption + commonoption
	computepxecfg = strings.Replace(computepxecfg, "CELLO_PXE_CONF_COMPUTE_INITRAMFS", "initrd.img-hcc-nfs", -1)

	computepxecfg = strings.Replace(computepxecfg, "CELLO_PXE_CONF_COMPUTE_ROOT", "/dev/nfs", -1)
	computepxecfg = strings.Replace(computepxecfg, "CELLO_PXE_CONF_Leader_IP", networkIP, -1)
	// logger.Logger.Println("computepxecfg => ", computepxecfg)
	fmt.Println("computepxecfg [" + computepxecfg + "]")

	// err = writeConfigFile(pxeDir, "Compute", computepxecfg)
	// if err != nil {
	// 	return false
	// }
	return true
}
