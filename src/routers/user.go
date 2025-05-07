package routers

import (
	"twitter/src/handlers"
	"twitter/src/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.RouterGroup) {
	h := handlers.GetUserHelper()
	r.GET("/get/otp", h.GetOtp)
	r.POST("/new", h.NewUser)
	r.PUT("/update", h.UpdateUser)
	r.DELETE("/delete", h.DeleteUser)
	r.GET("/get/all", middlewares.CheckAdmin, h.GetUsers)
	r.GET("/get/profile", h.GetProfile)
}