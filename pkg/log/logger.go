package log

import (
	"center.bojiu.com/config"
	"go.uber.org/zap"

	"strings"

	clog "common.bojiu.com/log"
)

var (
	ZapLog *zap.Logger
)

func InitLogger() *zap.Logger {
	log := config.GlobalCfg.Log
	level := log.Level
	logLevel := zap.DebugLevel
	if strings.EqualFold("debug", level) {
		logLevel = zap.DebugLevel
	}
	if strings.EqualFold("info", level) {
		logLevel = zap.InfoLevel
	}
	if strings.EqualFold("error", level) {
		logLevel = zap.ErrorLevel
	}
	if strings.EqualFold("warn", level) {
		logLevel = zap.WarnLevel
	}
	return clog.NewLogger(
		clog.SetPath(log.Path),
		clog.SetPrefix(log.Prefix),
		clog.SetDevelopment(log.Development),
		clog.SetDebugFileSuffix(log.DebugFileSuffix),
		clog.SetWarnFileSuffix(log.WarnFileSuffix),
		clog.SetErrorFileSuffix(log.ErrorFileSuffix),
		clog.SetInfoFileSuffix(log.InfoFileSuffix),
		clog.SetMaxAge(log.MaxAge),
		clog.SetMaxBackups(log.MaxBackups),
		clog.SetMaxSize(log.MaxSize),
		clog.SetLevel(logLevel),
	)
}
