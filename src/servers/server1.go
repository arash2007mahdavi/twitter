package servers

import (
	"fmt"
	"twitter/src/configs"
	"twitter/src/logger"
	"twitter/src/routers"
	"twitter/src/services"
	"twitter/src/validations"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var log = logger.NewLogger()

func Init_Server(cfg *configs.Config) {
	engine := gin.New()
	engine.Use(gin.Recovery(), gin.Logger())

	val, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		val.RegisterValidation("mobile", validations.ValidateMobileNumber, true)
		val.RegisterValidation("password", validations.ValidatePassword, true)
		val.RegisterValidation("username", validations.ValidateUsername, true)
	}

	twitter := engine.Group("/twitter")
	{
		user := twitter.Group("/user")
		routers.UserRouter(user)
		tweet := twitter.Group("/tweet")
		routers.TweetRouter(tweet)
	}
	log.Info(logger.Server, logger.Start, "started successfuly", nil)
	engine.Run(fmt.Sprintf("%v:%v", cfg.Server.Host, cfg.Server.Port))
}