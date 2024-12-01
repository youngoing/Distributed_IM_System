
# Distributed_IM_System

A scalable, real-time Instant Messaging (IM) system designed with WebSocket for efficient communication. It supports both group and direct chats, offline message handling, and dynamic load balancing using Nginx and Consul.

### Architecture

1. **Client**:  
   Users connect to the server via WebSocket over HTTP using Web, desktop, or mobile clients. Real-time message sending and receiving is handled through WebSocket.

2. **Load Balancing**:  
   WebSocket connection requests from clients are first directed through an Nginx load balancer. The Nginx configuration integrates with Consul to dynamically obtain healthy WebSocket nodes, balancing the load and ensuring high availability across the system.

3. **WebSocket Node Servers**:  
   After load balancing, requests are routed to multiple WebSocket node servers. Each node server uses Go’s goroutines to handle client connections concurrently. Clients are authenticated via Token validation before establishing a WebSocket connection for two-way communication. Additionally, each node runs an HTTP service for health checks, monitored by Consul.

4. **RabbitMQ Message Queue**:  
   WebSocket nodes communicate through RabbitMQ, which acts as a message broker. It employs a producer/consumer model, receiving messages from WebSocket nodes and forwarding them to other nodes or persisting them in the database. This ensures reliable message delivery and system scalability.

![Architecture](https://github.com/user-attachments/assets/eb54931d-5d11-449e-a678-4acd2ac82d6e)

### Message Forwarding Logic

1. **User Login and Connection Establishment**:  
   After login, clients establish a WebSocket connection with the server. Nginx routes the request to an available node. Each node keeps track of the user's ID and associated node information in Redis.

2. **Offline Message Handling**:  
   Once a connection is established, the node checks the database for any offline messages for the user. If there are any, they are sent to the client to ensure timely message delivery.

3. **Message Sending and Forwarding**:  
   Messages sent in group or private chats are routed to a designated message queue (e.g., `input` queue) via the `receive_id` field. A service consumes these messages and checks the user’s online status via Redis.

4. **Message Queue Consumption and Forwarding**:  
   - If the target user is online, the message is forwarded to their dedicated message queue.
   - Each node has its own dedicated message queue to consume and forward messages from other nodes. After processing, messages are pushed to the target user via WebSocket for real-time delivery.

### Group Chat Example

![Group Chat](https://github.com/user-attachments/assets/66265905-bb20-4b66-bca6-38828b869a5f)

### Private Chat Example

![Private Chat](https://github.com/user-attachments/assets/f8f38283-1f10-4918-a048-72b845771a2d)

### Implemented Features:
- Distributed node management for continuous operation
- User authentication (login, registration)
- Private and group chat functionality
- Offline message storage and delivery
- Chat room details display

### Unimplemented Features:
- Add Friend
- Accept Friend Requests
- Create Group Chat
- Invite Friends to Group
- Group Management
- User Profile Editing
- Delete Friend


