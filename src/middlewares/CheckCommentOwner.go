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

func CheckCommentOwner(ctx *gin.Context) {
	db := database.GetDB()
	username := ctx.Query("username")
	password := ctx.Query("password")
	comment_id := ctx.Query("comment_id")
	if len(username)==0 || len(password)==0 || len(comment_id)==0 {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, fmt.Errorf("error in query"), "enter username, password and tweet_id"))
		return
	}
	user := models.User{}
	tx := db.WithContext(ctx).Begin()
	err := tx.Model(&models.User{}).Where("username = ? AND deleted_at is null", username).First(&user).Error
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
	comment := models.Comment{}
	err = tx.Model(&models.Comment{}).Where("id = ? AND deleted_at is null", comment_id).First(&comment).Error
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "invalid comment"))
		return
	}
	if comment.UserId != user.Id {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, fmt.Errorf("error in user and comment"), "the comment isnt for the user"))
		return
	}
	ctx.Set("comment_id", comment_id)
	ctx.Set("modified_by", user.Id)
	ctx.Next()
}