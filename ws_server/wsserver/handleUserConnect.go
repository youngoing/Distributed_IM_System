package wsserver

import "github.com/sirupsen/logrus"

// 添加用户的写入通道到映射和添加用户的写入通道到映射和userList
func (server *WebSocketServerStruct) addUserToChannel(userId string, channel chan []byte) {
	server.userManager.mu.RLock() // 使用读锁
	defer server.userManager.mu.RUnlock()

	// 检查是否已存在通道
	if _, exists := server.userManager.userChannels[userId]; exists {
		logrus.Infof("User %s already has a channel", userId)
		return
	}
	logrus.Infof("Adding user with ID: %v", userId)

	// 添加用户通道
	server.userManager.userChannels[userId] = channel
	server.userManager.userList[userId] = server.nodeId
	logrus.Infof("Write channel for user %s has been added", userId)
	//查看userList和userChannels
	logrus.Infof("userList: %v", server.userManager.userList)
	//查看userChannels的key和value
	for k, v := range server.userManager.userChannels {
		logrus.Infof("key: %s, value: %v", k, v)
	}
}

// 从映射中移除用户的写入通道，从userList删除user
func (server *WebSocketServerStruct) removeUserFromChannel(userId string) {
	server.userManager.mu.RLock() // 使用读锁
	defer server.userManager.mu.RUnlock()
	delete(server.userManager.userList, userId)
	delete(server.userManager.userChannels, userId)
	err := removeUserFromRedis(userId)
	if err != nil {
		logrus.Error("Failed to remove user from Redis:", err)
	}

	logrus.Infof("Write channel for user %s has been removed", userId)
}
