package middlewares

import (
	"fmt"
	"net/http"
	"twitter/src/responses"

	"github.com/gin-gonic/gin"
)

func CheckAdmin(ctx *gin.Context) {
	admin_username := ctx.GetHeader("admin")
	admin_password := ctx.GetHeader("password")
	if admin_username == "arash2007mahdavi" && admin_password == "@rash2007" {
		ctx.Next()
		return
	}
	ctx.AbortWithStatusJSON(http.StatusLocked, responses.GenerateResponseWithError(http.StatusLocked, fmt.Errorf("invalid admin information"), "invalid admin"))
}