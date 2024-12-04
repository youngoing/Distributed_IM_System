import { useState, useEffect } from 'react';
import axios from 'axios';
import { isClientAuthenticated } from '../utils/auth';
import { login, register,websocketUrl } from '../api_list';
export default function AuthForm() {
  const [email, setEmail] = useState('');
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const [isLogin, setIsLogin] = useState(true); // 用于切换登录和注册

  useEffect(() => {
    const checkAuthentication = async () => {
      const isAuthenticated = await isClientAuthenticated();
      if (isAuthenticated) {
        // 如果用户已经登录，重定向到首页
        alert('你已经登录了');
        setTimeout(() => {
          window.location.href = '/';
        }, 1000); // 1000 毫秒 = 1 秒
      }
    };

    checkAuthentication();
  }, []);

  const handleAuth = async (e) => {
    e.preventDefault();
  
    if (!isLogin && password !== confirmPassword) {
      setError('两次输入的密码不一致');
      return;
    }
  
    try {
      const endpoint = isLogin ? login : register;
      const data = isLogin ? { username, password } : { username, email, password };
      const response = await axios.post(endpoint, data, {
        withCredentials: true,
      });
  
      // 从响应中提取用户信息
      const { user } = response.data;

      // 只存储用户信息到 sessionStorage
      sessionStorage.setItem('user', JSON.stringify(user));
  
      console.log('用户信息:', user);
  
      // 成功后的提示
      alert('认证成功');
  
      // 只有登录时才进行跳转
      if (isLogin) {
        setTimeout(() => {
          window.location.href = '/';
        }, 1000); // 1 秒延迟跳转
      }
    } catch (err) {
      setError(err.response?.data?.message || '认证失败');
    }
  };
  

  return (
    <div className="flex justify-center items-center min-h-screen bg-cover bg-center" style={{ backgroundImage: "url('/static/images/background.jpg')" }}>
      <div className="flex w-full">
        <div className="w-2/12"></div>
        <div className="w-8/12 bg-gray-800 text-white shadow-lg rounded-lg p-8">
          <h1 className="text-2xl font-bold mb-4">{isLogin ? '登录' : '注册'}页面</h1>
          <form onSubmit={handleAuth} className="space-y-4">
            {!isLogin && (
              <div>
                <input
                  type="email"
                  placeholder="邮箱"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="w-full p-2 border border-gray-600 rounded bg-gray-700 text-white focus:outline-none focus:ring focus:ring-blue-300"
                />
              </div>
            )}
            <div>
              <input
                type="text"
                placeholder="用户名"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                className="w-full p-2 border border-gray-600 rounded bg-gray-700 text-white focus:outline-none focus:ring focus:ring-blue-300"
              />
            </div>
            <div>
              <input
                type="password"
                placeholder="密码"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="w-full p-2 border border-gray-600 rounded bg-gray-700 text-white focus:outline-none focus:ring focus:ring-blue-300"
              />
            </div>
            {!isLogin && (
              <div>
                <input
                  type="password"
                  placeholder="确认密码"
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  className="w-full p-2 border border-gray-600 rounded bg-gray-700 text-white focus:outline-none focus:ring focus:ring-blue-300"
                />
              </div>
            )}
            <button type="submit" className="w-full bg-blue-500 text-white p-2 rounded hover:bg-blue-600">
              {isLogin ? '登录' : '注册'}
            </button>
          </form>
          {error && <p className="text-red-500 mt-4">{error}</p>}
          <button onClick={() => setIsLogin(!isLogin)} className="w-full mt-4 text-blue-400 hover:underline">
            {isLogin ? '切换到注册' : '切换到登录'}
          </button>
        </div>
        <div className="w-2/12"></div>
      </div>
    </div>
  );
}