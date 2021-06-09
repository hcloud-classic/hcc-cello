package handler

var grubdefault = "default Hcloud-Classic\n" +
	"label Hcloud-Classic\n" +
	"kernel CELLO_PXE_CONF_KERNEL\n"

var leaderoption = "append initrd=CELLO_PXE_CONF_LEADER_INITRAMFS" +
	" root=UUID=CELLO_PXE_CONF_LEADER_ROOT rw "
var computeoption = "append initrd=CELLO_PXE_CONF_COMPUTE_INITRAMFS" +
	" root=CELLO_PXE_CONF_COMPUTE_ROOT" +
	" nfsroot=CELLO_PXE_CONF_Leader_IP:/,rw "

var commonoption = " ip=dhcp session_id=CELLO_PXE_CONF_COMPUTE_SESSION_ID autonodeid=1 init=/sbin/init-krg biosdevname=0 net.ifnames=0 console=ttyS0,115200n8 rdshell"
var iscsioption = "netroot=iscsi:CELLO_PXE_CONF_ISCSI_SERVER_IP:tcp:3260:0:iqn.CELLO_PXE_CONF_ISCSI_TARGET_DOMAIN.target iscsi_initiator=iqn.CELLO_PXE_CONF_ISCSI_TARGET_DOMAIN.hcc withiscsi=1"

var defaultdir = "/root/boottp/HCC"
var pxecfgpath = "/pxelinux.cfg/default"
var configdir = "/etc/hcc/cello"
var iscsitarget = "target iqn.CELLO_PXE_CONF_ISCSI_TARGET_DOMAIN.target {" +
	" auth-group no-authentication" +
	" portal-group pg0" +
	" CELLO_PXE_CONF_ISCSI_LUN" +
	" }\n\n"
var iscsiportal = "portal-group pg0 { discovery-auth-group no-authentication listen 0.0.0.0 listen [::] }\n"
var iscsilun = "lun  LUN_NAME {" +
	"	path CELLO_PXE_CONF_ISCSI_ZVOLUME_PATH" +
	"	size CELLO_PXE_CONF_ISCSI_ZVOLUME_SIZE" +
	" }\n"

	// GlobalVolumesDB : Load volume info from DB
