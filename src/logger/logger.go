package logger

import (
	"twitter/src/configs"
)

type Logger interface {
	Init()

	Debug(cat Category, sub SubCategory, msg string, extra map[ExtraCategory]interface{})
	Debugf(template string, args ...interface{})

	Info(cat Category, sub SubCategory, msg string, extra map[ExtraCategory]interface{})
	Infof(template string, args ...interface{})

	Warn(cat Category, sub SubCategory, msg string, extra map[ExtraCategory]interface{})
	Warnf(template string, args ...interface{})

	Error(cat Category, sub SubCategory, msg string, extra map[ExtraCategory]interface{})
	Errorf(template string, args ...interface{})

	Fatal(cat Category, sub SubCategory, msg string, extra map[ExtraCategory]interface{})
	Fatalf(template string, args ...interface{})
}

func NewLogger() Logger {
	cfg := configs.GetConfig()
	if cfg.Logger.Type == "Zaplogger" {
		return NewZaplogger(cfg)
	} else {
		return NewZerologger(cfg)
	}
}