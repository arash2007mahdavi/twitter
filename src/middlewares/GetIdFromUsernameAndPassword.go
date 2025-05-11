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

func GetIdFromUsernameAndPassword(ctx *gin.Context) {
	username := ctx.Query("username")
	password := ctx.Query("password")
	if len(username) == 0 || len(password) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, responses.GenerateResponseWithError(http.StatusNotFound, fmt.Errorf("error in query"), "username, password didnt entered"))
		return
	}
	var user models.User
	db := database.GetDB()
	db.Model(&models.User{}).Where("username = ? AND deleted_by is null", username).First(&user)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, responses.GenerateResponseWithError(http.StatusNotFound, err, "invalid user"))
		return
	}
	ctx.Set("user_id", user.Id)
	ctx.Next()
}