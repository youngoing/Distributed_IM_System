import React, { useEffect, useState } from 'react';
import { group_detail, friend_detail, delete_group, exit_group, delete_friend } from '../api_list';
import axios from 'axios';

const ChatHeader = ({ chatName, loading, isRoom, currentChatId, handleChatChange }) => {
  const [isDetailsVisible, setDetailsVisible] = useState(false);
  const [chatDetail, setChatDetail] = useState({});

  // 获取聊天详情
  const getChatDetail = async (isRoom, currentChatId) => {
    try {
      const response = isRoom
        ? await axios.get(group_detail(currentChatId), { withCredentials: true })
        : await axios.get(friend_detail(currentChatId), { withCredentials: true });
      setChatDetail(response.data);
      console.log('获取聊天详情成功:', response.data);
    } catch (error) {
      console.error('获取聊天详情失败:', error);
    }
  };

  const handleDeleteFriend = async (friendId) => {
    if (!window.confirm('确定要删除该好友吗？')) {
      return;
    }

    try {
      const response = await axios.delete(delete_friend(friendId), { withCredentials: true });

      if (response.status === 200) {
        alert('好友删除成功');
        if (currentChatId === friendId) {
          handleChatChange(null, null, false);
        }
        window.location.reload();
      }
    } catch (error) {
      console.error('删除好友失败:', error);
      alert('删除好友失败，请稍后重试');
    }
  };

  const handleDeleteGroup = async (groupId) => {
    if (!window.confirm('确定要删除该群聊吗？')) {
      return;
    }

    try {
      const response = await axios.delete(delete_group(groupId), { withCredentials: true });

      if (response.status === 200) {
        alert('群聊删除成功');
        if (currentChatId === groupId) {
          handleChatChange(null, null, false);
        }
        window.location.reload();
      }
    } catch (error) {
      console.error('删除群聊失败:', error);
      if (error.response?.status === 400) {
        alert('你是群主不能删除群聊');
      } else {
        alert('删除群聊失败，请稍后重试');
      }
    }
  };

  const handleExitGroup = async (groupId) => {
    if (!window.confirm('确定要退出该群聊吗？')) {
      return;
    }

    try {
      const response = await axios.post(exit_group(groupId), null, { withCredentials: true });

      if (response.status === 200) {
        alert('退出群聊成功');
        if (currentChatId === groupId) {
          handleChatChange(null, null, false);
        }
        window.location.reload();
      }
    } catch (error) {
      console.error('退出群聊失败:', error);
      if (error.response?.status === 400) {
        alert('你是群主不能退出群聊');
      } else {
        alert('退出群聊失败，请稍后重试');
      }
    }
  };

  const toggleDetails = () => {
    setDetailsVisible(!isDetailsVisible);
  };

  const closeDetails = () => {
    setDetailsVisible(false);
  };

  useEffect(() => {
    if (currentChatId) {
      getChatDetail(isRoom, currentChatId);
    }
  }, [isRoom, currentChatId]);

  const renderDetails = () => {
    if (isRoom) {
      return (
        <>
          <h2 className="text-xl font-semibold text-center">群聊详情</h2>
          <p className="mt-2 text-center text-gray-600">{chatDetail.description || "暂无群聊描述"}</p>
          <div className="mt-4">
            <h3 className="font-semibold">群成员</h3>
            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 mt-2">
              {chatDetail.members && chatDetail.members.map((member) => (
                <div key={member.user_detail_id} className="flex flex-col items-center text-center">
                  <img
                    src={member.avatar_url || '/static/images/user.png'}
                    alt={member.nickname}
                    className="w-16 h-16 rounded-full border-2 border-blue-500 mb-2"
                  />
                  <p className="text-sm font-medium">{member.nickname}</p>
                </div>
              ))}
            </div>
          </div>
          <div className="mt-4 flex justify-around">
            <button
              onClick={() => handleExitGroup(currentChatId)}
              className="bg-gray-500 text-white p-2 rounded-md hover:bg-gray-600 transition duration-300"
            >
              退出群聊
            </button>
            <button
              onClick={() => handleDeleteGroup(currentChatId)}
              className="bg-gray-500 text-white p-2 rounded-md hover:bg-gray-600 transition duration-300"
            >
              删除群聊
            </button>
          </div>
        </>
      );
    } else {
      return (
        <>
          <h2 className="text-xl font-semibold text-center">好友详情</h2>
          <div className="flex flex-col items-center mt-4">
            <img
              src={chatDetail.avatar_url || '/static/images/user.png'}
              alt={chatDetail.nickname}
              className="w-24 h-24 rounded-full border-2 border-green-500 mb-3"
            />
            <p className="text-lg font-medium">{chatDetail.nickname}</p>
          </div>
          <div className="mt-4 flex justify-center">
            <button
              onClick={() => handleDeleteFriend(currentChatId)}
              className="bg-gray-500 text-white p-2 rounded-md hover:bg-gray-600 transition duration-300"
            >
              删除好友
            </button>
          </div>
        </>
      );
    }
  };

  if (!currentChatId) {
    return (
      <header className="p-4 bg-gray-900 text-white flex justify-between items-center">
        没有聊天
      </header>
    );
  }

  return (
    <header className="p-4 bg-gray-900 text-white flex justify-between items-center">
      {loading ? (
        <div>加载中...</div>
      ) : (
        <h1 className={`text-xl font-semibold ${isRoom ? "text-center text-2xl text-blue-600" : "text-center text-2xl text-green-600"}`}>
          {isRoom ? `群聊：${chatName}` : `与 ${chatName} 的聊天`}
        </h1>
      )}

      <button
        onClick={toggleDetails}
        className="bg-gray-600 text-white p-2 rounded-md hover:bg-gray-700 transition duration-300"
      >
        详情
      </button>

      {isDetailsVisible && (
        <div className="fixed inset-0 bg-black bg-opacity-50 backdrop-blur-sm flex justify-center items-center z-50">
          <div className="relative bg-white text-black p-6 rounded-lg shadow-lg w-96 max-w-lg">
            <button
              onClick={closeDetails}
              className="absolute top-2 right-2 text-gray-600 hover:text-gray-900"
            >
              &times;
            </button>
            {renderDetails()}
          </div>
        </div>
      )}
    </header>
  );
};

export default ChatHeader;
