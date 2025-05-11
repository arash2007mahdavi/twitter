package middlewares

import (
	"fmt"
	"net/http"
	"twitter/src/logger"
	"twitter/src/responses"

	"github.com/gin-gonic/gin"
)

var log = logger.NewLogger()

func CheckAdmin(ctx *gin.Context) {
	admin_username := ctx.GetHeader("admin")
	admin_password := ctx.GetHeader("password")
	if admin_username == "arash2007mahdavi" && admin_password == "arash2306" {
		log.Info(logger.Admin, logger.Enter, "admin entered", map[logger.ExtraCategory]interface{}{logger.Username: admin_username})
		ctx.Set("admin_username", admin_username)
		ctx.Next()
		return
	}
	ctx.AbortWithStatusJSON(http.StatusLocked, responses.GenerateResponseWithError(http.StatusLocked, fmt.Errorf("invalid admin information"), "invalid admin"))
}