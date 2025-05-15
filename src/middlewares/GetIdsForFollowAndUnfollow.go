package middlewares

import (
	"fmt"
	"net/http"
	"twitter/src/database"
	"twitter/src/database/models"
	"twitter/src/responses"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetIdsForFollowAndUnfollow(ctx *gin.Context) {
	db := database.GetDB()
	username := ctx.Query("username")
	password := ctx.Query("password")
	target_username := ctx.Query("target_username")
	if len(username)==0 || len(password)==0 || len(target_username)==0 {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, fmt.Errorf("error in query"), "enter username, password and target_username"))
		return
	}
	if username == target_username {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, fmt.Errorf("error in request"), "invalid query"))
		return
	}
	user := models.User{}
	tx := db.WithContext(ctx).Begin()
	err := tx.Model(&models.User{}).Where("username = ? AND deleted_by is null", username).First(&user).Error
	if err != nil {
		tx.Rollback()
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "invalid user"))
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		tx.Rollback()
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "invalid password"))
		return
	}
	target_user := models.User{}
	err = tx.Model(&models.User{}).Where("username = ? AND deleted_by is null", target_username).First(&target_user).Error
	if err != nil {
		tx.Rollback()
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "invalid target user"))
		return
	}
	ctx.Set("user_id", user.Id)
	ctx.Set("target_id", target_user.Id)
	ctx.Next()
}