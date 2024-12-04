package main

import (
	"msg_queue_server/mq"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	// 初始化消息队列
	mq.InitMq()
	run()
}

func run() {
	// 捕获终止信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 创建数据库和消息队列的 Handler
	dbHandler, err := mq.NewHandler()
	if err != nil {
		logrus.Fatalf("Failed to create database and message queue handler: %v", err)
	}
	defer dbHandler.Close()

	// 创建一个通道用于控制退出信号
	stopChan := make(chan struct{})

	// 启动消息队列处理和数据库处理
	go func() {
		mq.RunMq()
	}()

	go func() {
		dbHandler.RunDbHandler()
	}()

	// 等待终止信号
	sig := <-sigChan
	logrus.Infof("Received signal: %v, try shutting down...", sig)

	// 向 stopChan 发送停止信号
	close(stopChan)

	// 等待所有的服务关闭
	dbHandler.Close()

	logrus.Info("All services have been stopped")
}
