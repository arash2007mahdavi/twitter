package servers

import (
	"fmt"
	"twitter/src/configs"
	"twitter/src/docs"
	"twitter/src/logger"
	"twitter/src/metrics"
	"twitter/src/middlewares"
	"twitter/src/routers"
	"twitter/src/validations"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var log = logger.NewLogger()

func Init_Server(cfg *configs.Config) {
	engine := gin.New()
	engine.Use(gin.Recovery(), gin.Logger())
	engine.Use(middlewares.Prometheus)

	val, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		val.RegisterValidation("mobile", validations.ValidateMobileNumber, true)
		val.RegisterValidation("password", validations.ValidatePassword, true)
		val.RegisterValidation("username", validations.ValidateUsername, true)
	}

	err := prometheus.Register(metrics.DbCalls)
	if err != nil {
		log.Error(logger.Prometheus, logger.Start, "failed in start", nil)
	}
	err = prometheus.Register(metrics.HttpDuration)
	if err != nil {
		log.Error(logger.Prometheus, logger.Start, "failed in start", nil)
	}

	twitter := engine.Group("/twitter")
	{
		user := twitter.Group("/user")
		routers.UserRouter(user)
		tweet := twitter.Group("/tweet")
		routers.TweetRouter(tweet)
		comment := twitter.Group("/comment")
		routers.CommentRouter(comment)
		file := twitter.Group("/file")
		routers.FileRouter(file)

		twitter.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}

	RegisterSwagger(engine)

	log.Info(logger.Server, logger.Start, "started successfuly", nil)
	engine.Run(fmt.Sprintf("%v:%v", cfg.Server.Host, cfg.Server.Port))
}

func RegisterSwagger(r *gin.Engine) {
	docs.SwaggerInfo.Title = "twitter"
	docs.SwaggerInfo.Description = "twitter"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/twitter"
	docs.SwaggerInfo.Host = "localhost:2025"
	docs.SwaggerInfo.Schemes = []string{"http"}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}