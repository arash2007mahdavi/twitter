package logger

import (
	"twitter/src/configs"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Zaplogger struct {
	cfg *configs.Config
	logger *zap.SugaredLogger
}

func NewZaplogger(cfg *configs.Config) *Zaplogger {
	var zaplogger Zaplogger
	zaplogger.cfg = cfg
	zaplogger.Init()
	return &zaplogger
}

var ZaploggerLevels = map[string]zapcore.Level{
	"Debug": zapcore.DebugLevel,
	"Info": zapcore.InfoLevel,
	"Warn": zapcore.WarnLevel,
	"Error": zapcore.ErrorLevel,
	"Fatal": zapcore.FatalLevel,
}

func getLogLevelzap(cfg *configs.Config) zapcore.Level {
	level, ok := ZaploggerLevels[cfg.Logger.Level]
	if ok {
		return level
	}
	return zapcore.InfoLevel
}

func (z *Zaplogger) Init() {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename: "./logs.log",
		MaxSize: 50,
		MaxAge: 300,
		MaxBackups: 100,
		Compress: true,
		LocalTime: true,
	})

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		w,
		getLogLevelzap(z.cfg),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel)).Sugar()
	logger = logger.With("AppName", "Twitter", "Logger", "Zaplogger")
	z.logger = logger
}

func transformExtraToParams(extra map[ExtraCategory]interface{}) []interface{} {
	params := []interface{}{}
	for key, value := range extra {
		params = append(params, string(key))
		params = append(params, value)
	}
	return params
}

func (z *Zaplogger) Debug(cat Category, sub SubCategory, msg string, extra map[ExtraCategory]interface{}) {
	params := transformExtraToParams(extra)
	z.logger.Debugw(msg, params...)
}
func (z *Zaplogger) Debugf(template string, args ...interface{}) {
	z.logger.Debugf(template, args...)
}

func (z *Zaplogger) Info(cat Category, sub SubCategory, msg string, extra map[ExtraCategory]interface{}) {
	params := transformExtraToParams(extra)
	z.logger.Infow(msg, params...)
}
func (z *Zaplogger) Infof(template string, args ...interface{}) {
	z.logger.Infof(template, args...)
}

func (z *Zaplogger) Warn(cat Category, sub SubCategory, msg string, extra map[ExtraCategory]interface{}) {
	params := transformExtraToParams(extra)
	z.logger.Warnw(msg, params...)
}
func (z *Zaplogger) Warnf(template string, args ...interface{}) {
	z.logger.Warnf(template, args...)
}

func (z *Zaplogger) Error(cat Category, sub SubCategory, msg string, extra map[ExtraCategory]interface{}) {
	params := transformExtraToParams(extra)
	z.logger.Errorw(msg, params...)
}
func (z *Zaplogger) Errorf(template string, args ...interface{}) {
	z.logger.Errorf(template, args...)
}

func (z *Zaplogger) Fatal(cat Category, sub SubCategory, msg string, extra map[ExtraCategory]interface{}) {
	params := transformExtraToParams(extra)
	z.logger.Fatalw(msg, params...)
}
func (z *Zaplogger) Fatalf(template string, args ...interface{}) {
	z.logger.Fatalf(template, args...)
}