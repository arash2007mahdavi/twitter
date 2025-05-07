package handlers

import (
	"fmt"
	"net/http"
	"twitter/src/dtos"
	"twitter/src/logger"
	"twitter/src/responses"
	"twitter/src/services"

	"github.com/gin-gonic/gin"
)

type UserHelper struct {
	Logger logger.Logger
	Service *services.UserService
	Otp *services.OtpService
}

func GetUserHelper() *UserHelper {
	return &UserHelper{
		Logger: logger.NewLogger(),
		Service: services.NewUserService(),
		Otp: services.NewOtpService(),
	}
}

type OtpDto struct {
	MobileNumber string `json:"mobileNumber" binding:"required,mobile"`
}

func (h *UserHelper) GetOtp(ctx *gin.Context) {
	req := OtpDto{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, responses.GenerateResponseWithValidationError(http.StatusBadRequest, err, ""))
		return
	}
	otp := services.MakeOtp()
	err = h.Otp.SetOtp(req.MobileNumber, otp)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in seting otp"))
		return
	}
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, map[string]string{"otp": otp}, "otp set successfuly"))
}

func (h *UserHelper) NewUser(ctx *gin.Context) {
	req := dtos.UserCreate{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, responses.GenerateResponseWithValidationError(http.StatusBadRequest, err, ""))
		return
	}
	test_otp := ctx.Query("otp")
	if len(test_otp) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "otp doesnt entered"))
		return
	}
	err = h.Otp.ValidateOtp(req.MobileNumber, test_otp)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in validate otp"))
		return
	}
	res, err := h.Service.Create(ctx, &req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in creating new user"))
		return
	}
	h.Logger.Info(logger.User, logger.New, "new user added", nil)
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, res, "new user created successfuly"))
}

func (h *UserHelper) GetUsers(ctx *gin.Context) {
	users, err := h.Service.GetUsers(ctx)
	if err != nil || len(*users) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in get users"))
		return
	}
	h.Logger.Info(logger.Admin, logger.See, "admin saw users", nil)
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, users, "users list"))
}

func (h *UserHelper) GetProfile(ctx *gin.Context) {
	username := ctx.Query("username")
	password := ctx.Query("password")
	ctx.Set("Username", username)
	ctx.Set("Password", password)
	user, err := h.Service.GetProfile(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in get user"))
		return
	}
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, user, "your profile"))
}

func (h *UserHelper) UpdateUser(ctx *gin.Context) {
	req := dtos.UserUpdate{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, responses.GenerateResponseWithValidationError(http.StatusBadRequest, err, ""))
		return
	}
	id := ctx.Query("id")
	modified_by := ctx.Query("modified_by")
	if len(id) == 0 || len(modified_by) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, fmt.Errorf("invalid id and modified_by"), "invalid id and modified_by"))
		return
	}
	ctx.Set("id", id)
	ctx.Set("modified_by", modified_by)
	res, err := h.Service.Base.Update(ctx, &req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in updating model"))
		return
	}
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, res, "model updated successfuly"))
}