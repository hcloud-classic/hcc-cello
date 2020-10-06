package grpcsrv

import (
	pb "hcc/cello/action/grpc/pb/rpccello"
	"hcc/cello/dao"
	hccerr "hcc/cello/lib/errors"
	handler "hcc/cello/lib/handler"
	"hcc/cello/lib/logger"
	"hcc/cello/model"
	"strconv"
	"strings"

	gouuid "github.com/nu7hatch/gouuid"
)

func reformatPBReqtoPBVolume(contents *pb.ReqVolumeHandler) *pb.Volume {
	pbVolume := contents.GetVolume()
	return &pb.Volume{
		UUID:       pbVolume.GetUUID(),
		Size:       pbVolume.GetSize(),
		Filesystem: pbVolume.GetFilesystem(),
		ServerUUID: pbVolume.GetServerUUID(),
		UseType:    pbVolume.GetUseType(),
		UserUUID:   pbVolume.GetUserUUID(),
		Network_IP: pbVolume.GetNetwork_IP(),
		GatewayIp:  pbVolume.GetGatewayIp(),
		Pool:       pbVolume.GetPool(),
		Lun:        int64(pbVolume.GetLun()),
		Action:     pbVolume.GetAction(),
	}
}

func reformatPBReqtoModelVolume(contents *pb.ReqVolumeHandler, volume *model.Volume) {
	pbVolume := contents.GetVolume()
	onlysize := strings.Split(pbVolume.GetSize(), "G")
	size, _ := strconv.Atoi(onlysize[0])
	volume.Size = size
	volume.Filesystem = pbVolume.GetFilesystem()
	volume.ServerUUID = pbVolume.GetServerUUID()
	volume.UseType = pbVolume.GetUseType()
	volume.UserUUID = pbVolume.GetUserUUID()
	volume.NetworkIP = pbVolume.GetNetwork_IP()
	volume.GatewayIP = pbVolume.GetGatewayIp()
}

func createAction(pbVolume *pb.Volume, volume *model.Volume) *hccerr.HccErrorStack {
	errStack := hccerr.NewHccErrorStack()

	if pbVolume.UseType == "" || pbVolume.Filesystem == "" || pbVolume.Network_IP == "" || pbVolume.GatewayIp == "" {
		if errStack != nil {
			errStack.Push(&hccerr.HccError{ErrCode: hccerr.CelloGrpcArgumentError})
			goto ERROR
		}
	}

	if pbVolume.Pool == "" {
		volume.Pool = handler.AvailablePoolCheck()
		if len(volume.Pool) == 0 {
			errStack.Push(&hccerr.HccError{ErrCode: hccerr.CelloInternalStoragePoolError})
			goto ERROR
		}
	} else {
		volume.Pool = pbVolume.Pool
	}

	if volume.UseType == "os" {
		logger.Logger.Println("ActionHandle: Creating OS volume")

		createstatus, err := handler.CreateVolume(*volume)
		if !createstatus {
			errStack.Push(&hccerr.HccError{ErrCode: hccerr.CelloInternalCreateVolumeError})
			goto ERROR
		}
		lunNum := err.(string)
		volume.LunNum, _ = strconv.Atoi(lunNum)

		actionstatus, err := handler.PreparePxeSetting(volume.ServerUUID, volume.UseType, volume.NetworkIP, volume.GatewayIP)
		if !actionstatus {
			errStack.Push(&hccerr.HccError{ErrCode: hccerr.CelloInternalPreparePxeError})
			goto ERROR
		}

		iscsistatus, err := handler.WriteIscsiConfigObject(*volume)
		if !iscsistatus {
			errStack.Push(&hccerr.HccError{ErrCode: hccerr.CelloInternalWriteIscsiError})
			goto ERROR
		}
		logger.Logger.Println("[Action Result]  WriteIscsiConfigObject : ", actionstatus, " , CreateVolume : ", createstatus, "PrepareIscsiSetting : ", iscsistatus)

	}

	if pbVolume.UseType == "data" {
		logger.Logger.Println("ActionHandle: Creating data volume")

		createstatus, err := handler.CreateVolume(*volume)
		lunNum := err.(string)
		volume.LunNum, _ = strconv.Atoi(lunNum)

		if !createstatus {
			errStack.Push(&hccerr.HccError{ErrCode: hccerr.CelloInternalCreateVolumeError})
			goto ERROR
		}
		iscsistatus, err := handler.WriteIscsiConfigObject(*volume)
		if !iscsistatus {
			errStack.Push(&hccerr.HccError{ErrCode: hccerr.CelloInternalWriteIscsiError})
			goto ERROR
		}

		logger.Logger.Println("[Action Result]  createstatus  :", createstatus, " iscsistatus : ", iscsistatus)
	}
	pbVolume.Lun = int64(volume.LunNum)
	pbVolume.Pool = volume.Pool

	return errStack.ConvertReportForm()

ERROR:
	errStack.Push(&hccerr.HccError{
		ErrText: "createAction(): Failed to create volume",
	})

	return errStack.ConvertReportForm()

}

//VolumeHandler : Manipulate Volume Create
func VolumeHandler(contents *pb.ReqVolumeHandler) (*pb.Volume, *hccerr.HccErrorStack) {
	var err error
	var uuid string
	errStack := hccerr.NewHccErrorStack()
	var modelVolume model.Volume
	reformatPBReqtoModelVolume(contents, &modelVolume)
	retPbVolume := reformatPBReqtoPBVolume(contents)
	logger.Logger.Println("Resolving: create_volume")
	if retPbVolume.ServerUUID == "" || retPbVolume.UserUUID == "" {
		errStack.Push(&hccerr.HccError{ErrCode: hccerr.CelloGrpcArgumentError, ErrText: "Invalid UUID : { Server: " + retPbVolume.ServerUUID + "\n User : " + retPbVolume.UserUUID + "}"})
		goto ERROR
	}

	err = handler.ReloadPoolObject()
	if err != nil {
		errStack.Push(&hccerr.HccError{ErrText: "Can't Reload Object"})
		goto ERROR
	}

	switch strings.ToLower(retPbVolume.Action) {
	case "create":

		out, err := gouuid.NewV4()
		if err != nil {
			logger.Logger.Println("[VolumeHandler]Can't Create Volume UUID : ", err)
			goto ERROR
		}
		uuid = out.String()

		modelVolume.UUID = uuid
		retPbVolume.UUID = modelVolume.UUID

		tempErr := createAction(retPbVolume, &modelVolume)
		if tempErr != nil {
			logger.Logger.Println(tempErr)
			errStack.AppendStack(tempErr)
			goto ERROR
		}
		retPbVolume.Lun = int64(modelVolume.LunNum)

		errcode, errstr := dao.CreateVolume(&modelVolume)
		if errstr != "" {
			errStack.Push(&hccerr.HccError{ErrCode: errcode, ErrText: errstr})
			goto ERROR
		}

		logger.Logger.Println("[Create Volume] Success ")

	case "read":
	case "update":
	case "delete":
	default:
		errstr := "Invalid Action : " + retPbVolume.Action
		errStack.Push(&hccerr.HccError{ErrCode: hccerr.CelloGrpcArgumentError, ErrText: errstr})
	}
	return retPbVolume, errStack.ConvertReportForm()

ERROR:
	errStack.Push(&hccerr.HccError{
		ErrCode: hccerr.CelloInternalVolumeHandleError,
		ErrText: "VolumeHandler(): Failed to handle volume",
	})

	return nil, errStack.ConvertReportForm()
}
