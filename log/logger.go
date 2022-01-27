package log

import (
	"github.com/fighterlyt/log"
	"go.uber.org/zap"
)

func Info(logger log.Logger, msg string, fields ...zap.Field) {
	if logger == nil {
		return
	}

	logger.Info(msg, fields...)
}

func Warn(logger log.Logger, msg string, fields ...zap.Field) {
	if logger == nil {
		return
	}

	logger.Warn(msg, fields...)
}

func Debug(logger log.Logger, msg string, fields ...zap.Field) {
	if logger == nil {
		return
	}

	logger.Debug(msg, fields...)
}

/*Error 错误日志，之所以不用error,和类型名冲突
参数:
*	logger	log.Logger  	参数1
*	msg   	string      	参数2
*	fields	...zap.Field	参数3
返回值:
*/
func Error(logger log.Logger, msg string, fields ...zap.Field) {
	if logger == nil {
		return
	}

	logger.Error(msg, fields...)
}
