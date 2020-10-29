package grpcsrv

import (
	"fmt"
	pb "hcc/cello/action/grpc/pb/rpccello"
	"hcc/cello/dao"
	"hcc/cello/lib/config"
	hccerr "hcc/cello/lib/errors"
	"hcc/cello/lib/formatter"
	handler "hcc/cello/lib/handler"
	"hcc/cello/lib/logger"
	"hcc/cello/model"
	"sort"
	"strconv"
	"strings"

	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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
		CreatedAt:  pbVolume.GetCreatedAt(),
	}
}

func reformatPBReqtoPBPool(contents *pb.ReqPoolHandler) *pb.Pool {
	pbPool := contents.GetPool()
	return &pb.Pool{
		UUID:          pbPool.GetUUID(),
		Size:          pbPool.GetSize(),
		Free:          pbPool.GetFree(),
		Capacity:      pbPool.GetCapacity(),
		Health:        pbPool.GetHealth(),
		Name:          pbPool.GetName(),
		AvailableSize: pbPool.GetAvailableSize(),
		Action:        pbPool.GetAction(),
		Used:          pbPool.GetUsed(),
	}
}

func reformatModelVolumetoPBVolume(volume *model.Volume) *pb.Volume {

	return &pb.Volume{
		UUID:       volume.UUID,
		Size:       strconv.Itoa(volume.Size),
		Filesystem: volume.Filesystem,
		ServerUUID: volume.ServerUUID,
		UseType:    volume.UseType,
		UserUUID:   volume.UserUUID,
		Network_IP: volume.NetworkIP,
		GatewayIp:  volume.GatewayIP,
		Pool:       volume.Pool,
		Lun:        int64(volume.LunNum),
		CreatedAt:  &timestamp.Timestamp{Seconds: volume.CreatedAt.Unix(), Nanos: int32(volume.CreatedAt.Nanosecond())},
	}
}

func reformatModelPooltoPBPool(pool *model.Pool) *pb.Pool {
	return &pb.Pool{
		UUID:          pool.UUID,
		Size:          pool.Size,
		Free:          pool.Free,
		Capacity:      pool.Capacity,
		Health:        pool.Health,
		Name:          pool.Name,
		AvailableSize: pool.AvailableSize,
		Action:        pool.Action,
	}
}

func reformatPBReqtoModelVolume(contents *pb.ReqVolumeHandler, volume *model.Volume) {
	pbVolume := contents.GetVolume()
	onlysize := strings.Split(pbVolume.GetSize(), "G")
	size, _ := strconv.Atoi(onlysize[0])
	volume.Size = size
	volume.UUID = pbVolume.GetUUID()
	volume.Filesystem = pbVolume.GetFilesystem()
	volume.ServerUUID = pbVolume.GetServerUUID()
	volume.UseType = pbVolume.GetUseType()
	volume.UserUUID = pbVolume.GetUserUUID()
	volume.NetworkIP = pbVolume.GetNetwork_IP()
	volume.GatewayIP = pbVolume.GetGatewayIp()
}

