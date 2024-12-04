// 定义常量
const MsgStatusSent = "sent";         // 消息已发送
const MsgStatusDelivered = "delivered"; // 消息已送达
const MsgStatusRead = "read";         // 消息已读

const MsgTypeUser = "private";       // 用户消息
const MsgTypeGroup = "group";        // 群组消息
const MsgTypeAck = "ack";            // Ack 消息
const MsgTypeSystem = "system";      // 系统消息

const MsgContentTypeText = "text";   // 消息类型：文本消息
const MsgContentTypeImage = "image"; // 消息类型：图片消息
class WsMsg {
    constructor(msgType, status, wsHeader, wsBody) {
        this.msg_id = generateRandomID(); // 自动生成唯一的消息 ID
        this.msg_type = msgType;
        this.status = status;
        this.wsHeader = wsHeader;
        this.wsBody = wsBody;
    }
}

class WsHeader {
    constructor(senderId, receiverId = [], msgContentType, timestamp, groupId = null) {
        this.sender_id = senderId;
        this.receiver_id = receiverId; // 这里是数组，包含所有接收者的 ID
        this.msg_content_type = msgContentType;
        this.timestamp = timestamp;
        this.group_id = groupId; // 只有群消息时才会设置
    }
}
class WsBody {
    constructor(msgContent) {
        this.msg_content = msgContent;
    }
}

function newPrivateMessage(senderId, msgContentType, content,receiverId=[]) {
    return new WsMsg(
        MsgTypeUser,
        MsgStatusSent,
        new WsHeader(
            senderId,
            receiverId,
            msgContentType,
            Date.now()
        ),
        new WsBody(content)
    );
}
function newGroupMessage(senderId, groupId, msgContentType, content, receiverId=[]) {
    return new WsMsg(
        MsgTypeGroup,
        MsgStatusSent,
        new WsHeader(
            senderId,
            receiverId,
            msgContentType,
            Date.now(),
            groupId
        ),
        new WsBody(content)
    );
}
// flatten 函数


// 生成随机 ID 的函数
function generateRandomID() {
    // 生成 16 字节的随机数
    const bytes = crypto.getRandomValues(new Uint8Array(16));
    // 将字节数组转换为十六进制字符串
    return Array.from(bytes).map(byte => byte.toString(16).padStart(2, '0')).join('');
}
function newTestPrivateMessage(senderId, receiverId, content) {
    let msg = {
        //随机生成消息ID
        "msg_id": generateRandomID(),
        "msg_type": "private",
        "status": "sent",
        "sender_id": senderId.toString(),
        "receiver_id": [
            receiverId
        ],
        "msg_content_type": "text",
        "timestamp": Date.now(),
        "msg_content": content
    }
    return msg;


}
//receiverIds为数组
function newTestGroupMessage(senderId, groupId, receiverIds, content) {
    
    let msg ={
        "msg_id": generateRandomID(),
        "msg_type": "group",
        "status": "sent",
        "sender_id": senderId.toString(),
        "receiver_id": receiverIds, // 直接使用传递的数组
        "msg_content_type": "text",
        "timestamp": Date.now(),
        "group_id": groupId.toString(),
        "msg_content": content
      }
      return msg;


}
// 导出函数
export { newTestPrivateMessage, newTestGroupMessage };