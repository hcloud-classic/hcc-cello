package handler

var grubdefault = "default Hcloud-Classic\n" +
	"label Hcloud-Classic\n" +
	"kernel [Please Write vaild kernel which you want]\n"

var leaderoption = "append initrd=CELLO_PXE_CONF_LEADER_INITRAMFS" +
	" root=CELLO_PXE_CONF_LEADER_ROOT ro"
var computeoption = "append initrd=CELLO_PXE_CONF_COMPUTE_INITRAMFS" +
	" root=CELLO_PXE_CONF_COMPUTE_ROOT" +
	" nfsroot=CELLO_PXE_CONF_COMPUTE_NFS_IP:/,rw"
var commonoption = " [Please Write vaild option which you want]"

var defaultdir = "/root/boottp/HCC"
var pxecfgpath = "/pxelinux.cfg/default"
