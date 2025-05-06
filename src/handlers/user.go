package handlers

import (
	"net/http"
	"twitter/src/responses"
	"twitter/src/dtos"
	"twitter/src/logger"
	"twitter/src/services"

	"github.com/gin-gonic/gin"
)

type UserHelper struct {
	Logger logger.Logger
	Service *services.UserService
}

func GetUserHelper() *UserHelper {
	return &UserHelper{
		Logger: logger.NewLogger(),
		Service: services.NewUserService(),
	}
}

func (h *UserHelper) GetOtp(ctx *gin.Context) {

}

func (h *UserHelper) NewUser(ctx *gin.Context) {
	req := dtos.UserCreate{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, responses.GenerateResponseWithValidationError(http.StatusBadRequest, err, ""))
		return
	}
	res, err := h.Service.Create(ctx, &req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in creating new user"))
		return
	}

}