package handlers

import (
	"net/http"
	"twitter/src/dtos"
	"twitter/src/logger"
	"twitter/src/responses"
	"twitter/src/services"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

type UserHelper struct {
	Logger  logger.Logger
	Service *services.UserService
	Otp     *services.OtpService
}

func GetUserHelper() *UserHelper {
	return &UserHelper{
		Logger:  logger.NewLogger(),
		Service: services.NewUserService(),
		Otp:     services.NewOtpService(),
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
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, map[string]string{"otp": otp}, "otp set successfuly"))
}

func (h *UserHelper) NewUser(ctx *gin.Context) {
	req := dtos.UserCreate{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, responses.GenerateResponseWithValidationError(http.StatusBadRequest, err, "validation error"))
		return
	}
	HashPassword, _:= bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	req.Password = string(HashPassword)
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
	h.Logger.Info(logger.User, logger.New, "new user added", map[logger.ExtraCategory]interface{}{logger.Username: req.Username})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, res, "new user created successfuly"))
}

func (h *UserHelper) GetUsers(ctx *gin.Context) {
	users, err := h.Service.GetUsers(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in get users"))
		return
	}
	h.Logger.Info(logger.Admin, logger.See, "admin saw users", map[logger.ExtraCategory]interface{}{logger.Adminname: ctx.Value("admin_username")})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, users, "users list"))
}

func (h *UserHelper) GetProfile(ctx *gin.Context) {
	user, err := h.Service.GetProfile(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in get user"))
		return
	}
	h.Logger.Info(logger.User, logger.See, "user saw profile", map[logger.ExtraCategory]interface{}{logger.Userid: ctx.Value("user_id")})
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
	res, err := h.Service.Update(ctx, &req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in updating user"))
		return
	}
	h.Logger.Info(logger.User, logger.Update, "user updated successfuly", map[logger.ExtraCategory]interface{}{logger.Userid: ctx.Value("user_id")})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, res, "user updated successfuly"))
}

func (h *UserHelper) DeleteUser(ctx *gin.Context) {
	err := h.Service.Delete(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in deleting user"))
		return
	}
	h.Logger.Info(logger.User, logger.Delete, "user deleted successfuly", map[logger.ExtraCategory]interface{}{logger.Userid: ctx.Value("user_id")})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, nil, "user deleted successfuly"))
}