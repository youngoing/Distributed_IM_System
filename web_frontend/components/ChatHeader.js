import React, { useEffect, useState } from 'react';
import { group_detail, friend_detail } from '../api_list';
import axios from 'axios';

const ChatHeader = ({ chatName, loading, isRoom, currentChatId }) => {
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

  // 切换详细信息显示
  const toggleDetails = () => {
    setDetailsVisible(!isDetailsVisible);
  };

  // 关闭详细信息
  const closeDetails = () => {
    setDetailsVisible(false);
  };

  // 使用 useEffect 来获取聊天详情
  useEffect(() => {
    if (currentChatId) {
      getChatDetail(isRoom, currentChatId);
    }
  }, [isRoom, currentChatId]);

  // 渲染群聊或好友详情
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
        </>
      );
    }
  };
  if (!currentChatId){
    return (
      <header className="p-4 bg-gray-900 text-white flex justify-between items-center">
        没有聊天
      </header>

    )
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
        className="bg-blue-500 text-white p-2 rounded-md hover:bg-blue-600 transition duration-300"
      >
        详细
      </button>

      {isDetailsVisible && (
        // 弹出框背景
        <div className="fixed inset-0 bg-black bg-opacity-50 backdrop-blur-sm flex justify-center items-center z-50">
          {/* 弹出框内容 */}
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
