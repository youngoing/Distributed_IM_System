// components/WebSocketService.js
import {websocketUrl} from '../api_list';
class WebSocketService {
  constructor() {
      if (!WebSocketService.instance) {
          this.ws = null;
          this.sendMessageQueue = []; // 消息队列，存储未发送的消息
          this.listeners = []; // 消息监听器
          this.isLoggedIn = false; // 用户登录状态
          WebSocketService.instance = this;
      }
      return WebSocketService.instance;
  }

  // 连接 WebSocket
  connect(user_detail_id) {
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
          console.log('WebSocket 已连接');
          return; // 如果 WebSocket 已连接，直接返回
      }
      // ws://127.0.0.1:3000/_next/webpack-hmr
    //   const url = websocketUrl(user_detail_id);

    const url = `ws://172.25.59.171:8000?token=${user_detail_id}`;
      console.log('Connecting to WebSocket:', url);
      this.ws = new WebSocket(url);

      this.ws.onopen = () => {
          console.log('WebSocket 已连接');
          this.isLoggedIn = true; // 登录成功时，设置为已登录
          this.sendMessageQueue.forEach(msg => {
              this.send(msg); // 发送队列中的消息
          });
          this.sendMessageQueue = []; // 清空消息队列
          
      };
      // 监听 WebSocket 消息
      this.ws.onmessage = (event) => {
          const message = JSON.parse(event.data);
          // 调用所有注册的监听器
          this.listeners.forEach(listener => {
              listener(message);
          });

          // 保存消息到 sessionStorage
          this.saveMessageToLocalStorage(message);
      };

      this.ws.onclose = () => {
          console.log('WebSocket 连接关闭');
          this.isLoggedIn = false; // 用户退出时，设置为未登录
      };

      this.ws.onerror = (err) => {
          console.error('WebSocket 错误:', err);
      };
  }

  // 发送消息
  sendMessage(msg) {
      const msgString = JSON.stringify(msg);

      //存储到localStorage
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
          this.send(msgString); // WebSocket 连接打开时，直接发送消息
            // 保存自己消息到 localStorage
            this.saveUserMessageToLocalStorage(msgString);
      } else {
          console.warn('WebSocket 未连接，消息已存储:', msg);
          this.storeMessage(msgString); // 存储未发送的消息
      }
  }

  // 存储消息到队列
  storeMessage(msg) {
      this.sendMessageQueue.push(msg);
  }
  // 实际发送消息的方法
  send(msg) {
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        //    console.log('Sending message:', msg);
          this.ws.send(msg);
      } else {
          console.error('WebSocket not open, cannot send message');
      }
  }

  // 注册消息监听器
  onMessage(listener) {
      this.listeners.push(listener);
    //   console.log('Message listener registered', listener);
  }

  // 关闭 WebSocket 连接（用户退出时调用）
  disconnect() {
      if (this.ws) {
          this.ws.close();
          this.ws = null; // 清除 WebSocket 实例
          this.isLoggedIn = false; // 设置为未登录
      }
  }

  // 保存接收到的消息到 localStorage
saveMessageToLocalStorage(message) {
    try {
        const messageObj = typeof message === 'string' ? JSON.parse(message) : message;

        let storageKey;
        let chatId;

        switch (messageObj.msg_type) {
            case 'group':
                storageKey = 'group_msgs';
                chatId = messageObj.group_id;
                break;
            case 'private':
                storageKey = 'private_msgs';
                chatId = messageObj.sender_id;
                break;
            case 'invition':
                storageKey = 'notice_msgs';
                chatId = 'invition';
                break;
            case 'system':
                storageKey = 'notice_msgs';
                chatId = 'system';
                // 对于通知类消息，直接使用数组存储
                break; // 不执行后续的对象存储逻辑
            default:
                console.warn('Unknown message type:', messageObj.msg_type);
                return;
        }

        // 以下代码只处理 group、private、invition 和 system 消息
        const storedMessages = localStorage.getItem(storageKey);
        const messages = storedMessages ? JSON.parse(storedMessages) : {};

        // 如果没有该 chatId 的消息记录，初始化为空数组
        if (!messages[chatId]) {
            messages[chatId] = [];
        }

        // 将新的消息添加到对应的聊天记录中
        messages[chatId].push(messageObj);

        // 更新 localStorage 中的消息
        localStorage.setItem(storageKey, JSON.stringify(messages));

    } catch (error) {
        console.error('Failed to save message to localStorage:', error);
    }
}
    saveUserMessageToLocalStorage(message) {
        try {
            console.log("new websocketmessage",message);
            // 如果 message 是对象，直接使用；如果是 JSON 字符串，解析为对象
            const messageObj = typeof message === 'string' ? JSON.parse(message) : message;

            // 根据消息类型选择存储的键
            let storageKey;
            switch(messageObj.msg_type) {
                case 'group':
                    storageKey = 'group_msgs';
                    break;
                case 'private':
                    storageKey = 'private_msgs';
                    break;
                // case 'invition':
                //     storageKey = 'invition_msgs';
                //     break;
                // case 'system':
                //     storageKey = 'system_msgs';
                //     break;
                case 'invition':
                    storageKey = 'notice_msgs';
                    break;
                default:
                    storageKey = 'notice';
            }

            // 获取聊天ID
            let chatId;
            if (messageObj.msg_type === 'group') {
                chatId = messageObj.group_id;
            } else if (messageObj.msg_type === 'private') {
                chatId = messageObj.receiver_id[0];
            } else if (messageObj.msg_type === 'invition' || messageObj.msg_type === 'system') {
                chatId = 'notice'; // 系统消息和邀请消息统一使用 'notice' 作为 chatId
            }

            // 获取现有的存储消息
            const storedMessages = localStorage.getItem(storageKey);
            const messages = storedMessages ? JSON.parse(storedMessages) : {};

            // 如果没有该 chatId 的消息记录，初始化为空数组
            if (!messages[chatId]) {
                messages[chatId] = [];
            }

            // 将新的消息添加到对应的聊天记录中
            messages[chatId].push(messageObj);

            // 更新 localStorage 中的消息
            localStorage.setItem(storageKey, JSON.stringify(messages));

        } catch (error) {
            console.error('Failed to save message to localStorage:', error);
        }
    }
    removeMessageHandler(handler) {
        this.listeners = this.listeners.filter(listener => listener !== handler);
        console.log('Message listener removed', handler);
    }




}

const instance = new WebSocketService();
export default instance;
