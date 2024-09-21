package utils

import (
	"crispy-garbanzo/global"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap/zapcore"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetWriteSyncer(file string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file, //日志文件的位置
		MaxSize:    10,   //在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: 200,  //保留旧文件的最大个数
		MaxAge:     30,   //保留旧文件的最大天数
		Compress:   true, //是否压缩/归档旧文件
	}

	if global.FPG_CONFIG.Zap.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	}
	return zapcore.AddSync(lumberJackLogger)
}

// 获取当天/当周【周日开始】/当月的时间范围
func GetTimeRange(rangeType int) (time.Time, time.Time) {
	now := time.Now()

	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	if rangeType == 1 {
		return startOfDay, endOfDay
	}
	startOfWeek := now.Truncate(24*time.Hour).AddDate(0, 0, -int(now.Weekday())+int(time.Sunday))
	// 将时间设置为 0 点
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())
	// 结束时间为下周六的 23:59:59.999
	endOfWeek := startOfWeek.AddDate(0, 0, 6).Add(24*time.Hour - time.Nanosecond)

	if rangeType == 2 {
		return startOfWeek, endOfWeek
	}

	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	if rangeType == 3 {
		return startOfMonth, endOfMonth
	}
	return startOfDay, startOfDay
}
