package services

import (
	"fmt"
	"time"
	"twitter/src/configs"
	"twitter/src/database"
	"twitter/src/database/models"
	"twitter/src/logger"

	"github.com/golang-jwt/jwt"
)

type TokenService struct {
	logger logger.Logger
	cfg    *configs.Config
}

type TokenDto struct {
	Id           int
	Username     string
	Firstname    string
	Lastname     string
	MobileNumber string
}

func NewTokenService(cfg *configs.Config) *TokenService {
	logger := logger.NewLogger()
	return &TokenService{
		cfg:    cfg,
		logger: logger,
	}
}

type TokenDetail struct {
	AccessToken            string `json:"accessToken"`
	RefreshToken           string `json:"refreshToken"`
	AccessTokenExpireTime  int    `json:"accessTokenExpireTime"`
	RefreshTokenExpireTime int    `json:"refreshTokenExpireTime"`
}

func (s *TokenService) GenerateToken(token *TokenDto) (*TokenDetail, error) {
	accessToken := &TokenDetail{}
	accessToken.AccessTokenExpireTime = int(time.Now().Add(s.cfg.Jwt.AccessTokenExpireDuration * time.Minute).Unix())
	accessToken.RefreshTokenExpireTime = int(time.Now().Add(s.cfg.Jwt.RefreshTokenExpireDuration * time.Hour).Unix())

	atc := jwt.MapClaims{}
	atc["id"] = token.Id
	atc["first_name"] = token.Firstname
	atc["last_name"] = token.Lastname
	atc["username"] = token.Username
	atc["mobileNumber"] = token.MobileNumber
	atc["exp"] = accessToken.AccessTokenExpireTime

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atc)
	var err error
	accessToken.AccessToken, err = at.SignedString([]byte(s.cfg.Jwt.Secret))
	if err != nil {
		return nil, err
	}

	rtc := jwt.MapClaims{}
	rtc["id"] = token.Id
	rtc["exp"] = accessToken.RefreshTokenExpireTime

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtc)
	accessToken.RefreshToken, err = rt.SignedString([]byte(s.cfg.Jwt.Secret))
	if err != nil {
		return nil, err
	}

	return accessToken, nil
}

func (s *TokenService) GenerateAccessTokenByRefreshToken(refresh string) (*TokenDetail, error) {
	token, err := jwt.Parse(refresh, func(token *jwt.Token) (interface{}, error) {
        return []byte(s.cfg.Jwt.RefreshSecret), nil
    })
	if err != nil {
		return &TokenDetail{}, err
	}

	claims := token.Claims.(jwt.MapClaims)
	id := claims["id"].(string)

	DB := database.GetDB()
	var user models.User
	err = DB.Where("id = ? AND enable is true", id).First(&user).Error
	if err != nil {
		return &TokenDetail{}, err
	}
	dtos := TokenDto{
		Id: user.Id,
		Username: user.Username,
		Firstname: user.Firstname,
		Lastname: user.Lastname,
		MobileNumber: user.MobileNumber,
	}
	new_token, err := s.GenerateToken(&dtos)
	if err != nil {
		return &TokenDetail{}, err
	}
	return new_token, nil
}

func (s *TokenService) VerifyToken(token string) (*jwt.Token, error) {
	at, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("error in verify token")
		}
		return []byte(s.cfg.Jwt.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	return at, nil
}

func (s *TokenService) GetClaims(token string) (claimMap map[string]interface{}, err error) {
	claimMap = map[string]interface{}{}

	verifyToken, err := s.VerifyToken(token)
	if err != nil {
		return nil, err
	}
	claims, ok := verifyToken.Claims.(jwt.MapClaims)
	if ok && verifyToken.Valid {
		for k, v := range claims {
			claimMap[k] = v
		}
		return claimMap, nil
	}
	return nil, fmt.Errorf("claims not found")
}
