package service

import (
	"github.com/qinxiaogit/go_mic_project/frontend_api/config"
	hipstershop "github.com/qinxiaogit/go_mic_project/frontend_api/proto"
	"go-micro.dev/v4/client"
)

type baseService struct {
	cacheService    hipstershop.CacheService
	shippingService hipstershop.ShippingService
}

type GroupService struct {
	baseService
	LoginService *LoginService
}

func NewGroupService(c client.Client) *GroupService {
	cfg := config.Get()
	return &GroupService{
		baseService: baseService{
			cacheService:    hipstershop.NewCacheService(cfg.CacheService, c),
			shippingService: hipstershop.NewShippingService(cfg.ShippingService, c),
		},
		LoginService: new(LoginService),
	}
}
