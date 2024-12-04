package main

import (
	"ws_server/wsserver"

	"github.com/sirupsen/logrus"
)

func main() {
	// 调用确认环境变量函数
	if err := wsserver.ConfirmEnv(); err != nil {
		logrus.Infof("Environment variable check failed: %v", err)
	}

	wsserver.StartWsServer()
}
