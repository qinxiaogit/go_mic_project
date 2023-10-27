package login

import (
	"github.com/gin-gonic/gin"
	"github.com/qinxiaogit/go_mic_project/frontend_api/route"
	"net/http"
)

type service struct{}

func (s service) RegisterRoute(router gin.IRoutes) {
	router.POST("/api/user/login", s.Login)
	router = router.Use(route.AuthRequired())
	{
		router.POST("/api/user/logout", s.Logout)
	}
}

type loginResponse struct {
	Token string `json:"token" binding:"required"`
}

func (s service) Login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, loginResponse{
		Token: "123",
	})
}
