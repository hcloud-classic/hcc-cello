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

func (s *celloServer) VolumeHandler(_ context.Context, in *pb.ReqVolumeHandler) (*pb.ResVolumeHandler, error) {
	logger.Logger.Println("Request received: CreateVolume()")
	// fmt.Println("Grpc : \n", &pb.ResVolumeHandler{Volume: &pb.Volume{}, HccErrorStack: errconv.HccStackToGrpc(nil)})
	volume, errStack := grpcsrv.VolumeHandler(in)
	if volume == nil {
		return &pb.ResVolumeHandler{Volume: &pb.Volume{}, HccErrorStack: errconv.HccStackToGrpc(errStack)}, nil
	}

	return &pb.ResVolumeHandler{Volume: returnVolume(volume)}, nil
}
