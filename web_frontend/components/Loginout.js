import axios from 'axios';
import { logout } from '../api_list';  // 导入登出 API 端点
import WebSocketService from './WebsocketService';
const Loginout = async () => {
  try {
    await axios.get(logout, {
      withCredentials: true, // 确保请求中带上 Cookie
    });

    sessionStorage.removeItem('user');
    WebSocketService.disconnect()
    window.location.href = '/login';
  } catch (error) {
    console.error('退出登录失败:', error);
  }
};

export default Loginout;
