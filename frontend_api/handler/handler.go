package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/qinxiaogit/go_mic_project/frontend_api/route"
	"go-micro.dev/v4/client"
)

type Options struct {
	Client client.Client
	Router *gin.Engine
}

func Register(opts Options) error {
	router := opts.Router
	for _, r := range []route.Registrar{} {
		r.RegisterRoute(router.Group(""))
	}
	return nil
}
