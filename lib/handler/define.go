package handler

var grubdefault = "default Hcloud-Classic\n" +
	"label Hcloud-Classic\n" +
	"kernel vmlinuz-2.6.30-hcc\n"

var leaderoption = "append initrd=CELLO_PXE_CONF_LEADER_INITRAMFS" +
	" root=CELLO_PXE_CONF_LEADER_ROOT ro"
var computeoption = "append initrd=CELLO_PXE_CONF_COMPUTE_INITRAMFS" +
	" root=CELLO_PXE_CONF_COMPUTE_ROOT" +
	" nfsroot=CELLO_PXE_CONF_COMPUTE_NFS_IP:/,rw"
var commonoption = " ip=dhcp session_id=1 autonodeid=1 init=/sbin/init-hcc"

var defaultdir = "/root/boottp/HCC"
var pxecfgpath = "/pxelinux.cfg/default"
