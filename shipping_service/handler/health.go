package handler

import (
	"context"
	pb "github.com/qinxiaogit/go_mic_project/shiping_service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Health struct {
}

func (h *Health) Check(ctx context.Context, in *pb.HealthCheckRequest, out *pb.HealthCheckResponse) error {
	out.Status = pb.HealthCheckResponse_SERVING
	return nil
}

func (h *Health) Watch(ctx context.Context, req *pb.HealthCheckRequest, stream pb.Health_WatchStream) error {
	return status.Errorf(codes.Unimplemented, "health check via Watch not implemented")
}
