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

func CheckTweetOwner(ctx *gin.Context) {
	db := database.GetDB()
	username := ctx.Query("username")
	password := ctx.Query("password")
	tweet_id := ctx.Query("tweet_id")
	if len(username)==0 || len(password)==0 || len(tweet_id)==0 {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, fmt.Errorf("error in query"), "enter username, password and tweet_id"))
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
	tweet := models.Tweet{}
	err = tx.Model(&models.Tweet{}).Where("id = ? AND deleted_by is null", tweet_id).First(&tweet).Error
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "invalid tweet"))
		return
	}
	if tweet.UserId != user.Id {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, fmt.Errorf("error in user and tweet"), "the tweet isnt for the user"))
		return
	}
	ctx.Set("tweet_id", tweet_id)
	ctx.Set("modified_by", user.Id)
	ctx.Next()
}