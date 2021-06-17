package grpcsrv

import (
	"fmt"
	"hcc/cello/dao"
	"hcc/cello/lib/config"
	"hcc/cello/lib/formatter"
	handler "hcc/cello/lib/handler"
	"hcc/cello/lib/logger"
	"hcc/cello/model"
	"strconv"
	"strings"

	"github.com/golang/protobuf/ptypes"
	gouuid "github.com/nu7hatch/gouuid"
	"innogrid.com/hcloud-classic/hcc_errors"
	"innogrid.com/hcloud-classic/pb"
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
	tmpCreatedAt, _ := ptypes.TimestampProto(volume.CreatedAt.Local())
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
		// CreatedAt:  ptypes.time &ptypes.Timestamp{Seconds: volume.CreatedAt.Unix(), Nanos: int32(volume.CreatedAt.Nanosecond())},
		CreatedAt: tmpCreatedAt,
	}
}

func reformatModelVolumetoPBVolumeList(volume *[]model.Volume) []*pb.Volume {
	var retVol []*pb.Volume
	for _, args := range *volume {
		tmpCreatedAt, _ := ptypes.TimestampProto(args.CreatedAt.Local())

		tempPbVol := pb.Volume{
			UUID:       args.UUID,
			Size:       strconv.Itoa(args.Size),
			Filesystem: args.Filesystem,
			ServerUUID: args.ServerUUID,
			UseType:    args.UseType,
			UserUUID:   args.UserUUID,
			Network_IP: args.NetworkIP,
			GatewayIp:  args.GatewayIP,
			Pool:       args.Pool,
			Lun:        int64(args.LunNum),
			// CreatedAt:  &ptypes.Timestamp{Seconds: args.CreatedAt.Unix(), Nanos: int32(args.CreatedAt.Nanosecond())},
			CreatedAt: tmpCreatedAt,
		}
		retVol = append(retVol, &tempPbVol)
	}
	return retVol
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

func createAction(pbVolume *pb.Volume, volume *model.Volume) (uint64, string) {
	var ErrStr string
	var ErrCode uint64
	if pbVolume.UseType == "" {
		ErrStr += "UseType Null "
		ErrCode = hcc_errors.CelloGrpcArgumentError
		logger.Logger.Println(ErrStr)
		goto ERROR
	}

	if pbVolume.Pool == "" {
		ErrStr += "Pool Not Null "
		volume.Pool = handler.AvailablePoolCheck()
		if len(volume.Pool) == 0 {
			ErrStr += "Pool None of capacity "
			ErrCode = hcc_errors.CelloInternalStoragePoolError
			logger.Logger.Println(ErrStr)
			goto ERROR
		}
	} else {
		volume.Pool = pbVolume.Pool
	}

	if volume.UseType == "os" {
		logger.Logger.Println("ActionHandle: Creating OS volume")
		if pbVolume.Filesystem == "" || pbVolume.Network_IP == "" || pbVolume.GatewayIp == "" {
			ErrStr += "Failed : Filesystem||Network_IP||GatewayIp Check! "
			ErrCode = hcc_errors.CelloInternalStoragePoolError
			logger.Logger.Println(ErrStr)
			goto ERROR
		}
		createstatus, err := handler.CreateVolume(*volume)
		if !createstatus {
			ErrStr += "Failed : Create volume " + err.(error).Error()
			ErrCode = hcc_errors.CelloInternalAction
			logger.Logger.Println(ErrStr)
			goto ERROR
		}
		lunNum := err.(string)
		volume.LunNum, _ = strconv.Atoi(lunNum)

		actionstatus, err := handler.PreparePxeSetting(volume.ServerUUID, volume.UseType, volume.NetworkIP, volume.GatewayIP)
		if !actionstatus {
			ErrStr += "Failed : Prepare PxeSetting " + err.(error).Error()
			ErrCode = hcc_errors.CelloInternalPrepareDeploy
			logger.Logger.Println(ErrStr)
			goto ERROR
		}

		iscsistatus, err := handler.WriteIscsiConfigObject(*volume)
		if !iscsistatus {
			ErrStr += "Failed : Write Iscsi Config Object " + err.(error).Error()
			ErrCode = hcc_errors.CelloInternalPrepareDeploy
			logger.Logger.Println(ErrStr)
			goto ERROR
		}
		logger.Logger.Println("[Action Result]  WriteIscsiConfigObject : ", actionstatus, " , CreateVolume : ", createstatus, "PrepareIscsiSetting : ", iscsistatus)

	}

	if pbVolume.UseType == "data" {
		logger.Logger.Println("ActionHandle: Creating data volume")
		if volume.Filesystem == "" {
			for _, args := range formatter.GlobalVolumesDB {
				if args.ServerUUID == volume.ServerUUID {
					volume.Filesystem = args.Filesystem
					break
				}
			}
		}
		createstatus, err := handler.CreateVolume(*volume)
		lunNum := err.(string)
		volume.LunNum, _ = strconv.Atoi(lunNum)

		if !createstatus {
			ErrStr += "Failed : Create volume " + err.(error).Error()
			ErrCode = hcc_errors.CelloInternalAction
			logger.Logger.Println(ErrStr)
			goto ERROR
		}
		iscsistatus, err := handler.WriteIscsiConfigObject(*volume)
		if !iscsistatus {
			ErrStr += "Failed : Write Iscsi Config Object " + err.(error).Error()
			ErrCode = hcc_errors.CelloInternalPrepareDeploy
			logger.Logger.Println(ErrStr)
			goto ERROR
		}

		logger.Logger.Println("[Action Result]  createstatus  :", createstatus, " iscsistatus : ", iscsistatus)
	}
	pbVolume.Lun = int64(volume.LunNum)
	pbVolume.Pool = volume.Pool

	return 0, ""

ERROR:
	ErrStr += "createAction(): Failed to create volume " + ErrStr
	return ErrCode, ErrStr

}

func deleteAction(pbVolume *pb.Volume, volume *model.Volume) ([]formatter.Lun, uint64, string) {
	var retLunList []formatter.Lun
	var ErrStr string
	var ErrCode uint64
	if len(pbVolume.UseType) == 0 || len(pbVolume.ServerUUID) == 0 {
		ErrStr += "UseType||ServerUUID Null "
		ErrCode = hcc_errors.CelloGrpcArgumentError
		logger.Logger.Println(ErrStr)
		goto ERROR
	}
	logger.Logger.Println("ActionHandle: Delete  volume")

	switch strings.ToLower(volume.UseType) {
	case "os":

		deleteObjStatus, err := handler.DeleteVolumeObj(*volume)
		if !deleteObjStatus {
			ErrStr += "Failed : Delete Volume Object"
			ErrCode = hcc_errors.CelloInternalAction
			logger.Logger.Println(ErrStr)
			goto ERROR
		}
		lunInfo := (err).(*formatter.Clusterdomain)
		for _, args := range lunInfo.Lun {
			retLunList = append(retLunList, args)
		}

		iscsiStatus, err := handler.WriteIscsiConfigObject(*volume)
		if !iscsiStatus {
			ErrStr += "Failed : Write Iscsi Config Object " + err.(error).Error()
			ErrCode = hcc_errors.CelloInternalPrepareDeploy
			logger.Logger.Println(ErrStr)
			goto ERROR
		}
		logger.Logger.Println("WriteIscsiConfigObject : ", err)

		actionstatus, err := handler.DeletePxeSetting(volume)
		if !actionstatus {
			ErrStr += "Failed : Write Iscsi Config Object " + err.(error).Error()
			ErrCode = hcc_errors.CelloInternalPrepareDeploy
			logger.Logger.Println(ErrStr)
			goto ERROR
		}

	case "data":
		if volume.Filesystem == "" {
			for _, args := range formatter.GlobalVolumesDB {
				if args.ServerUUID == volume.ServerUUID {
					volume.Filesystem = args.Filesystem
					break
				}
			}
		}
		deleteObjStatus, err := handler.DeleteVolumeObj(*volume)
		if !deleteObjStatus {
			ErrStr += "Failed : Delete Volume Object"
			ErrCode = hcc_errors.CelloInternalAction
			logger.Logger.Println(ErrStr)
			goto ERROR
		}
		lunInfo := (err).(formatter.Lun)
		retLunList = append(retLunList, lunInfo)
		iscsiStatus, err := handler.WriteIscsiConfigObject(*volume)
		if !iscsiStatus {
			ErrStr += "Failed : Write Iscsi Config Object " + err.(error).Error()
			ErrCode = hcc_errors.CelloInternalPrepareDeploy
			logger.Logger.Println(ErrStr)

			goto ERROR
		}
		logger.Logger.Println("WriteIscsiConfigObject : ", err)

	default:
		ErrStr += "Use Type Of Volume Invalid"
		ErrCode = hcc_errors.CelloGrpcArgumentError
		logger.Logger.Println(ErrStr)
		goto ERROR
	}
	return retLunList, ErrCode, ErrStr

ERROR:
	ErrStr += "deleteAction(): Failed to delete volume " + ErrStr
	return retLunList, ErrCode, ErrStr

}

//VolumeHandler : Manipulate Volume Create
func VolumeHandler(contents *pb.ReqVolumeHandler) (*pb.Volume, uint64, string) {
	var recvErr error
	// var recvStr string
	var uuid string
	var modelVolume model.Volume
	var tempModelVolume model.Volume
	var ErrStr string
	var ErrCode uint64
	reformatPBReqtoModelVolume(contents, &modelVolume)
	retPbVolume := reformatPBReqtoPBVolume(contents)
	logger.Logger.Println("Resolving: Volume Handle")
	if len(retPbVolume.ServerUUID) == 0 || len(retPbVolume.UserUUID) == 0 {
		ErrStr += "reformatPBReqtoPBVolume() : " + "Invalid UUID : { Server: " + "\n User : " + "}"
		ErrCode = hcc_errors.CelloGrpcArgumentError
		logger.Logger.Println(ErrStr)
		goto ERROR
	}
	recvErr = handler.ReloadAllOfVolInfo()
	if recvErr != nil {
		// errStack.Push(&hccerr.HccError{ErrText: "Can't Preload "})
		ErrStr += "handl.() : " + recvErr.Error()
		ErrCode = hcc_errors.CelloInternalReloadObject
		logger.Logger.Println(ErrStr)
		goto ERROR
	}
	recvErr = handler.ReloadPoolObject()
	if recvErr != nil {
		// errStack.Push(&hccerr.HccError{ErrText: "Can't Reload Object"})
		ErrStr += "ReloadPoolObject() : " + recvErr.Error()
		ErrCode = hcc_errors.CelloInternalReloadObject
		logger.Logger.Println(ErrStr)
		goto ERROR
	}

	switch strings.ToLower(retPbVolume.Action) {
	case "create":

		out, recvErr := gouuid.NewV4()
		if recvErr != nil {
			ErrStr += "create(Can't Create Volume UUID) : " + recvErr.Error()
			ErrCode = hcc_errors.CelloInternalVolumeHandleError
			logger.Logger.Println(ErrStr)
			goto ERROR
		}
		uuid = out.String()

		modelVolume.UUID = uuid
		retPbVolume.UUID = modelVolume.UUID

		ErrCode, tempErr := createAction(retPbVolume, &modelVolume)
		if len(tempErr) > 0 {
			ErrStr += "Error createAction() : " + tempErr + ", ErrCode :" + strconv.FormatUint(ErrCode, 10)
			logger.Logger.Println(ErrStr)
			goto ERROR
		}
		retPbVolume.Lun = int64(modelVolume.LunNum)

		ErrCode, recvStr := dao.CreateVolume(&modelVolume)
		if recvStr != "" {
			ErrStr += "Error DB dao.CreateVolume() : " + recvStr + ", ErrCode : " + strconv.FormatUint(ErrCode, 10)
			logger.Logger.Println(ErrStr)
			goto ERROR
		}
		tempVolume, ErrCode, recvStr := dao.ReadVolume(&modelVolume)
		if tempVolume.UUID == "" {
			ErrStr += "Error DB dao.ReadVolume() : " + modelVolume.UUID + " is Not Exist " + recvStr + ", ErrCode : " + strconv.FormatUint(ErrCode, 10)
			logger.Logger.Println(ErrStr)
			goto ERROR
		}
		retPbVolume = reformatModelVolumetoPBVolume(&tempVolume)
		logger.Logger.Println("[Create Volume] Success ")

	case "read":
		if retPbVolume.UUID == "" {
			ErrStr += "Invalid UUID"
			ErrCode = hcc_errors.CelloGrpcArgumentError
			logger.Logger.Println(ErrStr)
			goto ERROR
		}

		tempVolume, ErrCode, recvErr := dao.ReadVolume(&modelVolume)
		if tempVolume.UUID == "" {
			ErrStr += "Error DB dao.ReadVolume() : " + modelVolume.UUID + " is Not Exist " + recvErr
			logger.Logger.Println(ErrStr, strconv.FormatUint(ErrCode, 10))
			goto ERROR
		}
		retPbVolume = reformatModelVolumetoPBVolume(&tempVolume)

	case "update":

	case "delete":
		lunList, ErrCode, tempErr := deleteAction(retPbVolume, &modelVolume)
		if len(tempErr) > 0 {
			ErrStr += "Error deleteAction() : " + tempErr + ", ErrCode :" + strconv.FormatUint(ErrCode, 10)
			logger.Logger.Println(ErrStr)
			goto ERROR
		}
		for i, args := range lunList {
			zfsDataSetVolName := strings.Split(args.Path, "/")
			deleteVolStatus, recvErr := handler.DeleteVolumeZFS(args.Pool + "/" + zfsDataSetVolName[len(zfsDataSetVolName)-1])
			fmt.Println("Delete ", i, " : ", args)
			if !deleteVolStatus {
				ErrStr += "Error DeleteVolumeZFS() : " + recvErr.(string)
				logger.Logger.Println(ErrStr)
				goto ERROR
			}
			tempModelVolume.UUID = args.UUID
			ErrCode, recvErr := dao.DeleteVolume(&tempModelVolume)
			if recvErr != nil {
				ErrStr += "Error DB dao.DeleteVolume() : " + recvErr.(string) + ", ErrCode : " + strconv.FormatUint(ErrCode, 10)
				logger.Logger.Println(ErrStr)
				goto ERROR
			}
		}

		logger.Logger.Println("[Delete Volume] Success ")

	default:
		ErrStr += "Invalid Action : " + retPbVolume.Action
		ErrCode = hcc_errors.CelloGrpcArgumentError
		goto ERROR
	}
	logger.Logger.Println("retPbVolume : ", retPbVolume)
	return retPbVolume, 0, ""

ERROR:
	ErrStr += "VolumeHandler(): Failed to handle volume {" + ErrStr + "code : " + strconv.FormatUint(ErrCode, 10) + "}"

	return nil, ErrCode, ErrStr
}

// GetPoolList : pool list
func GetPoolList(contents *pb.ReqGetPoolList) ([]*pb.Pool, uint64, string) {
	var recvErr error
	// var recvStr string
	var ErrStr string
	var ErrCode uint64
	var retPbPoolList []*pb.Pool
	singlePbPool := contents.GetPool()
	retPbPool := &pb.Pool{
		UUID:          singlePbPool.GetUUID(),
		Size:          singlePbPool.GetSize(),
		Free:          singlePbPool.GetFree(),
		Capacity:      singlePbPool.GetCapacity(),
		Health:        singlePbPool.GetHealth(),
		Name:          singlePbPool.GetName(),
		AvailableSize: singlePbPool.GetAvailableSize(),
		Action:        singlePbPool.GetAction(),
		Used:          singlePbPool.GetUsed(),
	}
	logger.Logger.Println("Resolving: Pool List")
	// if retPbPool.Name == "" {
	// 	errStack.Push(&hccerr.HccError{ErrCode: hccerr.CelloGrpcArgumentError, ErrText: "Invalid Pool name : " + retPbPool.Name + "}"})
	// 	goto ERROR
	// }
	recvErr = handler.ReloadAllOfVolInfo()
	if recvErr != nil {
		ErrStr += "ReloadAllOfVolInfo() : " + recvErr.Error()
		ErrCode = hcc_errors.CelloInternalReloadObject
		logger.Logger.Println(ErrStr)
		goto ERROR
	}
	recvErr = handler.ReloadPoolObject()
	if recvErr != nil {
		ErrStr += "ReloadPoolObject() : " + recvErr.Error()
		ErrCode = hcc_errors.CelloInternalReloadObject
		logger.Logger.Println(ErrStr)
		goto ERROR
	}

	switch strings.ToLower(retPbPool.Action) {

	case "read":

		for _, args := range formatter.PoolObjectMap.PoolMap {
			if strings.Contains(args.Name, config.VolumeConfig.VOLUMEPOOL) {
				var tempretPbPool pb.Pool
				tempretPbPool.AvailableSize = args.AvailableSize
				tempretPbPool.Capacity = args.Capacity
				tempretPbPool.Free = args.Free
				tempretPbPool.Size = args.Size
				tempretPbPool.Health = args.Health
				tempretPbPool.Name = args.Name
				tempretPbPool.Used = args.Used
				retPbPoolList = append(retPbPoolList, &tempretPbPool)
			}
		}

	default:
		ErrStr += "Invalid Action : " + retPbPool.Action
		ErrCode = hcc_errors.CelloGrpcArgumentError
		goto ERROR
	}
	logger.Logger.Println("retPbPoolList : ", retPbPoolList)

	return retPbPoolList, ErrCode, ErrStr

ERROR:
	ErrStr += "GetPoolList(): Failed to handle Pool {" + ErrStr + "code : " + strconv.FormatUint(ErrCode, 10) + "}"

	return nil, ErrCode, ErrStr
}
func PoolHandler(contents *pb.ReqPoolHandler) (*pb.Pool, uint64, string) {
	var recvErr error
	// var recvStr string
	var uuid string
	var modelPool model.Pool
	var ErrStr string
	var ErrCode uint64
	reformatPBReqtoModelPool(contents, &modelPool)
	retPbPool := reformatPBReqtoPBPool(contents)
	logger.Logger.Println("Resolving: Pool Handle")

	recvErr = handler.ReloadAllOfVolInfo()
	if recvErr != nil {
		// errStack.Push(&hccerr.HccError{ErrText: "Can't Preload "})
		ErrStr += "ReloadAllOfVolInfo() : " + recvErr.Error()
		ErrCode = hcc_errors.CelloInternalReloadObject
		logger.Logger.Println(ErrStr)
		goto ERROR
	}
	recvErr = handler.ReloadPoolObject()
	if recvErr != nil {
		ErrStr += "ReloadPoolObject() : " + recvErr.Error()
		ErrCode = hcc_errors.CelloInternalReloadObject
		logger.Logger.Println(ErrStr)
		goto ERROR
	}

	switch strings.ToLower(retPbPool.Action) {
	case "create":

		out, recvErr := gouuid.NewV4()
		if recvErr != nil {
			ErrStr += "create(Can't Create Pool UUID) : " + recvErr.Error()
			ErrCode = hcc_errors.CelloGrpcArgumentError
			logger.Logger.Println(ErrStr)
			goto ERROR
		}
		uuid = out.String()

		modelPool.UUID = uuid
		retPbPool.UUID = modelPool.UUID

	case "read":
		for _, args := range formatter.PoolObjectMap.PoolMap {
			fmt.Println("formatter.PoolObjectMap.PoolMap\n\n", args)
			if strings.Contains(args.Name, config.VolumeConfig.VOLUMEPOOL) {
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
		ErrStr += "Invalid Action : " + retPbPool.Action
		ErrCode = hcc_errors.CelloGrpcArgumentError
		goto ERROR
	}
	logger.Logger.Println("retPbPool : ", retPbPool)
	return retPbPool, 0, ""

ERROR:
	ErrStr += "PoolHandler(): Failed to handle Pool {" + ErrStr + "code : " + strconv.FormatUint(ErrCode, 10) + "}"

	return nil, ErrCode, ErrStr

}

//GetVolumeList : Manipulate Volume Create
func GetVolumeList(contents *pb.ReqGetVolumeList) ([]*pb.Volume, uint64, string) {
	var recvErr error
	// var recvStr string
	var ErrStr string
	var ErrCode uint64
	var modelVolume model.Volume
	var retPbVolumeList []*pb.Volume

	tempContents := contents.GetVolume()
	// modelVolume.ServerUUID = contents
	// retPbVolume := reformatPBReqtoPBVolume(contents)
	logger.Logger.Println("Resolving: GetVolumeList ", tempContents.ServerUUID)
	if len(tempContents.ServerUUID) == 0 {
		ErrStr += "GetVolume() : " + "Invalid UUID : { Server: " + "\n User : " + "}"
		ErrCode = hcc_errors.CelloGrpcArgumentError
		logger.Logger.Println(ErrStr)
		goto ERROR
	}
	recvErr = handler.ReloadAllOfVolInfo()
	if recvErr != nil {
		ErrStr += "ReloadAllOfVolInfo() : " + recvErr.Error()
		ErrCode = hcc_errors.CelloInternalReloadObject
		logger.Logger.Println(ErrStr)
		goto ERROR
	}
	recvErr = handler.ReloadPoolObject()
	if recvErr != nil {
		ErrStr += "ReloadPoolObject() : " + recvErr.Error()
		ErrCode = hcc_errors.CelloInternalReloadObject
		logger.Logger.Println(ErrStr)
		goto ERROR
	}
	switch strings.ToLower(tempContents.Action) {
	case "read_list":
		modelVolume.ServerUUID = tempContents.ServerUUID
		if tempContents.UserUUID != "" {
			modelVolume.UserUUID = tempContents.UserUUID
		}
		// fmt.Println("Codex : ", int(contents.GetRow()), int(contents.GetPage()))
		tempVolume, ErrCode, recvStr := dao.ReadVolumeList(&modelVolume, int(contents.GetRow()), int(contents.GetPage()))
		if tempVolume == nil {
			ErrStr += "Error DB : " + modelVolume.UUID + " is Not Exist " + recvStr + ", ErrCode : " + strconv.FormatUint(ErrCode, 10)
			logger.Logger.Println(ErrStr)
			goto ERROR
		}
		retPbVolumeList = reformatModelVolumetoPBVolumeList(&tempVolume)

	case "update":

	case "delete":

	default:
		ErrStr += "Invalid Action : " + contents.GetVolume().Action
		ErrCode = hcc_errors.CelloGrpcArgumentError
		goto ERROR
	}
	logger.Logger.Println("retPbVolumeList : ", retPbVolumeList)
	return retPbVolumeList, ErrCode, ErrStr

ERROR:
	ErrStr += "GetVolumeList(): Failed to get volume List {" + ErrStr + "code : " + strconv.FormatUint(ErrCode, 10) + "}"

	return nil, ErrCode, ErrStr
}
