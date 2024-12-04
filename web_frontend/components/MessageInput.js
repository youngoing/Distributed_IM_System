// components/MessageInput.js

import React from 'react';
import data from '@emoji-mart/data';
import Picker from '@emoji-mart/react';

const MessageInput = ({ messageInput, setMessageInput, handleSendMessage, showEmojiPicker, setShowEmojiPicker, handleEmojiSelect,currentChatId }) => {
  if (!currentChatId) return null;

  return (
    
    <footer className="p-4 bg-gray-200 flex flex-col relative">
      <div className="flex items-center space-x-2">
        <button onClick={() => setShowEmojiPicker((prev) => !prev)} className="p-2 bg-gray-300 rounded-l-lg">
          😊
        </button>
        <input
          type="text"
          value={messageInput}
          onChange={(e) => setMessageInput(e.target.value)}
          onKeyDown={(e) => e.key === 'Enter' && handleSendMessage()}
          className="flex-grow p-2 border border-gray-300 rounded-l-lg"
          placeholder="输入消息..."
        />
        <button onClick={handleSendMessage} className="bg-blue-500 hover:bg-blue-600 text-white py-2 px-4 rounded-r-lg">
          发送
        </button>
      </div>
      {showEmojiPicker && (
        <div className="absolute bottom-0 left-0 w-full z-50 mt-2">
          <Picker data={data} onEmojiSelect={handleEmojiSelect} />
        </div>
      )}
    </footer>
  );
};

export default MessageInput;