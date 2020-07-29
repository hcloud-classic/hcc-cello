package handler

var grubdefault = "default Hcloud-Classic\n" +
	"label Hcloud-Classic\n" +
	"kernel CELLO_PXE_CONF_KERNEL\n"

var leaderoption = "append initrd=CELLO_PXE_CONF_LEADER_INITRAMFS" +
	" root=CELLO_PXE_CONF_LEADER_ROOT rw "
var computeoption = "append initrd=CELLO_PXE_CONF_COMPUTE_INITRAMFS" +
	" root=CELLO_PXE_CONF_COMPUTE_ROOT" +
	" nfsroot=CELLO_PXE_CONF_Leader_IP:/,rw "

var commonoption = " ip=dhcp session_id=CELLO_PXE_CONF_COMPUTE_SESSION_ID autonodeid=1 init=/sbin/init-krg biosdevname=0 net.ifnames=0 console=ttyS0,115200n8 rdshell"
var iscsioption = "netroot=iscsi:CELLO_PXE_CONF_ISCSI_SERVER_IP:tcp:3260:0:iqn.CELLO_PXE_CONF_ISCSI_TARGET_DOMAIN.target iscsi_initiator=iqn.CELLO_PXE_CONF_ISCSI_TARGET_DOMAIN.hcc withiscsi=1"

var defaultdir = "/root/boottp/HCC"
var pxecfgpath = "/pxelinux.cfg/default"

var iscsitarget = "target iqn.CELLO_PXE_CONF_ISCSI_TARGET_DOMAIN.target {\n" +
	"auth-group no-authentication\n" +
	"portal-group pg0\n" +
	"CELLO_PXE_CONF_ISCSI_LUN\n" +
	"}"
var iscsilun = "lun CELLO_PXE_CONF_ISCSI_LUN_ORDER {\n" +
	"	path /dev/zvol/CELLO_PXE_CONF_ISCSI_ZVOLUME_PATH\n" +
	"	size CELLO_PXE_CONF_ISCSI_ZVOLUME_SIZE\n" +
	"}"