func reformatPBReqtoModelPool(contents *pb.ReqPoolHandler, pool *model.Pool) {
	pbPool := contents.GetPool()
	pool.UUID = pbPool.GetUUID()
	pool.Size = pbPool.GetSize()
	pool.Free = pbPool.GetFree()
	pool.Capacity = pbPool.GetCapacity()
	pool.Health = pbPool.GetHealth()
	pool.Name = pbPool.GetName()
	pool.AvailableSize = pbPool.GetAvailableSize()
	pool.Action = pbPool.GetAction()
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

func deleteAction(pbVolume *pb.Volume, volume *model.Volume) ([]formatter.Lun, *hccerr.HccErrorStack) {
	var retLunList []formatter.Lun
	errStack := hccerr.NewHccErrorStack()

	if pbVolume.UseType == "" || pbVolume.Filesystem == "" || pbVolume.ServerUUID == "" {
		if errStack != nil {
			errStack.Push(&hccerr.HccError{ErrCode: hccerr.CelloGrpcArgumentError})
			goto ERROR
		}
	}
	logger.Logger.Println("ActionHandle: Delete OS volume")

	switch strings.ToLower(volume.UseType) {
	case "os":

		deleteObjStatus, err := handler.DeleteVolumeObj(*volume)
		if !deleteObjStatus {
			errStack.Push(&hccerr.HccError{ErrCode: hccerr.CelloInternalCreateVolumeError})
			goto ERROR
		}
		lunInfo := (err).(*formatter.Clusterdomain)
		for _, args := range lunInfo.Lun {
			retLunList = append(retLunList, args)
		}

		iscsiStatus, err := handler.WriteIscsiConfigObject(*volume)
		if !iscsiStatus {
			errStack.Push(&hccerr.HccError{ErrCode: hccerr.CelloInternalWriteIscsiError})

			goto ERROR
		}
		logger.Logger.Println("WriteIscsiConfigObject : ", err)

		actionstatus, err := handler.DeletePxeSetting(volume)
		if !actionstatus {
			errStack.Push(&hccerr.HccError{ErrCode: hccerr.CelloInternalPreparePxeError, ErrText: err.(string)})

			goto ERROR
		}

	case "data":

		deleteObjStatus, err := handler.DeleteVolumeObj(*volume)
		if !deleteObjStatus {
			errStack.Push(&hccerr.HccError{ErrCode: hccerr.CelloInternalCreateVolumeError})
			goto ERROR
		}
		lunInfo := (err).(formatter.Lun)
		retLunList = append(retLunList, lunInfo)
		iscsiStatus, err := handler.WriteIscsiConfigObject(*volume)
		if !iscsiStatus {
			errStack.Push(&hccerr.HccError{ErrCode: hccerr.CelloInternalWriteIscsiError})

			goto ERROR
		}
		logger.Logger.Println("WriteIscsiConfigObject : ", err)

	default:
		errstr := "Use Type Invalid"
		errStack.Push(&hccerr.HccError{ErrText: errstr})
		goto ERROR
	}
	return retLunList, errStack.ConvertReportForm()

ERROR:
	errStack.Push(&hccerr.HccError{
		ErrText: "deleteAction(): Failed to delete volume",
	})

	return retLunList, errStack.ConvertReportForm()

}
func ReloadAllofVolInfo() error {

	celloParams := make(map[string]interface{})
	celloParams["row"] = 254
	celloParams["page"] = 1
	dbVol, err := dao.ReadVolumeAll(celloParams)
	if err != nil {
		fmt.Println("Error")
	}

	formatter.GlobalVolumesDB = dbVol.([]model.Volume)
	fmt.Println("ReloadAllofVolInfo", formatter.GlobalVolumesDB)
	sort.Slice(formatter.GlobalVolumesDB, func(i, j int) bool {
		return formatter.GlobalVolumesDB[i].LunNum < formatter.GlobalVolumesDB[j].LunNum
	})
	handler.PreLoad()
	fmt.Println("ReloadAllofVolInfo : \n", formatter.VolObjectMap.GetIscsiMap())
	return err
}

//VolumeHandler : Manipulate Volume Create
func VolumeHandler(contents *pb.ReqVolumeHandler) (*pb.Volume, *hccerr.HccErrorStack) {
	var err error
	var uuid string
	errStack := hccerr.NewHccErrorStack()
	var modelVolume model.Volume
	var tempModelVolume model.Volume

	reformatPBReqtoModelVolume(contents, &modelVolume)
	retPbVolume := reformatPBReqtoPBVolume(contents)
	logger.Logger.Println("Resolving: Volume Handle")
	if retPbVolume.ServerUUID == "" || retPbVolume.UserUUID == "" {
		errStack.Push(&hccerr.HccError{ErrCode: hccerr.CelloGrpcArgumentError, ErrText: "Invalid UUID : { Server: " + retPbVolume.ServerUUID + "\n User : " + retPbVolume.UserUUID + "}"})
		goto ERROR
	}
	err = ReloadAllofVolInfo()
	if err != nil {
		errStack.Push(&hccerr.HccError{ErrText: "Can't Preload "})
		logger.Logger.Println("Preload", errStack)
		goto ERROR
	}
	err = handler.ReloadPoolObject()
	if err != nil {
		errStack.Push(&hccerr.HccError{ErrText: "Can't Reload Object"})
		logger.Logger.Println("ReloadPoolObject", errStack)
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
		if tempErr.Len() > 0 {
			logger.Logger.Println("Error createAction: ", tempErr)
			errStack.AppendStack(tempErr)
			goto ERROR
		}
		retPbVolume.Lun = int64(modelVolume.LunNum)

		errcode, errstr := dao.CreateVolume(&modelVolume)
		if errstr != "" {
			logger.Logger.Println("Error DB : ", errstr)
			errStack.Push(&hccerr.HccError{ErrCode: errcode, ErrText: errstr})
			goto ERROR
		}
		errcode, tempVolume := dao.ReadVolume(&modelVolume)
		if tempVolume.UUID == "" {
			errStr := "Error DB : " + modelVolume.UUID + " is Not Exist"
			logger.Logger.Println()
			errStack.Push(&hccerr.HccError{ErrCode: errcode, ErrText: errStr})
			goto ERROR
		}
		retPbVolume = reformatModelVolumetoPBVolume(&tempVolume)
		logger.Logger.Println("[Create Volume] Success ")

	case "read_single":
		if retPbVolume.UUID == "" {
			errStack.Push(&hccerr.HccError{ErrCode: hccerr.CelloGrpcArgumentError, ErrText: "Invalid UUID : { Server: " + retPbVolume.ServerUUID})
			goto ERROR
		}

		errcode, tempVolume := dao.ReadVolume(&modelVolume)
		if tempVolume.UUID == "" {
			errStr := "Error DB : " + modelVolume.UUID + " is Not Exist"
			logger.Logger.Println()
			errStack.Push(&hccerr.HccError{ErrCode: errcode, ErrText: errStr})
			goto ERROR
		}
		retPbVolume = reformatModelVolumetoPBVolume(&tempVolume)
	case "read_list":
		if retPbVolume.ServerUUID == "" {
			errStack.Push(&hccerr.HccError{ErrCode: hccerr.CelloGrpcArgumentError, ErrText: "Invalid UUID : { Server: " + retPbVolume.ServerUUID})
			goto ERROR
		}

		errcode, tempVolume := dao.ReadVolume(&modelVolume)
		if tempVolume.UUID == "" {
			errStr := "Error DB : " + modelVolume.UUID + " is Not Exist"
			logger.Logger.Println()
			errStack.Push(&hccerr.HccError{ErrCode: errcode, ErrText: errStr})
			goto ERROR
		}
		retPbVolume = reformatModelVolumetoPBVolume(&tempVolume)
	case "update":

	case "delete":
		lunList, tempErr := deleteAction(retPbVolume, &modelVolume)
		if tempErr.Len() > 0 {
			logger.Logger.Println("Error deleteAction: ", tempErr)
			errStack.AppendStack(tempErr)
			goto ERROR
		}
		for i, args := range lunList {
			zfsDataSetVolName := strings.Split(args.Path, "/")
			deleteVolStatus, err := handler.DeleteVolumeZFS(args.Pool + "/" + zfsDataSetVolName[len(zfsDataSetVolName)-1])
			fmt.Println("Delete ", i, " : ", args)
			if !deleteVolStatus {
				errStack.Push(&hccerr.HccError{ErrCode: hccerr.CelloInternalCreateVolumeError, ErrText: err.(string)})
				goto ERROR
			}
			tempModelVolume.UUID = args.UUID
			errcode, errstr := dao.DeleteVolume(&tempModelVolume)
			if errstr != nil {
				logger.Logger.Println("Error DB : ", errstr)
				errStack.Push(&hccerr.HccError{ErrCode: errcode, ErrText: errstr.Error()})
				goto ERROR
			}
		}

		logger.Logger.Println("[Delete Volume] Success ")

	default:
		errstr := "Invalid Action : " + retPbVolume.Action
		errStack.Push(&hccerr.HccError{ErrCode: hccerr.CelloGrpcArgumentError, ErrText: errstr})
	}
	logger.Logger.Println("retPbVolume : ", retPbVolume)
	return retPbVolume, errStack.ConvertReportForm()

ERROR:
	errStack.Push(&hccerr.HccError{
		ErrCode: hccerr.CelloInternalVolumeHandleError,
		ErrText: "VolumeHandler(): Failed to handle volume",
	})

	return nil, errStack.ConvertReportForm()
}

func PoolHandler(contents *pb.ReqPoolHandler) (*pb.Pool, *hccerr.HccErrorStack) {
	var err error
	var uuid string
	errStack := hccerr.NewHccErrorStack()
	var modelPool model.Pool

	reformatPBReqtoModelPool(contents, &modelPool)
	retPbPool := reformatPBReqtoPBPool(contents)
	logger.Logger.Println("Resolving: Pool Handle")
	// if retPbPool.Name == "" {
	// 	errStack.Push(&hccerr.HccError{ErrCode: hccerr.CelloGrpcArgumentError, ErrText: "Invalid Pool name : " + retPbPool.Name + "}"})
	// 	goto ERROR
	// }
	err = ReloadAllofVolInfo()
	if err != nil {
		errStack.Push(&hccerr.HccError{ErrText: "Can't Preload "})
		logger.Logger.Println("Preload", errStack)
		goto ERROR
	}
	err = handler.ReloadPoolObject()
	if err != nil {
		errStack.Push(&hccerr.HccError{ErrText: "Can't Reload Object"})
		logger.Logger.Println("ReloadPoolObject", errStack)
		goto ERROR
	}

	switch strings.ToLower(retPbPool.Action) {
	case "create":

		out, err := gouuid.NewV4()
		if err != nil {
			logger.Logger.Println("[VolumeHandler]Can't Create Volume UUID : ", err)
			goto ERROR
		}
		uuid = out.String()

		modelPool.UUID = uuid
		retPbPool.UUID = modelPool.UUID

	case "read":

		for _, args := range formatter.PoolObjectMap.PoolMap {
			fmt.Println("formatter.PoolObjectMap.PoolMap\n\n", args)

			if args.Name == config.VolumeConfig.VOLUMEPOOL {
				retPbPool.AvailableSize = args.AvailableSize
				retPbPool.Capacity = args.Capacity
				retPbPool.Free = args.Free
				retPbPool.Size = args.Size
				retPbPool.Health = args.Health
				retPbPool.Name = args.Name
				retPbPool.Used = args.Used
			}
		}
	case "update":

	case "delete":

	default:
		errstr := "Invalid Action : " + retPbPool.Action
		errStack.Push(&hccerr.HccError{ErrCode: hccerr.CelloGrpcArgumentError, ErrText: errstr})
	}
	logger.Logger.Println("retPbPool : ", retPbPool)
	return retPbPool, errStack.ConvertReportForm()

ERROR:
	errStack.Push(&hccerr.HccError{
		ErrCode: hccerr.CelloInternalVolumeHandleError,
		ErrText: "PoolHandler(): Failed to handle Pool",
	})

	return nil, errStack.ConvertReportForm()

}
