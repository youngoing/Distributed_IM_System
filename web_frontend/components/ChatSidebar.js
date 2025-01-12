import React, { useState, useEffect, useRef, useCallback } from 'react';
import { AiOutlineArrowRight, AiOutlineArrowLeft, AiOutlinePlus, AiOutlineMessage, AiOutlineUserAdd, AiOutlineTeam, AiOutlineUsergroupAdd } from 'react-icons/ai'; // 导入图标
import { FaArrowLeft } from 'react-icons/fa';
import { search_group_or_friend, applicationUrl, invitionUrl, auth_inviteUrl, createGroupUrl } from '../api_list';
import axios from 'axios';
const ChatSidebar = ({ groups, friends, handleChatChange, user, onLogout, isSidebarOpen, toggleSidebar, currentChatId }) => {
  const [showProfileMenu, setShowProfileMenu] = useState(false);
  const [isGroupsCollapsed, setIsGroupsCollapsed] = useState(false);
  const [isFriendsCollapsed, setIsFriendsCollapsed] = useState(false);
  const [showModal, setShowModal] = useState(false);
  const [modalType, setModalType] = useState(''); // 'addFriend' or 'joinGroup'
  const [searchInput, setSearchInput] = useState('');
  const [searchResults, setSearchResults] = useState({ users: [], groups: [] });
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [isNoticeCollapsed, setIsNoticeCollapsed] = useState(false);
  const [notices, setNotices] = useState({ invition: [], system: [] });
  const [showNoticeModal, setShowNoticeModal] = useState(false);
  const [selectedNotice, setSelectedNotice] = useState(null);
  const [showCreateGroupModal, setShowCreateGroupModal] = useState(false);
  const [groupForm, setGroupForm] = useState({
    name: '',
    description: '',
    avatar_url: ''
  });
  // 发送申请函数
  const handleApply = async (targetId, isGroup) => {
    const data = {
      action: isGroup ? 'group' : 'friend',
      sender_id: user.user_detail_id,
      receiver_id: isGroup ? null : targetId,
      group_id: isGroup ? targetId : null
    };

    try {
      const response = await axios.post(applicationUrl, data, {
        withCredentials: true,
        headers: {
          'Content-Type': 'application/json'
        }
      });

      if (response.status === 200) {
        alert(isGroup ? '已发送入群申请' : '已发送好友申请');
      }
    } catch (error) {
      console.error('申请发送失败:', error);
      alert('申请发送失败，请稍后重试');
    }
  };
  const handleProfileClick = () => {
    setShowProfileMenu((prevState) => !prevState);
  };

  const toggleGroupsCollapse = () => {
    setIsGroupsCollapsed((prevState) => !prevState);
  };

  const toggleFriendsCollapse = () => {
    setIsFriendsCollapsed((prevState) => !prevState);
  };
  // 在打开模态框时自动搜索
  useEffect(() => {
    if (modalType) {
      handleSearch();
    }
  }, [modalType]);

  // 修改搜索函数
  const handleSearch = async () => {
    try {
      if (modalType === 'addFriend') {
        const response = await axios.get(search_group_or_friend(searchInput || '', 'user'), {
          withCredentials: true
        });
        const users = Array.isArray(response.data) ? response.data :
          (response.data.users || []);
        setSearchResults({
          users: users,
          groups: []
        });
        console.log('用户搜索结果:', users);
      } else if (modalType === 'joinGroup') {
        const response = await axios.get(search_group_or_friend(searchInput || '', 'group'), {
          withCredentials: true
        });
        const groups = Array.isArray(response.data) ? response.data :
          (response.data.groups || []);
        setSearchResults({
          users: [],
          groups: groups
        });
        console.log('群组搜索结果:', groups);
      }
    } catch (error) {
      console.error('搜索失败:', error);
      alert('搜索失败，请稍后重试');
    }
  };

  // 处理用户按 Enter 键
  const handleKeyPress = (e) => {
    if (e.key === 'Enter') {
      handleSearch();
    }
  };

  // 搜索框渲染
  const renderSearchBox = () => (
    <div className="relative">
      <input
        type="text"
        value={searchInput}
        onChange={(e) => setSearchInput(e.target.value)}
        onKeyDown={handleKeyPress}
        placeholder={modalType === 'addFriend' ? '搜索用户...' : '搜索群组...'}
        className="w-full bg-gray-800 text-white px-4 py-3 rounded-lg focus:ring-2 focus:ring-blue-500 focus:outline-none"
      />
      <div className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 cursor-pointer" onClick={handleSearch}>
        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
      </div>
    </div>
  );

  // 搜索结果渲染
  const renderSearchResults = () => {
    if (isLoading) {
      return <p className="text-gray-400 text-center py-8">加载中...</p>;
    }

    console.log('当前搜索结果:', searchResults);
    console.log('modalType:', modalType);

    const hasNoResults = modalType === 'addFriend'
      ? (!searchResults.users || searchResults.users.length === 0)
      : (!searchResults.groups || searchResults.groups.length === 0);

    if (hasNoResults) {
      return (
        <div className="text-center text-gray-400 py-8">
          <p className="text-lg">未找到相关{modalType === 'addFriend' ? '用户' : '群组'}</p>
          <p className="text-sm mt-2">试试其他关键词吧</p>
        </div>
      );
    }

    return (
      <div className="max-h-[400px] overflow-y-auto space-y-2 mt-4">
        {/* 根据 modalType 渲染不同的结果 */}
        {modalType === 'addFriend' && searchResults.users.map((user) => (
          <div key={user.user_detail_id} className="bg-gray-800 p-4 rounded-lg hover:bg-gray-700 transition-colors animate-fadeIn">
            <div className="flex items-center space-x-3">
              <img src={user.avatar_url || '/static/images/user.png'} className="w-12 h-12 rounded-full border-2 border-blue-500 object-cover" alt={user.nickname} />
              <div className="flex-1 min-w-0">
                <h3 className="text-white font-medium truncate">{user.nickname}</h3>
                <p className="text-gray-400 text-sm">ID: {user.user_detail_id}</p>
              </div>
              <button
                onClick={() => handleApply(user.user_detail_id, false)}
                className="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white rounded-lg transition-colors flex-shrink-0"
              >
                添加好友
              </button>
            </div>
          </div>
        ))}

        {modalType === 'joinGroup' && searchResults.groups.map((group) => (
          <div key={group.group_id} className="bg-gray-800 p-4 rounded-lg hover:bg-gray-700 transition-colors animate-fadeIn">
            <div className="flex items-center space-x-3">
              <img src={group.avatar_url || '/static/images/group.png'} className="w-12 h-12 rounded-full border-2 border-green-500 object-cover" alt={group.name} />
              <div className="flex-1 min-w-0">
                <h3 className="text-white font-medium truncate">{group.name}</h3>
                <p className="text-gray-400 text-sm">ID: {group.group_id}</p>
              </div>
              <button
                onClick={() => handleApply(group.group_id, true)}
                className="px-4 py-2 bg-green-600 hover:bg-green-500 text-white rounded-lg transition-colors flex-shrink-0"
              >
                加入群组
              </button>
            </div>
          </div>
        ))}
      </div>
    );
  };

  // 添加切换通知栏的函数
  const toggleNoticeCollapse = () => {
    setIsNoticeCollapsed(!isNoticeCollapsed);
  };

  // 处理邀请的函数
  const handleInvitation = async (notice, action) => {
    try {
      // 根据通知类型构建请求数据
      const requestData = {
        msg_id: notice.msg_id,
        token: notice.extra.token,
        action: action, // 'accept' 或 'reject'
        type: notice.extra.invition_type === 'friend' ? 'friend' : 'group',
        sender_id: parseInt(notice.sender_id),
        receiver_id: parseInt(notice.receiver_id[0]),
        group_id: notice.extra.group_id ? parseInt(notice.extra.group_id) : null
      };

      // 发送请求
      const response = await axios.post(auth_inviteUrl, requestData, {
        withCredentials: true,
        headers: {
          'Content-Type': 'application/json'
        }
      });

      // 处理成功响应
      if (response.status === 200) {
        // 从 localStorage 中移除已处理的通知
        const storedNotices = JSON.parse(localStorage.getItem('notice_msgs') || '{}');
        const updatedInvitations = storedNotices.invition.filter(
          msg => msg.msg_id !== notice.msg_id
        );

        // 更新 localStorage
        localStorage.setItem('notice_msgs', JSON.stringify({
          ...storedNotices,
          invition: updatedInvitations
        }));

        // 更新状态
        setNotices(prev => ({
          ...prev,
          invition: updatedInvitations
        }));

        // 关闭弹窗并显示成功消息
        setShowNoticeModal(false);
        alert(response.data.message || (action === 'accept' ? '已接受请求' : '已拒绝请求'));

        // 刷新页面以更新列表
        window.location.reload();
      }
    } catch (error) {
      // 错误处理
      console.error('处理邀请失败:', error);

      if (error.response) {
        // 处理特定的错误状态
        switch (error.response.status) {
          case 400:
            alert(`请求错误: ${error.response.data.details || '参数无效'}`);
            break;
          case 500:
            alert(`服务器错误: ${error.response.data.details || '处理请求失败'}`);
            break;
          default:
            alert('操作失败，请稍后重试');
        }
      } else if (error.request) {
        // 请求发送失败
        alert('网络错误，请检查网络连接');
      } else {
        // 其他错误
        alert('操作失败，请稍后重试');
      }

      // 关闭弹窗
      setShowNoticeModal(false);
    }
  };

  // 添加 useInterval 自定义 Hook
  const useInterval = (callback, delay) => {
    const savedCallback = useRef();

    useEffect(() => {
      savedCallback.current = callback;
    }, [callback]);

    useEffect(() => {
      function tick() {
        savedCallback.current();
      }
      if (delay !== null) {
        let id = setInterval(tick, delay);
        return () => clearInterval(id);
      }
    }, [delay]);
  };

  // 在组件中添加更新通知的函数
  const updateNotices = useCallback(() => {
    const storedNotices = localStorage.getItem('notice_msgs');
    if (storedNotices) {
      const parsedNotices = JSON.parse(storedNotices);
      setNotices(parsedNotices);
    }
  }, []);

  // 使用 useInterval 定时更新通知
  useInterval(() => {
    updateNotices();
  }, 1000); // 每秒更新一次

  // 添加处理通知点击的函数
  const handleNoticeClick = (notice) => {
    setSelectedNotice(notice);
    setShowNoticeModal(true);
  };

  // 添加创建群聊的处理函数
  const handleCreateGroup = async (e) => {
    e.preventDefault();

    try {
      const requestData = {
        user_detail_id: user.user_detail_id,
        name: groupForm.name,
        description: groupForm.description,
        avatar_url: groupForm.avatar_url
      };

      const response = await axios.post(createGroupUrl, requestData, {
        withCredentials: true,
        headers: {
          'Content-Type': 'application/json'
        }
      });

      if (response.status === 200) {
        alert('群聊创建成功');
        setShowCreateGroupModal(false);
        // 重置表单
        setGroupForm({
          name: '',
          description: '',
          avatar_url: ''
        });
        // 刷新页面以更新群聊列表
        window.location.reload();
      }
    } catch (error) {
      console.error('创建群聊失败:', error);
      alert(error.response?.data?.error || '创建群聊失败，请稍后重试');
    }
  };

  return (
    <div className="relative flex">
      <aside className={`transition-all duration-300 ${isSidebarOpen ? 'w-72' : 'w-0'} bg-gray-800 text-white flex flex-col h-screen overflow-hidden`}>
        {/* 个人信息区域 - 优化 */}
        <header className="p-4 bg-gray-900 border-b border-gray-700">
          <div className="flex items-center space-x-4">
            <div className="relative group">
              <img
                src={user?.avatar_url || '/static/images/me.png'}
                alt="用户头像"
                className="w-12 h-12 rounded-full cursor-pointer border-2 border-blue-500 hover:border-blue-400 transition-all"
                onClick={handleProfileClick}
              />
              {showProfileMenu && (
                <div className="absolute top-14 left-0 bg-gray-700 rounded-lg shadow-lg w-48 z-50 py-2">
                  {/* 个人信息预览 */}
                  <div className="px-4 py-2 border-b border-gray-600">
                    <div className="text-sm font-medium">{user?.nickname || '用户'}</div>
                    <div className="text-xs text-gray-400">ID: {user?.user_detail_id}</div>
                  </div>
                  {/* 操作选项 */}
                  <div className="py-1">
                    <button className="w-full px-4 py-2 text-left text-sm hover:bg-gray-600 flex items-center space-x-2">
                      <span>个人设置</span>
                    </button>
                    <button className="w-full px-4 py-2 text-left text-sm hover:bg-gray-600 flex items-center space-x-2">
                      <span>修改头像</span>
                    </button>
                    <button
                      className="w-full px-4 py-2 text-left text-sm hover:bg-gray-600 flex items-center space-x-2 text-red-400"
                      onClick={onLogout}
                    >
                      <span>退出登录</span>
                    </button>
                  </div>
                </div>
              )}
            </div>
            {isSidebarOpen && (
              <div className="flex-1">
                <div className="text-lg font-medium truncate">{user?.nickname || '用户'}</div>
                <div className="text-xs text-gray-400">在线</div>
              </div>
            )}
          </div>
        </header>

        {/* 搜索和添加区域 - 简化布局 */}
        <div className="p-4 bg-gray-750 border-b border-gray-700">
          <div className="flex items-center justify-between">
            <h2 className="text-base font-semibold text-gray-200 flex items-center">
              <AiOutlineMessage className="w-5 h-5 mr-2 text-blue-400" />
              聊天列表
            </h2>
            <button
              onClick={() => setShowModal(true)}
              className="p-2 hover:bg-gray-600 rounded-lg transition-all duration-200 flex items-center space-x-1 text-sm text-gray-300 hover:text-white group"
              title="添加"
            >
              <AiOutlinePlus className="w-5 h-5 group-hover:text-blue-400 transition-colors" />
            </button>
          </div>
        </div>

        {/* 聊天列表区域 - 优化群聊和私聊显示 */}
        <div className="flex-1 overflow-y-auto scrollbar-thin scrollbar-thumb-gray-600">
          {/* 群聊列表 */}
          {groups.length > 0 && (
            <div className="mb-2">

              <div
                className="flex items-center justify-between px-4 py-3 bg-gray-750 cursor-pointer hover:bg-gray-700"
                onClick={toggleGroupsCollapse}
              >
                <h2 className="text-sm font-medium flex items-center space-x-2">
                  <span>群聊列表</span>
                  <span className="text-xs text-gray-400">({groups.length})</span>
                </h2>
                <span className="text-xs transform transition-transform duration-200">
                  {isGroupsCollapsed ? '▶' : '▼'}
                </span>
              </div>
              {!isGroupsCollapsed && (
                <div className="py-1">
                  {groups.map((group) => (
                    <div
                      key={group.group_id}
                      onClick={() => handleChatChange(group.group_id, group.name, true)}
                      className="flex items-center px-4 py-2 hover:bg-gray-700 cursor-pointer group relative"
                    >
                      <img
                        src={group.avatar_url || '/static/images/group.png'}
                        alt={group.name}
                        className="w-10 h-10 rounded-full border-2 border-blue-500"
                      />
                      {isSidebarOpen && (
                        <div className="ml-3 flex-1 min-w-0">
                          <div className="flex items-center justify-between">
                            <span className="text-sm font-medium truncate">{group.name}</span>
                            <span className="text-xs text-gray-400">{group.members?.length + 1 || 1}人</span>
                          </div>
                          <p className="text-xs text-gray-400 truncate">
                            {group.description || '暂无群介绍'}
                          </p>
                        </div>
                      )}
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}

          {/* 好友列表 */}
          {friends.length > 0 && (
            <div className="mb-2">
              <div
                className="flex items-center justify-between px-4 py-3 bg-gray-750 cursor-pointer hover:bg-gray-700"
                onClick={toggleFriendsCollapse}
              >
                <h2 className="text-sm font-medium flex items-center space-x-2">
                  <span>好友列表</span>
                  <span className="text-xs text-gray-400">({friends.length})</span>
                </h2>
                <span className="text-xs transform transition-transform duration-200">
                  {isFriendsCollapsed ? '▶' : '▼'}
                </span>
              </div>
              {!isFriendsCollapsed && (
                <div className="py-1">
                  {friends.map((friend) => (
                    <div
                      key={friend.user_detail_id}
                      onClick={() => handleChatChange(friend.user_detail_id, friend.nickname, false)}
                      className="flex items-center px-4 py-2 hover:bg-gray-700 cursor-pointer group relative"
                    >
                      <img
                        src={friend.avatar_url || '/static/images/user.png'}
                        alt={friend.nickname}
                        className="w-10 h-10 rounded-full border-2 border-green-500"
                      />
                      {isSidebarOpen && (
                        <div className="ml-3 flex-1 min-w-0">
                          <div className="flex items-center justify-between">
                            <span className="text-sm font-medium truncate">{friend.nickname}</span>
                          </div>
                          <p className="text-xs text-gray-400 truncate">
                            {friend.signature || '这个人很懒，什么都没写~'}
                          </p>
                        </div>
                      )}
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}

          {/* 通知消息列表 */}
          {notices.invition?.length > 0 && (
            <div className="mb-2">
              <div
                className="flex items-center justify-between px-4 py-3 bg-gray-750 cursor-pointer hover:bg-gray-700"
                onClick={toggleNoticeCollapse}
              >
                <h2 className="text-sm font-medium flex items-center space-x-2">
                  <span>通知消息</span>
                  <span className="text-xs text-gray-400">({notices.invition.length})</span>
                </h2>
                <span className="text-xs transform transition-transform duration-200">
                  {isNoticeCollapsed ? '▶' : '▼'}
                </span>
              </div>
              {!isNoticeCollapsed && (
                <div className="py-1">
                  {notices.invition.map((notice) => (
                    <div
                      key={notice.msg_id}
                      onClick={() => handleNoticeClick(notice)}
                      className="px-4 py-2 hover:bg-gray-700 cursor-pointer group relative"
                    >
                      <div className="flex items-center space-x-3">
                        <img
                          src={notice.extra.sender_avatar_url || '/static/images/user.png'}
                          alt={notice.extra.sender_nickname}
                          className="w-10 h-10 rounded-full"
                        />
                        <div className="flex-1">
                          <p className="text-sm text-white">{notice.msg_content}</p>
                          <p className="text-xs text-gray-400">
                            {new Date(notice.timestamp / 1000000).toLocaleString()}
                          </p>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}

          {/* 空状态显示 */}
          {groups.length === 0 && friends.length === 0 && (
            <div className="flex flex-col items-center justify-center h-full text-gray-400">
              <p className="text-sm">暂无聊天</p>
              <button
                onClick={() => setShowModal(true)}
                className="mt-4 px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 transition-colors"
              >
                添加好友/群聊
              </button>
            </div>
          )}
        </div>
      </aside>

      {/* 折叠按钮 */}
      <button
        className="absolute top-1/2 -right-4 transform -translate-y-1/2 bg-gray-700 text-white p-2 rounded-full hover:bg-gray-600 transition-colors focus:outline-none shadow-lg"
        onClick={toggleSidebar}
      >
        {isSidebarOpen ? <AiOutlineArrowLeft size={20} /> : <AiOutlineArrowRight size={20} />}
      </button>

      {/* 添加操作选择弹窗 */}
      {showModal && (
        <div className="fixed inset-0 bg-gray-800/70 backdrop-blur-sm flex justify-center items-center z-50">
          <div className="bg-gray-900 p-8 rounded-xl w-96 shadow-2xl border border-gray-700 relative">
            <h2 className="text-2xl font-bold mb-6 text-white text-center">选择操作</h2>
            <div className="space-y-4">
              <button
                className="block w-full text-center py-3 bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-500 hover:to-blue-600 text-white rounded-lg transition-all duration-200 transform hover:scale-[1.02]"
                onClick={() => {
                  setModalType('addFriend');
                  setShowModal(false);
                }}
              >
                <div className="flex items-center justify-center space-x-2">
                  <AiOutlineUserAdd className="w-5 h-5" />
                  <span>添加好友</span>
                </div>
              </button>
              <button
                className="block w-full text-center py-3 bg-gradient-to-r from-green-600 to-green-700 hover:from-green-500 hover:to-green-600 text-white rounded-lg transition-all duration-200 transform hover:scale-[1.02]"
                onClick={() => {
                  setModalType('joinGroup');
                  setShowModal(false);
                }}
              >
                <div className="flex items-center justify-center space-x-2">
                  <AiOutlineTeam className="w-5 h-5" />
                  <span>加入群聊</span>
                </div>
              </button>
              <button
                className="block w-full text-center py-3 bg-gradient-to-r from-purple-600 to-purple-700 hover:from-purple-500 hover:to-purple-600 text-white rounded-lg transition-all duration-200 transform hover:scale-[1.02]"
                onClick={() => {
                  setShowCreateGroupModal(true);
                  setShowModal(false);
                }}
              >
                <div className="flex items-center justify-center space-x-2">
                  <AiOutlineUsergroupAdd className="w-5 h-5" />
                  <span>创建群聊</span>
                </div>
              </button>
            </div>
            <button
              className="absolute top-4 right-4 text-gray-400 hover:text-white transition-colors"
              onClick={() => setShowModal(false)}
            >
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
      )}

      {modalType && (
        <div className="fixed inset-0 bg-gray-800/70 backdrop-blur-sm flex justify-center items-center z-50">
          <div className="bg-gray-900 p-8 rounded-xl w-[480px] shadow-2xl border border-gray-700">
            <h2 className="text-2xl font-bold mb-6 text-white">
              {modalType === 'addFriend' ? '添加好友' : '加入群聊'}
            </h2>
            <div className="space-y-4">
              {renderSearchBox()}
              {renderSearchResults()}
              <div className="flex justify-end space-x-3 mt-6">
                <button
                  className="px-6 py-2 bg-gray-700 text-white rounded-lg hover:bg-gray-600 transition-colors"
                  onClick={() => {
                    setModalType(null);
                    setSearchInput('');
                    setSearchResults({ users: [], groups: [] });
                  }}
                >
                  关闭
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* 通知详情弹窗 */}
      {showNoticeModal && selectedNotice && (
        <div className="fixed inset-0 bg-gray-800/70 backdrop-blur-sm flex justify-center items-center z-50">
          <div className="bg-gray-900 p-6 rounded-xl w-96 max-w-lg border border-gray-700 relative">
            <div className="flex items-center space-x-4 mb-4">
              <img
                src={selectedNotice.extra.sender_avatar_url || '/static/images/user.png'}
                alt={selectedNotice.extra.sender_nickname}
                className="w-16 h-16 rounded-full border-2 border-blue-500"
              />
              <div>
                <h3 className="text-lg font-semibold text-white">
                  {selectedNotice.extra.sender_nickname}
                </h3>
                <p className="text-sm text-gray-400">
                  {selectedNotice.extra.invition_type === 'friend' ? '好友请求' : '群组邀请'}
                </p>
              </div>
            </div>

            <div className="mb-6">
              <p className="text-white text-lg">{selectedNotice.msg_content}</p>
              {selectedNotice.extra.invition_type === 'user' && (
                <p className="text-sm text-gray-400 mt-2">
                  群组名称：{selectedNotice.extra.group_name}
                </p>
              )}
              <p className="text-xs text-gray-500 mt-2">
                {new Date(selectedNotice.timestamp / 1000000).toLocaleString()}
              </p>
            </div>

            <div className="flex justify-end space-x-3">
              <button
                onClick={() => handleInvitation(selectedNotice, 'reject')}
                className="px-4 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700 transition-colors"
              >
                拒绝
              </button>
              <button
                onClick={() => handleInvitation(selectedNotice, 'accept')}
                className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
              >
                接受
              </button>
              <button
                onClick={() => setShowNoticeModal(false)}
                className="px-4 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700 transition-colors"
              >
                关闭
              </button>
            </div>

            {/* 关闭按钮 */}
            <button
              onClick={() => setShowNoticeModal(false)}
              className="absolute top-4 right-4 text-gray-400 hover:text-white"
            >
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
      )}

      {/* 创建群聊模态框 */}
      {showCreateGroupModal && (
        <div className="fixed inset-0 bg-gray-800/70 backdrop-blur-sm flex justify-center items-center z-50">
          <div className="bg-gray-900 p-6 rounded-xl w-96 max-w-lg border border-gray-700">
            <h3 className="text-xl font-semibold text-white mb-4">创建新群聊</h3>

            <form onSubmit={handleCreateGroup}>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-300 mb-1">
                    群聊名称 *
                  </label>
                  <input
                    type="text"
                    value={groupForm.name}
                    onChange={(e) => setGroupForm(prev => ({ ...prev, name: e.target.value }))}
                    className="w-full px-3 py-2 bg-gray-800 text-white rounded-lg border border-gray-700 focus:outline-none focus:border-blue-500"
                    required
                    placeholder="请输入群聊名称"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-300 mb-1">
                    群聊描述
                  </label>
                  <textarea
                    value={groupForm.description}
                    onChange={(e) => setGroupForm(prev => ({ ...prev, description: e.target.value }))}
                    className="w-full px-3 py-2 bg-gray-800 text-white rounded-lg border border-gray-700 focus:outline-none focus:border-blue-500"
                    placeholder="请输入群聊描述"
                    rows="3"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-300 mb-1">
                    群聊头像URL
                  </label>
                  <input
                    type="url"
                    value={groupForm.avatar_url}
                    onChange={(e) => setGroupForm(prev => ({ ...prev, avatar_url: e.target.value }))}
                    className="w-full px-3 py-2 bg-gray-800 text-white rounded-lg border border-gray-700 focus:outline-none focus:border-blue-500"
                    placeholder="请输入群聊头像URL"
                  />
                </div>
              </div>

              <div className="flex justify-end space-x-3 mt-6">
                <button
                  type="button"
                  onClick={() => setShowCreateGroupModal(false)}
                  className="px-4 py-2 bg-gray-700 text-white rounded-lg hover:bg-gray-600 transition-colors"
                >
                  取消
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-500 transition-colors"
                >
                  创建
                </button>
              </div>
            </form>

            {/* 关闭按钮 */}
            <button
              onClick={() => setShowCreateGroupModal(false)}
              className="absolute top-4 right-4 text-gray-400 hover:text-white"
            >
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
      )}

    </div>
  );
};

export default ChatSidebar;
