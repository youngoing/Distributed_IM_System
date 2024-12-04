// utils/auth.js
import axios from 'axios';
import { auth_login } from '../api_list';

export async function isClientAuthenticated() {
  // 确保只有在浏览器环境下访问 sessionStorage
  if (typeof window !== 'undefined') {
    const user = sessionStorage.getItem('user');
  
    if (user) {
      try {
        const res = await axios.get(auth_login, {
          withCredentials: true,
        });
  
        console.log('API 响应数据:', res.data); // 添加调试日志
  
        if (res.status === 200) {
          return true;
        } else {
          return false;
        }
      } catch (err) {
        console.error('验证用户登录失败:', err);
        return false;
      }
    } else {
      return false;
    }
  } else {
    return false;
  }
  
}