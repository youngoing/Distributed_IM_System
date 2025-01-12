import { useState, useEffect, useCallback } from 'react';
import axios from 'axios';
import { newTestPrivateMessage, newTestGroupMessage } from '../utils/msgModel';
import { friend_list, group_list } from '../api_list';
import Loginout from '../components/Loginout';
import ProtectedComponent from '../components/ProtectedComponent';
import WebSocketService from '../components/WebsocketService';
import ChatSidebar from '../components/ChatSidebar';
import ChatHeader from '../components/ChatHeader';
import ChatMessages from '../components/ChatMessages';
import MessageInput from '../components/MessageInput';

export default function Msg({ children }) {
  const [isSidebarOpen, setSidebarOpen] = useState(true);  // 默认侧边栏展开
  const [currentChatId, setCurrentChatId] = useState(null);
  const [chatName, setChatName] = useState(null);
  const [isRoom, setIsRoom] = useState(false);
  const [loading, setLoading] = useState(false);
  const [chatMessages, setChatMessages] = useState([]);
  // 通知消息
  const [noticeMessages, setNoticeMessages] = useState([]);
  const [members, setMembers] = useState({});
  const [groupMemberIds, setGroupMemberIds] = useState([]);
  const [friends, setFriends] = useState([]);
  const [groups, setGroups] = useState([]);
  const [messageInput, setMessageInput] = useState('');
  const [user, setUser] = useState({});
  const [showLogoutConfirm, setShowLogoutConfirm] = useState(false);
  const [showEmojiPicker, setShowEmojiPicker] = useState(false);
  const [showDropdown, setShowDropdown] = useState(false);

  const handleDropdownToggle = () => {
    setShowDropdown(!showDropdown);
  };

  const handleEmojiSelect = (emojiData) => {
    setMessageInput(prev => prev + emojiData.emoji);
    setShowEmojiPicker(false);
  };

  const handleChatChange = async (id, name, isRoom) => {
    if (currentChatId === id) return;
    setCurrentChatId(id);
    setChatName(name);
    setIsRoom(isRoom);
    setLoading(true);
    try {
      let loadedMessages = [];
      const storageKey = isRoom ? 'group_msgs' : 'private_msgs';
      const storedMessages = localStorage.getItem(storageKey);
      if (storedMessages) {
        const messages = JSON.parse(storedMessages);
        loadedMessages = messages[id] || [];
      }
      setChatMessages(loadedMessages);
    } catch (error) {
      console.error('加载聊天消息失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSendMessage = async () => {
    if (!messageInput.trim()) return;
    console.log(currentChatId)

    // 创建一个与接收消息格式一致的消息对象
    const newMessage = {
      sender_id: user.user_detail_id,
      receiver_id: isRoom ? null : currentChatId,
      group_id: isRoom ? currentChatId : null,
      msg_content: messageInput,
      timestamp: new Date().toISOString(),
      msg_type: isRoom ? 'group' : 'private'
    };
    // 更新聊天消息列表
    setChatMessages(prevMessages => [...prevMessages, newMessage]);
    setMessageInput('');
    if (isRoom) {
      console.log("groupMemberIds",groupMemberIds)
      const groupMessage = newTestGroupMessage(user.user_detail_id, String(currentChatId), groupMemberIds, messageInput);
      WebSocketService.sendMessage(groupMessage);
    } else {
      const privateMessage = newTestPrivateMessage(user.user_detail_id, String(currentChatId), messageInput);
      WebSocketService.sendMessage(privateMessage);
    }
  };

  const handleIncomingMessage = useCallback(async (message) => {
    setChatMessages((prevMessages) => {
        // 检查消息是否已存在
        const isDuplicate = prevMessages.some(msg =>
            msg.timestamp === message.timestamp &&
            msg.sender_id === message.sender_id &&
            msg.msg_content === message.msg_content
        );

        if (isDuplicate) {
            console.log("消息重复，不更新");
            return prevMessages;
        }

        // 群聊消息处理
        if (isRoom && message.group_id) {
            if (String(currentChatId) === String(message.group_id)) {
                console.log("群聊消息匹配成功，更新消息列表");
                return [...prevMessages, message];
            }
        } 
        // 私聊消息处理
        else if (!isRoom && !message.group_id) {
            if (String(currentChatId) === String(message.sender_id)) {
                console.log("私聊消息匹配成功，更新消息列表");
                return [...prevMessages, message];
            }
        }
        console.log("消息不匹配当前聊天");
        return prevMessages;
    });
  }, [currentChatId, isRoom]);

  const addMembers = (data) => {
    const newMembers = { ...members };
    data.forEach(member => {
      const id = member.user_detail_id;
      // 跳过当前用户
      if (id === user.user_detail_id) {
        return;
      }
      if (!newMembers[id]) {
        newMembers[id] = {
          nickname: member.nickname,
          avatar_url: member.avatar_url
        };
      }
    });
    setMembers(prevMembers => ({
      ...prevMembers,
      ...newMembers
    }));

    if (isRoom) {
      // 使用 useRef 来保存完整的群成员列表
      const currentGroupMembers = groups.find(g => g.group_id === currentChatId)?.members || [];
      const memberIds = currentGroupMembers
        .filter(member => String(member.user_detail_id) !== String(user.user_detail_id))
        .map(member => String(member.user_detail_id));

      console.log("群成员IDs:", memberIds);
      setGroupMemberIds(memberIds);
    }
  };
// 添加useEffect来监听状态变化
useEffect(() => {
  console.log("状态已更新:", {
      currentChatId,
      isRoom,
      chatName
  });
}, [currentChatId, isRoom, chatName]);

  useEffect(() => {
    if (user && user.user_detail_id) {
      // 确保只在组件挂载时连接一次
      WebSocketService.connect(user.user_detail_id);
      WebSocketService.onMessage(handleIncomingMessage);

      return () => {
        // 清理时移除特定的处理函数
        WebSocketService.removeMessageHandler(handleIncomingMessage);
      };
    }
  }, [user, handleIncomingMessage]);

  useEffect(() => {
    if (!localStorage.getItem('group_msgs')) {
      localStorage.setItem('group_msgs', JSON.stringify({}));
    }
    if (!localStorage.getItem('private_msgs')) {
      localStorage.setItem('private_msgs', JSON.stringify({}));
    }
  }, []);

  useEffect(() => {
    const storedUser = sessionStorage.getItem('user');
    if (storedUser) {
      const parsedUser = JSON.parse(storedUser);
      setUser(parsedUser);
    }
  }, []);

  useEffect(() => {
    const fetchFriendsAndGroups = async () => {
      if (user && user.user_detail_id) {
        try {
          const [friendsRes, groupsRes] = await Promise.all([
            axios.get(friend_list(user.user_detail_id), { withCredentials: true }),
            axios.get(group_list(user.user_detail_id), { withCredentials: true }),
          ]);
          setFriends(friendsRes.data);
          setGroups(groupsRes.data);
        } catch (error) {
          console.error('获取好友或群聊列表失败:', error);
        }
      }
    };
    if (user && user.user_detail_id) {
      fetchFriendsAndGroups();
    }
  }, [user]);

  useEffect(() => {
    if (Array.isArray(groups)) {
      groups.forEach(group => {
        if (group.members) {
          addMembers(group.members, true);
        }
      });
    }
    if (Array.isArray(friends)) {
      friends.forEach(friend => {
        addMembers([friend], false);
      });
    }
  }, [groups, currentChatId, friends]);

  return (
    <ProtectedComponent>
      {showLogoutConfirm && (
        <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50">
          <div className="bg-white p-6 rounded-lg shadow-lg">
            <p className="text-lg font-semibold mb-4">确认要退出登录吗？</p>
            <div className="flex space-x-4">
              <button onClick={() => Loginout()} className="bg-red-600 hover:bg-red-700 text-white py-2 px-4 rounded-lg">
                确认
              </button>
              <button onClick={() => setShowLogoutConfirm(false)} className="bg-gray-300 hover:bg-gray-400 text-gray-800 py-2 px-4 rounded-lg">
                取消
              </button>
            </div>
          </div>
        </div>
      )}
      <div className="flex justify-center items-center min-h-screen bg-cover bg-center" style={{ backgroundImage: "url('/static/images/bg.jpg')" }}>
        <div className="w-full max-w-6xl bg-white shadow-lg rounded-lg overflow-hidden flex h-screen">
          <div className="relative">

            <ChatSidebar
              groups={groups}
              friends={friends}
              handleChatChange={handleChatChange}
              user={user}
              onLogout={Loginout}
              isSidebarOpen={isSidebarOpen}  // 控制侧边栏的展开与折叠
              toggleSidebar={() => setSidebarOpen(!isSidebarOpen)}  // 切换侧边栏展开与折叠
            />
          </div>

          <main className="flex-grow flex flex-col">
            <ChatHeader chatName={chatName} loading={loading} isRoom={isRoom} currentChatId={currentChatId} />
            <ChatMessages chatMessages={chatMessages} members={members} user={user} currentChatId={currentChatId} />
            <MessageInput
              messageInput={messageInput}
              setMessageInput={setMessageInput}
              handleSendMessage={handleSendMessage}
              showEmojiPicker={showEmojiPicker}
              setShowEmojiPicker={setShowEmojiPicker}
              handleEmojiSelect={handleEmojiSelect}
              currentChatId={currentChatId}
            />
          </main>
        </div>
      </div>
    </ProtectedComponent>
  );
}
