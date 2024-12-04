// components/ChatMessages.js
import { FaArrowLeft } from 'react-icons/fa';
import React from 'react';

const ChatMessages = ({ chatMessages, members, user,currentChatId }) => {
  if (!currentChatId) {
    return (
      <div className="flex flex-col items-center justify-center h-full text-center p-4 bg-gray-500 text-white">
        <FaArrowLeft className="text-6xl text-gray-400 mb-4 animate-bounce" />
        <p className="text-xl font-semibold">点击左侧开始聊天</p>
        <p className="text-sm">选择一个群聊或好友开始聊天</p>
      </div>
    );
  }
  return (
    <div
    className="flex-grow overflow-y-auto p-4"
    style={{
      backgroundImage: "url('static/images/chat-bg2.jpg')",  // 使用相对路径
      backgroundSize: 'cover',
      backgroundPosition: 'center',
    }}
    >
      {chatMessages.length > 0 ? (
        chatMessages.map((message, index) => {
          const isOwnMessage = String(message.sender_id) === String(user.user_detail_id);
          const senderAvatar = members[message.sender_id]?.avatar_url || 'static/images/user.png';
          const senderNickname = members[message.sender_id]?.nickname || '匿名用户';
          const mineNickName = user?.nickname || '我';
          const mineAvatar = user?.avatar_url || 'static/images/me.png';

          return (
            <div key={index} className={`flex ${isOwnMessage ? 'justify-end' : 'justify-start'} mb-4`}>
              <div className={`flex ${isOwnMessage ? 'flex-row-reverse' : 'flex-row'} items-start space-x-3`}>
                <img src={isOwnMessage ? mineAvatar : senderAvatar} alt="Avatar" className="w-10 h-10 rounded-full" />
                <div className={`flex flex-col ${isOwnMessage ? 'items-end' : 'items-start'}`}>
                  <div className={`${isOwnMessage ? 'text-blue-600' : 'text-white-800'} text-sm font-semibold`}>
                    {isOwnMessage ? mineNickName : senderNickname}
                  </div>
                  <div className={`${isOwnMessage ? 'bg-blue-500 text-white' : 'bg-gray-200 text-gray-800'} p-3 rounded-lg max-w-xs`}>
                    <p>{message.msg_content}</p>
                    <small className="block text-xs text-right">{new Date(message.timestamp).toLocaleTimeString()}</small>
                  </div>
                </div>
              </div>
            </div>
          );
        })
      ) : (
        <div className="flex flex-col items-center justify-center text-center text-gray-400 py-8">
          <p className="text-xl font-semibold">没有消息</p>
          <p className="text-sm">你还没有收到任何消息。</p>
        </div>
      )}
    </div>
  );
};

export default ChatMessages;