package configs

import "time"

type Config struct {
	Server   ConfigServer
	Logger   ConfigLogger
	Postgres ConfigPostgres
	Redis    ConfigRedis
	Jwt      ConfigJWT
}

type ConfigServer struct {
	Host string
	Port int
}

type ConfigLogger struct {
	Type  string
	Level string
}

type ConfigPostgres struct {
	Host            string
	Port            int
	User            string
	Password        string
	Dbname          string
	Sslmode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifeTime time.Duration
}

type ConfigRedis struct {
	WriteTimeOut time.Duration
	ReadTimeOut  time.Duration
	DialTimeOut  time.Duration
	Host         string
	Port         int
	DB           int
	Password     string
	Poolsize     int
	PoolTimeOut  time.Duration
}

type ConfigJWT struct {
	Secret                     string
	RefreshSecret              string
	AccessTokenExpireDuration  time.Duration
	RefreshTokenExpireDuration time.Duration
}
