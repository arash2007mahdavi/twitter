package middlewares

import (
	"fmt"
	"net/http"
	"strconv"
	"twitter/src/database"
	"twitter/src/database/models"
	"twitter/src/responses"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetIdForUpdate(ctx *gin.Context) {
	username := ctx.Query("username")
	password := ctx.Query("password")
	modified_by := ctx.Query("modified_by")
	if len(username)==0 || len(password)==0 {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, fmt.Errorf("invalid query"), "enter username, password"))
		return
	}
	db := database.GetDB()
	user := models.User{}
	db.Model(&models.User{}).Where("username = ? AND enabled is true", username).First(&user)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "invalid User"))
		return
	}
	ctx.Set("user_id", user.Id)
	user_id := strconv.Itoa(user.Id)
	if len(modified_by) != 0 {
		ctx.Set("modified_by", modified_by)
	} else {
		ctx.Set("modified_by", user_id)
	}
	ctx.Next()
}