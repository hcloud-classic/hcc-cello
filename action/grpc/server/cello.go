package server

import (
	"context"
	errconv "hcc/cello/action/grpc/errconv"
	pb "hcc/cello/action/grpc/pb/rpccello"
	"hcc/cello/driver/grpcsrv"
	"hcc/cello/lib/logger"
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
	volume, errStack := grpcsrv.VolumeHandler(in)
	if volume == nil {
		return &pb.ResVolumeHandler{Volume: &pb.Volume{}, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResVolumeHandler{Volume: returnVolume(volume)}, nil
}

func (s *celloServer) PoolHandler(_ context.Context, in *pb.ReqPoolHandler) (*pb.ResPoolHandler, error) {
	logger.Logger.Println("Request received: Pool Handler()")
	// fmt.Println("Grpc : \n", &pb.ResVolumeHandler{Volume: &pb.Volume{}, HccErrorStack: errconv.HccStackToGrpc(nil)})
	pool, errStack := grpcsrv.PoolHandler(in)
	if pool == nil {
		return &pb.ResPoolHandler{Pool: &pb.Pool{}, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResPoolHandler{Pool: pool}, nil
}

func (s *celloServer) GetVolumeList(_ context.Context, in *pb.ReqGetVolumeList) (*pb.ResGetVolumeList, error) {
	logger.Logger.Println("Request received: GetVolumeList()")
	// fmt.Println("Grpc : \n", &pb.ResVolumeHandler{Volume: &pb.Volume{}, HccErrorStack: errconv.HccStackToGrpc(nil)})
	volumeList, errStack := grpcsrv.GetVolumeList(in)
	if volumeList == nil {
		return &pb.ResGetVolumeList{Volume: nil, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResGetVolumeList{Volume: volumeList}, nil
}

func (s *celloServer) GetPoolList(_ context.Context, in *pb.ReqGetPoolList) (*pb.ResGetPoolList, error) {
	logger.Logger.Println("Request received: GetPoolList()")
	// fmt.Println("Grpc : \n", &pb.ResVolumeHandler{Volume: &pb.Volume{}, HccErrorStack: errconv.HccStackToGrpc(nil)})
	pool, errStack := grpcsrv.GetPoolList(in)
	if pool == nil {
		return &pb.ResGetPoolList{Pool: nil, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResGetPoolList{Pool: pool}, nil
}
