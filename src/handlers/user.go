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
	MobileNumber string `json:"mobile_number" binding:"required,mobile"`
}

func (h *UserHelper) GetOtp(ctx *gin.Context) {
	req := OtpDto{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, responses.GenerateResponseWithValidationError(http.StatusBadRequest, err, "validation error"))
		return
	}
	otp := services.MakeOtp()
	err = h.Otp.SetOtp(req.MobileNumber, otp)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in seting otp"))
		return
	}
	h.Logger.Info(logger.Otp, logger.Set, "new otp set in redis", nil)
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, map[string]string{"otp": otp}, "otp set successfuly"))
}

func (h *UserHelper) NewUser(ctx *gin.Context) {
	req := dtos.UserCreate{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, responses.GenerateResponseWithValidationError(http.StatusBadRequest, err, "validation error"))
		return
	}
	test_otp := ctx.Query("otp")
	if len(test_otp) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "otp didnt enter"))
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
	if err != nil {
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
	h.Logger.Info(logger.User, logger.See, "user saw profile", nil)
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, user, "your profile"))
}

func (h *UserHelper) UpdateUser(ctx *gin.Context) {
	req := dtos.UserUpdate{}
	req.Enabled = true
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
	res, err := h.Service.Update(ctx, &req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in updating user"))
		return
	}
	h.Logger.Info(logger.User, logger.Update, "user updated successfuly", nil)
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, res, "user updated successfuly"))
}

func (h *UserHelper) DeleteUser(ctx *gin.Context) {
	id := ctx.Query("id")
	deleted_by := ctx.Query("deleted_by")
	if len(id) == 0 || len(deleted_by) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, fmt.Errorf("invalid id and deleted_by"), "invalid id and deleted_by"))
		return
	}
	ctx.Set("id", id)
	ctx.Set("deleted_by", deleted_by)
	err := h.Service.Delete(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in deleting user"))
		return
	}
	h.Logger.Info(logger.User, logger.Delete, "user deleted successfuly", nil)
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, nil, "user deleted successfuly"))
}