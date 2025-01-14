package api

import (
	"github.com/gin-gonic/gin"
	_ "github.com/skyrocketOoO/serverx/docs/openapi"
	"github.com/skyrocketOoO/serverx/internal/controller"
	"github.com/skyrocketOoO/serverx/internal/controller/middleware"
	"github.com/skyrocketOoO/serverx/internal/global"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Bind(r *gin.Engine, h *controller.Handler) {
	vr := r.Group(global.ApiVersion)
	vr.GET("/ping", h.Ping)
	vr.GET("/healthy", h.Healthy)
	vr.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	vr.POST("/login", h.Login)
	vr.POST("/register", h.Register)

	pR := vr.Group("/")
	pR.Use(middleware.Jwt())
}
