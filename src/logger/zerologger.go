package logger

import (
	"os"
	"twitter/src/configs"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type Zerologger struct {
	cfg *configs.Config
	logger *zerolog.Logger
}

func NewZerologger(cfg *configs.Config) *Zerologger {
	var zerologger Zerologger
	zerologger.cfg = cfg
	zerologger.Init()
	return &zerologger
}

var ZerologLevels = map[string]zerolog.Level{
	"Debug": zerolog.DebugLevel,
	"Info": zerolog.InfoLevel,
	"Warn": zerolog.WarnLevel,
	"Error": zerolog.ErrorLevel,
	"Fatal": zerolog.FatalLevel,
}

func getLogLevelzero(cfg *configs.Config) zerolog.Level {
	level, ok := ZerologLevels[cfg.Logger.Level]
	if ok {
		return level
	}
	return zerolog.InfoLevel
}

func (z *Zerologger) Init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.SetGlobalLevel(getLogLevelzero(z.cfg))

	file, err := os.OpenFile("./logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("problem in open or create logs.log file")
	}

	logger := zerolog.New(file).With().Timestamp().
		Str("AppName", "Twitter").
		Str("Logger", "Zerologger").
		Logger()
	z.logger = &logger
}

func transformExtra(extra map[ExtraCategory]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	for key, value := range extra {
		result[string(key)] = value
	}
	return result
}

func (z *Zerologger) Debug(cat Category, sub SubCategory, msg string, extra map[ExtraCategory]interface{}) {
	transformed_extra := transformExtra(extra)
	z.logger.Debug().Fields(transformed_extra).Msg(msg)
}
func (z *Zerologger) Debugf(template string, args ...interface{}) {
	z.logger.Debug().Msgf(template, args...)
}

func (z *Zerologger) Info(cat Category, sub SubCategory, msg string, extra map[ExtraCategory]interface{}) {
	transformed_extra := transformExtra(extra)
	z.logger.Info().Fields(transformed_extra).Msg(msg)
}
func (z *Zerologger) Infof(template string, args ...interface{}) {
	z.logger.Info().Msgf(template, args...)
}

func (z *Zerologger) Warn(cat Category, sub SubCategory, msg string, extra map[ExtraCategory]interface{}) {
	transformed_extra := transformExtra(extra)
	z.logger.Warn().Fields(transformed_extra).Msg(msg)
}
func (z *Zerologger) Warnf(template string, args ...interface{}) {
	z.logger.Warn().Msgf(template, args...)
}

func (z *Zerologger) Error(cat Category, sub SubCategory, msg string, extra map[ExtraCategory]interface{}) {
	transformed_extra := transformExtra(extra)
	z.logger.Error().Fields(transformed_extra).Msg(msg)
}
func (z *Zerologger) Errorf(template string, args ...interface{}) {
	z.logger.Error().Msgf(template, args...)
}

func (z *Zerologger) Fatal(cat Category, sub SubCategory, msg string, extra map[ExtraCategory]interface{}) {
	transformed_extra := transformExtra(extra)
	z.logger.Fatal().Fields(transformed_extra).Msg(msg)
}
func (z *Zerologger) Fatalf(template string, args ...interface{}) {
	z.logger.Fatal().Msgf(template, args...)
}