package server

import (
	"context"
	"hcc/cello/action/grpc/errconv"
	"hcc/cello/driver/grpcsrv"
	"hcc/cello/lib/logger"

	errh "innogrid.com/hcloud-classic/hcc_errors"
	"innogrid.com/hcloud-classic/pb"
)

type celloServer struct {
	pb.UnimplementedCelloServer
}

func returnVolume(volume *pb.Volume) *pb.Volume {
	return &pb.Volume{
		UUID:       volume.UUID,
		Size:       volume.Size,
		Filesystem: volume.Filesystem,
		ServerUUID: volume.ServerUUID,
		UseType:    volume.UseType,
		UserUUID:   volume.UserUUID,
		Network_IP: volume.Network_IP,
		GatewayIp:  volume.GatewayIp,
		Pool:       volume.Pool,
		Lun:        volume.Lun,
		CreatedAt:  volume.CreatedAt,
	}
}

func returnPool(pool *pb.Pool) *pb.Pool {
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

func (s *celloServer) VolumeHandler(_ context.Context, in *pb.ReqVolumeHandler) (*pb.ResVolumeHandler, error) {
	logger.Logger.Println("Request received: Volume Handler()")
	// fmt.Println("Grpc : \n", &pb.ResVolumeHandler{Volume: &pb.Volume{}, HccErrorStack: errconv.HccStackToGrpc(nil)})
	volume, errCode, errStr := grpcsrv.VolumeHandler(in)
	if volume == nil {
		errStack := errh.NewHccErrorStack(errh.NewHccError(errCode, errStr))

		return &pb.ResVolumeHandler{Volume: &pb.Volume{}, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResVolumeHandler{Volume: returnVolume(volume)}, nil
}

func (s *celloServer) PoolHandler(_ context.Context, in *pb.ReqPoolHandler) (*pb.ResPoolHandler, error) {
	logger.Logger.Println("Request received: Pool Handler()")
	// fmt.Println("Grpc : \n", &pb.ResVolumeHandler{Volume: &pb.Volume{}, HccErrorStack: errconv.HccStackToGrpc(nil)})
	pool, errCode, errStr := grpcsrv.PoolHandler(in)
	if pool == nil {
		errStack := errh.NewHccErrorStack(errh.NewHccError(errCode, errStr))

		return &pb.ResPoolHandler{Pool: &pb.Pool{}, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResPoolHandler{Pool: pool}, nil
}

func (s *celloServer) GetVolumeList(_ context.Context, in *pb.ReqGetVolumeList) (*pb.ResGetVolumeList, error) {
	logger.Logger.Println("Request received: GetVolumeList()")
	// fmt.Println("Grpc : \n", &pb.ResVolumeHandler{Volume: &pb.Volume{}, HccErrorStack: errconv.HccStackToGrpc(nil)})
	volumeList, errCode, errStr := grpcsrv.GetVolumeList(in)
	if volumeList == nil {
		errStack := errh.NewHccErrorStack(errh.NewHccError(errCode, errStr))

		return &pb.ResGetVolumeList{Volume: nil, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResGetVolumeList{Volume: volumeList}, nil
}

func (s *celloServer) GetPoolList(_ context.Context, in *pb.ReqGetPoolList) (*pb.ResGetPoolList, error) {
	logger.Logger.Println("Request received: GetPoolList()")
	// fmt.Println("Grpc : \n", &pb.ResVolumeHandler{Volume: &pb.Volume{}, HccErrorStack: errconv.HccStackToGrpc(nil)})
	pool, errCode, errStr := grpcsrv.GetPoolList(in)
	if pool == nil {
		errStack := errh.NewHccErrorStack(errh.NewHccError(errCode, errStr))

		return &pb.ResGetPoolList{Pool: nil, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResGetPoolList{Pool: pool}, nil
}
