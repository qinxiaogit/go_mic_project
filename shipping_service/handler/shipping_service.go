package handler

import (
	"context"
	"fmt"
	pb "github.com/qinxiaogit/go_mic_project/shiping_service/proto"
	"go-micro.dev/v4/logger"
)

type ShippingService struct {
}

func (s ShippingService) GetQuote(ctx context.Context, in *pb.GetQuoteRequest, out *pb.GetQuoteResponse) error {
	logger.Info("[GetQuo] received request")
	defer logger.Info("[GetQuo] completed request")
	//1. Generate a quote base on the total number of items to be shipped
	quote := CreateQuoteFromCount(0)
	out.CostUsd = &pb.Money{CurrencyCode: "USD",
		Units: int64(quote.Dollars),
		Nanos: int32(quote.Cents),
	}
	return nil
}

func (s *ShippingService) ShipOrder(ctx context.Context, in *pb.ShipOrderRequest, out *pb.ShipOrderResponse) error {
	logger.Info("[ShipOrder] received request")
	defer logger.Info("[ShipOrder] completed request")
	// 1. Create a Tracing ID
	baseAddress := fmt.Sprintf("%s, %s, %s", in.Address.StreetAddress, in.Address.City, in.Address.State)
	id := CreateTrackingId(baseAddress)
	out.TrackingId = id
	return nil
}
