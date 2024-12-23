import React, { useState } from 'react';
import { AiOutlineArrowRight, AiOutlineArrowLeft, AiOutlinePlus } from 'react-icons/ai'; // 导入图标
import { FaArrowLeft } from 'react-icons/fa';
const ChatSidebar = ({ groups, friends, handleChatChange, user, onLogout, isSidebarOpen, toggleSidebar,currentChatId }) => {
  const [showProfileMenu, setShowProfileMenu] = useState(false);
  const [isGroupsCollapsed, setIsGroupsCollapsed] = useState(false);
  const [isFriendsCollapsed, setIsFriendsCollapsed] = useState(false);
  const [showModal, setShowModal] = useState(false);
  const [modalType, setModalType] = useState(''); // 'addFriend' or 'joinGroup'
  const [searchQuery, setSearchQuery] = useState('');
  const [filteredResults, setFilteredResults] = useState([]);

  const handleProfileClick = () => {
    setShowProfileMenu((prevState) => !prevState);
  };

  const toggleGroupsCollapse = () => {
    setIsGroupsCollapsed((prevState) => !prevState);
  };

  const toggleFriendsCollapse = () => {
    setIsFriendsCollapsed((prevState) => !prevState);
  };

  const handleSearch = (query) => {
    setSearchQuery(query);
    const results = [...groups, ...friends].filter(item =>
      item.name.toLowerCase().includes(query.toLowerCase()) ||
      item.nickname.toLowerCase().includes(query.toLowerCase())
    );
    setFilteredResults(results);
  };

  const handleAddFriend = () => {
    // 发送添加好友请求的逻辑
    console.log('发送添加好友请求');
  };

  const handleJoinGroup = () => {
    // 发送加入群聊请求的逻辑
    console.log('发送加入群聊请求');
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

        {/* 搜索和添加区域 */}
        <div className="p-4 bg-gray-750 border-b border-gray-700">
          <div className="flex items-center space-x-2">
            <input
              type="text"
              placeholder="搜索聊天..."
              value={searchQuery}
              onChange={(e) => handleSearch(e.target.value)}
              className="flex-1 bg-gray-700 text-white text-sm rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <button
              onClick={() => setShowModal(true)}
              className="p-2 hover:bg-gray-600 rounded-md transition-colors"
              title="添加"
            >
              <AiOutlinePlus className="w-5 h-5" />
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
                            <span className="text-xs text-gray-400">{group.members?.length || 0}人</span>
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
                            <span className="text-xs text-gray-400">
                              {friend.online ? '在线' : '离线'}
                            </span>
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

      {/* 弹出框：添加好友或加入群聊 */}
      {showModal && (
        <div
          className="fixed inset-0 bg-gray-800 bg-opacity-50 flex justify-center items-center z-50"
          onClick={() => setShowModal(false)} // 点击外部区域关闭弹出框
        >
          <div
            className="bg-white p-6 rounded-lg w-1/3"
            onClick={(e) => e.stopPropagation()} // 阻止点击弹出框内部时关闭弹出框
          >
            <h2 className="text-xl font-semibold mb-4">选择操作</h2>
            <button
              className="block w-full text-center py-2 bg-blue-600 text-white rounded-lg mb-4"
              onClick={() => {
                setModalType('addFriend');
                setShowModal(false); // 关闭弹出框
              }}
            >
              添加好友
            </button>
            <button
              className="block w-full text-center py-2 bg-green-600 text-white rounded-lg mb-4"
              onClick={() => {
                setModalType('joinGroup');
                setShowModal(false); // ���闭弹出框
              }}
            >
              加入群聊
            </button>
            {/* 添加关闭按钮 */}
            <button
              className="absolute top-2 right-2 text-gray-600 hover:text-gray-800"
              onClick={() => setShowModal(false)} // 点击关闭按钮关闭弹出框
            >
              &times; {/* 使用 x 图标作为关闭按钮 */}
            </button>
          </div>
        </div>
      )}


      {/* 处理添加好友或加入群聊的弹出框 */}
      {modalType === 'addFriend' && (
        <div className="fixed inset-0 bg-gray-800 bg-opacity-50 flex justify-center items-center z-50">
          <div className="bg-white p-6 rounded-lg w-1/3">
            <h2 className="text-xl font-semibold mb-4">添加好友</h2>
            <input
              type="text"
              placeholder="搜索好友..."
              value={searchQuery}
              onChange={(e) => handleSearch(e.target.value)}
              className="border p-2 rounded w-full mb-4"
            />
            <ul className="space-y-2 max-h-60 overflow-y-auto">
              {filteredResults.map((item) => (
                <li key={item.id} className="p-2 border-b">
                  <span>{item.nickname}</span>
                  <button className="ml-4 text-blue-600" onClick={handleAddFriend}>添加</button>
                </li>
              ))}
            </ul>
            <button className="mt-4 bg-red-600 text-white py-2 px-4 rounded" onClick={() => setShowModal(false)}>
              关闭
            </button>
          </div>
        </div>
      )}

 {/* 处理添加好友或加入群聊的弹出框 */}
 {modalType === 'addFriend' && (
        <div className="fixed inset-0 bg-gray-800 bg-opacity-50 flex justify-center items-center z-50">
          <div className="bg-white p-6 rounded-lg w-1/3">
            <h2 className="text-xl font-semibold mb-4">添加好友</h2>
            <input
              type="text"
              placeholder="搜索好友..."
              value={searchQuery}
              onChange={(e) => handleSearch(e.target.value)}
              className="border p-2 rounded w-full mb-4"
            />
            <ul className="space-y-2 max-h-60 overflow-y-auto">
              {filteredResults.map((item) => (
                <li key={item.id} className="p-2 border-b">
                  <span>{item.nickname}</span>
                  <button className="ml-4 text-blue-600" onClick={() => handleAddFriend(item.id)}>添加</button>
                </li>
              ))}
            </ul>
            <button className="mt-4 bg-red-600 text-white py-2 px-4 rounded" onClick={() => setModalType(null)}>
              关闭
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default ChatSidebar;
