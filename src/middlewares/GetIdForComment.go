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

func GetIdForComment(ctx *gin.Context) {
	username := ctx.Query("username")
	password := ctx.Query("password")
	tweet_id := ctx.Query("tweet_id")
	if len(username) == 0 || len(password) == 0 || len(tweet_id) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, fmt.Errorf("Missing required fields"), "error in query"))
		return
	}
	db := database.GetDB()
	tx := db.WithContext(ctx).Begin()
	user := models.User{}
	err := tx.Model(&models.User{}).Where("username = ? AND deleted_at is null", username).First(&user).Error
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "invalid user"))
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "invalid password"))
		return
	}
	tx.Commit()
	ctx.Set("user_id", user.Id)
	ctx.Set("tweet_id", tweet_id)
	ctx.Next()
}