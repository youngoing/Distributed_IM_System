import { useState, useEffect } from 'react';
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

  const handleEmojiSelect = (emoji) => {
    setMessageInput((prev) => prev + emoji.native);
    setShowEmojiPicker(false);
  };

  const handleChatChange = async (id, name, isRoom) => {
    if (currentChatId === id) return;
    setCurrentChatId(id);
    setChatName(name);
    setIsRoom(isRoom);
    setLoading(true);
    // console.log("isRoom:", isRoom);
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
    setChatMessages((prevMessages) => {
      const newMessage = {
        sender_id: user.user_detail_id,
        msg_content: messageInput,
        timestamp: new Date().toISOString(),
      };
      return [...prevMessages, newMessage];
    });
    setMessageInput('');
    if (isRoom) {
      const groupMessage = newTestGroupMessage(user.user_detail_id, currentChatId, groupMemberIds, messageInput);
      WebSocketService.sendMessage(groupMessage);
      console.log('groupMessage:', groupMessage);
    } else {
      const privateMessage = newTestPrivateMessage(user.user_detail_id, currentChatId, messageInput);
      WebSocketService.sendMessage(privateMessage);
      console.log('privateMessage:', privateMessage);
    }
  };

  const handleIncomingMessage = async (message) => {
    setChatMessages((prevMessages) => {
      if (message.msg_type === 'group' && isRoom && String(currentChatId) === message.group_id) {
        return [...prevMessages, message];
      }
      if (message.msg_type === 'private' && !isRoom && String(currentChatId) === message.sender_id) {
        return [...prevMessages, message];
      }
      return prevMessages;
    });
  };

  const addMembers = (data) => {
    const newMembers = { ...members };
    data.forEach(member => {
      const id = member.user_detail_id;
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
      const MemberIds = data.map(member => String(member.user_detail_id));
      setGroupMemberIds(MemberIds);
    }
  };
  // console.log("members:", members);
  // console.log("groupMemberIds:", groupMemberIds);

  useEffect(() => {
    if (user && user.user_detail_id) {
      WebSocketService.connect(user.user_detail_id);
      WebSocketService.onMessage(handleIncomingMessage);
      return () => {
        WebSocketService.onMessage(() => { });
      };
    }
  }, [user, currentChatId, isRoom]);

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

  // useEffect(() => {
  //   if (!currentChatId && groups.length > 0) {
  //     handleChatChange(groups[0].group_id, groups[0].name, true);
  //   }
  // }, [groups, currentChatId]);

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
            <ChatHeader chatName={chatName} loading={loading} isRoom={isRoom} currentChatId={currentChatId}/>
            <ChatMessages chatMessages={chatMessages} members={members} user={user} currentChatId = {currentChatId} />
            <MessageInput
              messageInput={messageInput}
              setMessageInput={setMessageInput}
              handleSendMessage={handleSendMessage}
              showEmojiPicker={showEmojiPicker}
              setShowEmojiPicker={setShowEmojiPicker}
              handleEmojiSelect={handleEmojiSelect}
              currentChatId = {currentChatId}
            />
          </main>
        </div>
      </div>
    </ProtectedComponent>
  );
}
