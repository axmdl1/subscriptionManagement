package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func Init() error {
	l, err := zap.NewProduction()
	if err != nil {
		return err
	}

	Log = l
	return nil
}

func Sync() {
	_ = Log.Sync()
}
