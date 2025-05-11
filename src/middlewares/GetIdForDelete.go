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

func GetIdForDelete(ctx *gin.Context) {
	username := ctx.Query("username")
	password := ctx.Query("password")
	deleted_by := ctx.Query("deleted_by")
	if len(deleted_by) == 0 {
		deleted_by = "0"
	}
	if len(username)==0 || len(password)==0 {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, fmt.Errorf("invalid query"), "enter username, password"))
		return
	}
	db := database.GetDB()
	user := models.User{}
	db.Model(&models.User{}).Where("username = ? AND deleted_at is null", username).First(&user)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "invalid User"))
		return
	}
	ctx.Set("user_id", user.Id)
	ctx.Set("deleted_by", deleted_by)
	ctx.Next()
}