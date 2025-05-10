package middlewares

import (
	"fmt"
	"net/http"
	"twitter/src/database"
	"twitter/src/database/models"
	"twitter/src/responses"

	"github.com/gin-gonic/gin"
)

func CheckForUpdate(ctx *gin.Context) {
	username := ctx.Query("username")
	password := ctx.Query("password")
	modified_by := ctx.Query("modified_by")
	if len(username)==0 || len(password)==0 || len(modified_by)==0 {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, fmt.Errorf("invalid query"), "enter username, password and modified_by"))
		return
	}
	db := database.GetDB()
	user := models.User{}
	err := db.Model(&models.User{}).Where("username = ? AND password = ? AND deleted_at is null", username, password).First(&user).Error
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "invalid User"))
		return
	}
	ctx.Set("user_id", user.Id)
	ctx.Set("modified_by", modified_by)
	ctx.Next()
}