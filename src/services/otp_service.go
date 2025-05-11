package services

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"twitter/src/cache"
	"twitter/src/logger"

	"github.com/redis/go-redis/v9"
)

var log = logger.NewLogger()

type OtpService struct {
	Redis *redis.Client
	Logger logger.Logger
}

func NewOtpService() *OtpService {
	return &OtpService{
		Redis: cache.GetRedis(), Logger: logger.NewLogger(),
	}
}

func MakeOtp() string {
	rand.Seed(time.Now().UnixMilli())
	min := 100000
	max := 999999
	otp := rand.Intn(max - min) + min
	string_otp := strconv.Itoa(otp)
	return string_otp
}

func (s *OtpService) SetOtp(mobileNumber string, otp string) error {
	redisValue := cache.RedisValue{Value: otp, Valid: true}
	err := cache.Set(s.Redis, mobileNumber, &redisValue, 2)
	if err != nil {
		return err
	}
	log.Info(logger.Otp, logger.Set, "new otp set", map[logger.ExtraCategory]interface{}{logger.MobileNumber: mobileNumber})
	return nil
}

func (s *OtpService) ValidateOtp(mobileNumber string, test_otp string) error {
	res, err := cache.Get(s.Redis, mobileNumber)
	if err != nil {
		return fmt.Errorf("this mobileNumber doesnt exists")
	} else if !res.Valid {
		return fmt.Errorf("otp used")
	} else if res.Value != test_otp {
		return fmt.Errorf("invalid otp")
	}
	log.Info(logger.Otp, logger.Validate, "otp validated successfuly", map[logger.ExtraCategory]interface{}{logger.MobileNumber: mobileNumber})
	return nil
}