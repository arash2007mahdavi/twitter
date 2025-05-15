package handlers

import (
	"fmt"
	"net/http"
	"twitter/src/database"
	"twitter/src/database/models"
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
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, fmt.Errorf("enter otp"), "otp didnt enter"))
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

func (h *UserHelper) Follow(ctx *gin.Context) {
	db := database.GetDB()
	username := ctx.Query("username")
	password := ctx.Query("password")
	target_username := ctx.Query("target_username")
	if len(username)==0 || len(password)==0 || len(target_username)==0 {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, fmt.Errorf("error in query"), "enter username, password and target_username"))
		return
	}
	if username == target_username {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, fmt.Errorf("error in request"), "invalid query"))
		return
	}
	user := models.User{}
	tx := db.WithContext(ctx).Begin()
	err := tx.Model(&models.User{}).Where("username = ? AND deleted_by is null", username).First(&user).Error
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
	target_user := models.User{}
	err = tx.Model(&models.User{}).Where("username = ? AND deleted_by is null", target_username).First(&target_user).Error
	if err != nil {
		tx.Rollback()
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "invalid target user"))
		return
	}
	try_sample := models.UserFollowers{}
	err = tx.Model(&models.UserFollowers{}).Where("user_id = ? AND follower_id = ? AND deleted_at is null", target_user.Id, user.Id).First(&try_sample).Error
	if err == nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, fmt.Errorf("already following"), "user is already a follower"))
		tx.Rollback()
		return
	}
	user_follower := models.UserFollowers{
		UserId: target_user.Id,
		FollowerId: user.Id,
	}
	err = tx.Model(&models.UserFollowers{}).Create(&user_follower).Error
	if err != nil {
		tx.Rollback()
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in add follower"))
		return
	}
	tx.Commit()
	h.Logger.Info(logger.User, logger.Follow, "user followed other one", nil)
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, nil, "followed successfuly"))
}

func (h *UserHelper) GetFollowers(ctx *gin.Context) {
	followers, err := h.Service.GetFollowers(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusAccepted, responses.GenerateResponseWithError(http.StatusAccepted, err, "error in getting followers"))
		return
	}
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, followers, "followers list"))
}

func (h *UserHelper) GetFollowings(ctx *gin.Context) {
	followings, err := h.Service.GetFollowings(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in get followings"))
		return
	}
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, followings, "followings get"))
}