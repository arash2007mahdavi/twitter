package middlewares

import (
	"fmt"
	"net/http"
	"twitter/src/database"
	"twitter/src/database/models"
	"twitter/src/responses"

	"github.com/gin-gonic/gin"
)

func UserCheck(ctx *gin.Context) {
	username := ctx.Query("username")
	password := ctx.Query("password")
	if len(username) == 0 || len(password) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, responses.GenerateResponseWithError(http.StatusNotFound, nil, "username, password didnt entered"))
		return
	}
	var count int64
	db := database.GetDB()
	db.Model(&models.User{}).
		Where("username = ? AND deleted_by is null", username).Count(&count)
	if count > 0 {
		var count1 int64
		db.Model(&models.User{}).
			Where("username = ? AND password = ? AND deleted_by is null", username, password).Count(&count1)
		if count1 > 0 {
			ctx.Set("username", username)
			ctx.Next()
		} else {
			ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, fmt.Errorf("wrong password"), "wrong password"))
			return
		}
	} else {
		ctx.AbortWithStatusJSON(http.StatusNotFound, responses.GenerateResponseWithError(http.StatusNotFound, fmt.Errorf("invalid username"), "invalid username"))
		return
	}
}