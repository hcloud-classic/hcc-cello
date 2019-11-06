package handler

var grubdefault = "default kerrighed\n" +
	"label kerrighed\n" +
	"kernel vmlinuz-2.6.30-krg\n"

var leaderoption = "initrd=CELLO_PXE_CONF_LEADER_INITRAMFS" +
	" root=CELLO_PXE_CONF_LEADER_ROOT ro"
var computeoption = "initrd=CELLO_PXE_CONF_COMPUTE_INITRAMFS" +
	" root=CELLO_PXE_CONF_COMPUTE_ROOT" +
	" rootnfs=CELLO_PXE_CONF_COMPUTE_NFS_IP:/,rw"
var commonoption = " ip=dhcp session_id=1 autonodeid=1 init=/sbin/init-krg"

var defaultdir = "/root/boottp/HCC"
var pxecfgpath = "/pxelinux.cfg/default"
