package middlewares

import (
	"fmt"
	"net/http"
	"twitter/src/database"
	"twitter/src/database/models"
	"twitter/src/responses"

	"github.com/gin-gonic/gin"
)

func GetIdFromUsername(ctx *gin.Context) {
	username := ctx.Query("username")
	if len(username) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, responses.GenerateResponseWithError(http.StatusNotFound, nil, "username didnt entered"))
		return
	}
	var count int64
	db := database.GetDB()
	db.Model(&models.User{}).
		Where("username = ? AND enabled is true", username).Count(&count)
	if count > 0 {
		user := models.User{}
		db.Model(&models.User{}).Where("username = ? AND enabled is true", username).First(&user)
		ctx.Set("user_id", user.Id)
		ctx.Next()
	} else {
		ctx.AbortWithStatusJSON(http.StatusNotFound, responses.GenerateResponseWithError(http.StatusNotFound, fmt.Errorf("invalid username"), "invalid username"))
		return
	}
}