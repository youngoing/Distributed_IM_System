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
      <aside className={`transition-all duration-300 ${isSidebarOpen ? 'w-64' : 'w-0'} bg-gray-800 text-white flex flex-col min-h-screen overflow-hidden`}>
        <header className="p-4 bg-gray-900 flex justify-between items-center">
          <div className="flex items-center">
            <div className="relative">
              <img
                src={user?.avatar_url ? user.avatar_url : '/static/images/me.png'}
                alt={user?.nickname || 'User Avatar'}
                className="w-10 h-10 rounded-full cursor-pointer"
                onClick={handleProfileClick}
                aria-label="User Profile"
              />
              {isSidebarOpen && <span className="text-lg font-semibold text-yellow-400 ml-3">{user?.nickname}</span>}
              {showProfileMenu && (
                <div className="absolute right-0 mt-2 bg-gray-700 rounded-md shadow-lg">
                  <div className="p-4 cursor-pointer hover:bg-gray-600" onClick={onLogout}>
                    退出登录
                  </div>
                </div>
              )}
            </div>
          </div>
        </header>

        <div className="flex-grow overflow-y-auto">
          <ul className="space-y-2">
            {/* 群聊搜索框和加号按钮 */}
            <div className="flex justify-between items-center p-4 bg-gray-700">
              <input
                type="text"
                placeholder="搜索群聊..."
                value={searchQuery}
                onChange={(e) => handleSearch(e.target.value)}
                className="bg-gray-600 text-white rounded-md p-2 w-3/4"
              />
              <AiOutlinePlus
                size={24}
                className="cursor-pointer text-white ml-4"
                onClick={() => setShowModal(true)} // 点击加号按钮，显示弹出框
              />
            </div>

            {/* 群聊部分 */}
            {groups.length === 0 && friends.length === 0 && (
              <div className="text-center text-gray-400 mt-8">
                <p>还没有群聊</p>
                <p>还没有朋友</p>
              </div>
            )}

            {groups.length > 0 && (
              <li className="p-4 bg-gray-700 cursor-pointer flex justify-between items-center" onClick={toggleGroupsCollapse}>
                <h2 className={`${!isSidebarOpen ? 'text-sm' : 'text-lg'} font-semibold`}>
                  群聊列表 ({groups.length})
                </h2>
                <span>{isGroupsCollapsed ? '▶' : '▼'}</span>
              </li>
            )}

            {!isGroupsCollapsed && groups.length > 0 && (
              <div className={`transition-all duration-300 ${isSidebarOpen ? 'max-h-96' : 'max-h-0'} overflow-hidden`}>
                {groups.map((group) => (
                  <li key={group.group_id} onClick={() => handleChatChange(group.group_id, group.name, true)} className="p-4 border-b border-gray-700 hover:bg-gray-700 cursor-pointer flex items-center">
                    {isSidebarOpen && <img src={group.avatar_url ? group.avatar_url : '/static/images/group.png'} alt={group.name} className="w-10 h-10 rounded-full mr-4" />}
                    {isSidebarOpen ? group.name : ''}
                  </li>
                ))}
              </div>
            )}

            {friends.length > 0 && (
              <li className="p-4 bg-gray-700 cursor-pointer flex justify-between items-center" onClick={toggleFriendsCollapse}>
                <h2 className={`${!isSidebarOpen ? 'text-sm' : 'text-lg'} font-semibold`}>
                  私聊列表 ({friends.length})
                </h2>
                <span>{isFriendsCollapsed ? '▶' : '▼'}</span>
              </li>
            )}

            {!isFriendsCollapsed && friends.length > 0 && (
              <div className={`transition-all duration-300 ${isSidebarOpen ? 'max-h-96' : 'max-h-0'} overflow-hidden`}>
                {friends.map((friend) => (
                  <li key={friend.user_detail_id} onClick={() => handleChatChange(friend.user_detail_id, friend.nickname, false)} className="p-4 border-b border-gray-700 hover:bg-gray-700 cursor-pointer flex items-center">
                    {isSidebarOpen && <img src={friend.avatar_url ? friend.avatar_url : '/static/images/user.png'} alt={friend.nickname} className="w-10 h-10 rounded-full mr-4" />}
                    {isSidebarOpen ? friend.nickname : ''}
                  </li>
                ))}
              </div>
            )}
          </ul>
        </div>
      </aside>

      {/* 折叠按钮 */}
      <button className="absolute top-1/2 left-full transform -translate-x-1/2 translate-y-[-50%] bg-gray-700 text-white p-3 rounded-full focus:outline-none shadow-lg hover:bg-gray-600 transition-all duration-300" onClick={toggleSidebar}>
        {isSidebarOpen ? <AiOutlineArrowLeft size={24} /> : <AiOutlineArrowRight size={24} />}
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
                setShowModal(false); // 关闭弹出框
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
