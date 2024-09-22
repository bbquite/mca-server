package utils

import "go.uber.org/zap"

func InitLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	sugar := logger.Sugar()
	defer logger.Sync()

	return sugar, nil
}
