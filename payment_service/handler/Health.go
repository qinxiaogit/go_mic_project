package handler

import (
	"context"
	pb "github.com/qinxiaogit/go_mic_project/payment_service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Health struct {
}

func (h *Health) Check(ctx context.Context, req *pb.HealthCheckRequest, resp *pb.HealthCheckResponse) error {
	if req.Service == "aaaa" {
		resp.Status = pb.HealthCheckResponse_SERVER_UNKNOWN

	} else {
		resp.Status = pb.HealthCheckResponse_SERVING

	}
	return nil
}

func (h *Health) Watch(ctx context.Context, req *pb.HealthCheckRequest, stream pb.Health_WatchStream) error {
	return status.Errorf(codes.Unimplemented, "health check via watch not implemented")
}
