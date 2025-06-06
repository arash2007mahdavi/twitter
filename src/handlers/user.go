package handlers

import (
	"fmt"
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

// GetOtp godoc
// @Summary Get Otp
// @Description get otp by valid mobile_number
// @Tags User
// @Accept json
// @Produce json
// @Param OtpDto body OtpDto true "body for get mobile_number"
// @Success 200 {object} responses.Response{result=string} "Success"
// @Failure 400 {object} responses.Response{} "Validation Error"
// @Failure 503 {object} responses.Response{} "Redis Error"
// @Router /user/get/otp [post]
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
		ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, responses.GenerateResponseWithError(http.StatusServiceUnavailable, err, "error in seting otp"))
		return
	}
	h.Logger.Info(logger.Otp, logger.Set, "new otp set", map[logger.ExtraCategory]interface{}{logger.MobileNumber: req.MobileNumber})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, map[string]string{"otp": otp}, "otp set successfuly"))
}

// NewUser godoc
// @Summary Make New User
// @Description make new user with username and password and mobile_number and also the otp has gotten by mobile_number
// @Tags User
// @Accept json
// @Produce json
// @Param UserCreate body dtos.UserCreate true "body for create user"
// @Param otp query string true "the otp has been gotten for the mobile_number"
// @Success 200 {object} responses.Response{result=dtos.UserResponse} "Success"
// @Failure 400 {object} responses.Response{} "Validation Error"
// @Failure 406 {object} responses.Response{} "Error"
// @Router /user/new [post]
func (h *UserHelper) NewUser(ctx *gin.Context) {
	req := dtos.UserCreate{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, responses.GenerateResponseWithValidationError(http.StatusBadRequest, err, "validation error"))
		return
	}
	HashPassword, err:= bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in hashing password"))
		return
	}
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

// GetUsers godoc
// @Summary Get Users
// @Description get users just with admin username and password
// @Tags User
// @Produce json
// @Param admin_username header string true "get admin username"
// @Param admin_password header string true "get admin password"
// @Success 200 {object} responses.Response{result=[]dtos.UserResponse} "Success"
// @Failure 406 {object} responses.Response{} "Error"
// @Router /user/get/users [get]
func (h *UserHelper) GetUsers(ctx *gin.Context) {
	users, err := h.Service.GetUsers(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in get users"))
		return
	}
	h.Logger.Info(logger.Admin, logger.See, "admin saw users", map[logger.ExtraCategory]interface{}{logger.Adminname: ctx.Value("admin_username")})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, users, "users list"))
}

// GetProfile godoc
// @Summary Get Profile
// @Description get user profile
// @Tags User
// @Produce json
// @Param username query string true "get user's username"
// @Success 200 {object} responses.Response{result=models.User} "Success"
// @Failure 500 {object} responses.Response{} "Internal Server Error"
// @Router /user/get/profile [get]
func (h *UserHelper) GetProfile(ctx *gin.Context) {
	user, err := h.Service.GetProfile(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, responses.GenerateResponseWithError(http.StatusInternalServerError, err, "error in get user"))
		return
	}
	h.Logger.Info(logger.User, logger.Profile, "get profile", map[logger.ExtraCategory]interface{}{logger.Userid: ctx.Value("user_id")})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, user, "your profile"))
}

// UpdateUser godoc
// @Summary Update User
// @Description update user's firstname, lastname and ...
// @Tags User
// @Accept json
// @Produce json
// @Param UserUpdate body dtos.UserUpdate true "body for update user"
// @Param username query string true "get user's username"
// @Param password query string true "get user's password"
// @Param modified_by query string false "get modified_by user id"
// @Success 200 {object} responses.Response{result=dtos.UserResponse} "Success"
// @Failure 400 {object} responses.Response{} "Validation Error"
// @Failure 500 {object} responses.Response{} "Internal Server Error"
// @Router /user/update [put]
func (h *UserHelper) UpdateUser(ctx *gin.Context) {
	req := dtos.UserUpdate{}
	req.Enabled = true
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, responses.GenerateResponseWithValidationError(http.StatusBadRequest, err, ""))
		return
	}
	if req.Password != "" {
		hashedPassword, _:= bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		req.Password = string(hashedPassword)
	}
	res, err := h.Service.Update(ctx, &req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, responses.GenerateResponseWithError(http.StatusInternalServerError, err, "error in updating user"))
		return
	}
	h.Logger.Info(logger.User, logger.Update, "user updated successfuly", map[logger.ExtraCategory]interface{}{logger.Userid: ctx.Value("user_id")})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, res, "user updated successfuly"))
}

