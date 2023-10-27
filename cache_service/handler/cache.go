package handler

import (
	"context"
	"fmt"
	"github.com/go-micro/plugins/v4/cache/redis"
	pb "github.com/qinxiaogit/go_mic_project/cache_service/proto"
	"go-micro.dev/v4/cache"
	log "go-micro.dev/v4/logger"
	"time"
)

type CacheService struct {
	cache cache.Cache
}

func NewCache(opts ...cache.Option) *CacheService {
	c := redis.NewCache(opts...)
	return &CacheService{c}
}

func (c *CacheService) Get(ctx context.Context, in *pb.GetRequest, out *pb.GetResponse) error {
	log.Info("Received a cache.Get request:%v", in)

	v, e, err := c.cache.Get(ctx, in.Key)
	if err != nil {
		return err
	}
	vby, er := v.([]byte)
	if er {
		out.Value = string(vby)
	} else {
		fmt.Sprintf("%v", v)
	}
	out.Expiration = e.String()
	return nil
}

func (c *CacheService) Put(ctx context.Context, in *pb.PutRequest, out *pb.PutResponse) error {
	log.Info("Received a cache.Put request:%v", in)

	d, err := time.ParseDuration(in.Duration)
	if err != nil {
		return err
	}
	if err := c.cache.Put(ctx, in.Key, in.Value, d); err != nil {
		return nil
	}
	return nil
}

func (c *CacheService) Delete(ctx context.Context, in *pb.DeleteRequest, out *pb.DeleteResponse) error {
	log.Info("Received a cache.Delete request:%v", in)

	if err := c.cache.Delete(ctx, in.Key); err != nil {
		return err
	}
	return nil
}
