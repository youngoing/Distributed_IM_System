// components/MessageInput.js

import React from 'react';
import { FaSmile, FaPaperPlane } from 'react-icons/fa';
import EmojiPicker from 'emoji-picker-react';

const MessageInput = ({
  messageInput,
  setMessageInput,
  handleSendMessage,
  showEmojiPicker,
  setShowEmojiPicker,
  handleEmojiSelect,
  currentChatId
}) => {
  const handleKeyPress = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  };

  const onEmojiClick = (emojiData, event) => {
    setMessageInput(prev => prev + emojiData.emoji);
    setShowEmojiPicker(false);
  };

  if (!currentChatId) return null;

  return (
    <div className="relative border-t border-gray-200 bg-white p-4">
      {/* 表情选择器 */}
      {showEmojiPicker && (
        <div className="absolute bottom-full mb-2 left-4">
          <div className="relative">
            <div className="absolute bottom-0 left-0 transform -translate-y-2">
              <EmojiPicker
                onEmojiClick={onEmojiClick}
                disableSearchBar
                native
              />
            </div>
          </div>
        </div>
      )}

      <div className="flex items-center space-x-2">
        {/* 表情按钮 */}
        <button
          onClick={() => setShowEmojiPicker(!showEmojiPicker)}
          className="p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-full transition-colors"
        >
          <FaSmile className="w-5 h-5" />
        </button>

        {/* 消息输入框 */}
        <div className="flex-1">
          <textarea
            value={messageInput}
            onChange={(e) => setMessageInput(e.target.value)}
            onKeyPress={handleKeyPress}
            placeholder="输入消息..."
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 resize-none"
            rows="2"
          />
        </div>

        {/* 发送按钮 */}
        <button
          onClick={handleSendMessage}
          disabled={!messageInput.trim()}
          className={`p-2 rounded-full transition-colors ${
            messageInput.trim()
              ? 'bg-blue-500 hover:bg-blue-600 text-white'
              : 'bg-gray-200 text-gray-400 cursor-not-allowed'
          }`}
        >
          <FaPaperPlane className="w-5 h-5" />
        </button>
      </div>
    </div>
  );
};

export default MessageInput;