// DeleteUser godoc
// @Summary Delete User
// @Description delete user with user's username and password
// @Tags User
// @Produce json
// @Param username query string true "get user's username"
// @Param password query string true "get user's password"
// @Param deleted_by query string false "get deleted_by user id"
// @Success 200 {object} responses.Response{result=string} "Success"
// @Failure 500 {object} responses.Response{} "Internal Server Error"
// @Router /user/delete [delete]
func (h *UserHelper) DeleteUser(ctx *gin.Context) {
	err := h.Service.Delete(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, responses.GenerateResponseWithError(http.StatusInternalServerError, err, "error in deleting user"))
		return
	}
	h.Logger.Info(logger.User, logger.Delete, "user deleted successfuly", map[logger.ExtraCategory]interface{}{logger.Userid: ctx.Value("user_id")})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, "deleted", "user deleted successfuly"))
}

// Follow godoc
// @Summary Follow An User
// @Description following an user by another with it's username and password and also the target user's id
// @Tags User
// @Produce json
// @Param username query string true "get user's username"
// @Param password query string true "get user's password"
// @Param target_username query string false "get target's username"
// @Success 200 {object} responses.Response{result=string} "Success"
// @Failure 500 {object} responses.Response{} "Internal Server Error"
// @Router /user/follow [post]
func (h *UserHelper) Follow(ctx *gin.Context) {
	err := h.Service.Follow(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, responses.GenerateResponseWithError(http.StatusInternalServerError, err, "error in follow"))
		return
	}
	h.Logger.Info(logger.User, logger.Follow, "user followed other one", map[logger.ExtraCategory]interface{}{logger.Userid: ctx.Value("user_id"), logger.Targetid: ctx.Value("target_id")})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, "followed", "followed successfuly"))
}

// UnFollow godoc
// @Summary UnFollow An User
// @Description unfollowing an user by another with it's username and password and also the target user's id
// @Tags User
// @Produce json
// @Param username query string true "get user's username"
// @Param password query string true "get user's password"
// @Param target_username query string false "get target's username"
// @Success 200 {object} responses.Response{result=string} "Success"
// @Failure 500 {object} responses.Response{} "Internal Server Error"
// @Router /user/unfollow [delete]
func (h *UserHelper) UnFollow(ctx *gin.Context) {
	err := h.Service.UnFollow(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in unfollowing target"))
		return
	}
	h.Logger.Info(logger.User, logger.UnFollow, "user unfollowed other one", map[logger.ExtraCategory]interface{}{logger.Userid: ctx.Value("user_id"), logger.Targetid: ctx.Value("target_id")})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, "unfollowed", "target unfollowed successfuly"))
}

// GetFollowers godoc
// @Summary Get Followers
// @Description get followers of an user
// @Tags User
// @Produce json
// @Param username query string true "get user's username"
// @Success 200 {object} responses.Response{result=[]dtos.UserResponse} "Success"
// @Failure 500 {object} responses.Response{} "Internal Server Error"
// @Router /user/get/followers [get]
func (h *UserHelper) GetFollowers(ctx *gin.Context) {
	followers, err := h.Service.GetFollowers(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, responses.GenerateResponseWithError(http.StatusInternalServerError, err, "error in getting followers"))
		return
	}
	h.Logger.Info(logger.User, logger.Follower, "get followers", map[logger.ExtraCategory]interface{}{logger.Userid: ctx.Value("user_id")})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, followers, "followers list"))
}

// GetFollowings godoc
// @Summary Get Followings
// @Description get followings of an user
// @Tags User
// @Produce json
// @Param username query string true "get user's username"
// @Success 200 {object} responses.Response{result=[]dtos.UserResponse} "Success"
// @Failure 500 {object} responses.Response{} "Internal Server Error"
// @Router /user/get/followings [get]
func (h *UserHelper) GetFollowings(ctx *gin.Context) {
	followings, err := h.Service.GetFollowings(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, responses.GenerateResponseWithError(http.StatusNotAcceptable, err, "error in get followings"))
		return
	}
	h.Logger.Info(logger.User, logger.Following, "get followings", map[logger.ExtraCategory]interface{}{logger.Userid: ctx.Value("user_id")})
	ctx.JSON(http.StatusOK, responses.GenerateNormalResponse(http.StatusOK, followings, "followings get"))
}