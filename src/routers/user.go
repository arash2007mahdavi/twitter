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
	r.PUT("/update", middlewares.GetIdForUpdate, h.UpdateUser)
	r.DELETE("/delete", middlewares.GetIdForDelete, h.DeleteUser)
	r.GET("/get/users", middlewares.CheckAdmin, h.GetUsers)
	r.GET("/get/profile", middlewares.GetIdFromUsername, h.GetProfile)
	r.POST("/follow", h.Follow)
}