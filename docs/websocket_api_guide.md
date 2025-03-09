# WebSocket API 接口文档

> **文档版本**：1.1  
> **最后更新**：2024-07-01  
> **适用版本**：SLA2 v1.0+

本文档描述了 SLA2 系统中 AI 服务的 WebSocket 接口，实现实时交互和流式消息传输。

## 连接信息

### 连接地址
```
ws://<服务器地址>:9102/ws/chat
```
> **注意**：9102 是系统默认的 HTTP 网关端口

### 认证方式
WebSocket连接支持两种认证方式：

1. **HTTP请求头**：`Authorization: Bearer <YOUR_JWT_TOKEN>`

2. **URL查询参数**：`ws://<服务器地址>:9102/ws/chat?token=<YOUR_JWT_TOKEN>`

**认证说明**：
- 认证仅在WebSocket连接建立阶段（HTTP握手阶段）进行，使用HTTP请求头或URL参数传递令牌
- 浏览器环境推荐使用URL查询参数方式
- 连接建立后无需在消息中包含令牌，用户身份会保持在整个连接期间
- 认证失败返回401状态码
- WebSocket消息的context字段不用于认证，仅用于传递聊天上下文信息

## 消息格式

支持两种消息格式：**JSON** 和 **Protobuf**，系统会自动检测并处理。

### JSON 格式

#### 请求消息
```json
{
  "type": "chat",
  "action": "message",
  "message": "你好，AI助手！",
  "sessionId": "session-123",
  "streamId": "stream-456",
  "context": {
    "sessionId": "session-123",
    "history": ["用户之前的消息", "AI之前的回复"]
  }
}
```

#### 终止消息
```json
{
  "type": "stream",
  "action": "stop",
  "streamId": "stream-456"
}
```

> **重要**：前端可以通过发送终止消息，请求后端停止当前正在进行的流式消息推送。终止消息必须包含正确的streamId以识别要终止的流。

#### 响应消息

**开始消息**:
```json
{
  "type": "stream",
  "action": "start",
  "sessionId": "session-123",
  "streamId": "stream-456",
  "timestamp": 1615453545
}
```

**内容消息**:
```json
{
  "type": "stream",
  "action": "message",
  "message": "你好，我是AI助手，很高兴为你服务！",
  "sessionId": "session-123",
  "streamId": "stream-456",
  "role": "assistant",
  "timestamp": 1615453546,
  "isFinal": false
}
```

**结束消息**:
```json
{
  "type": "stream",
  "action": "end",
  "sessionId": "session-123",
  "streamId": "stream-456",
  "timestamp": 1615453547,
  "isFinal": true
}
```

**错误消息**:
```json
{
  "type": "stream",
  "action": "error",
  "error": "AI 服务暂时不可用",
  "timestamp": 1615453547,
  "isFinal": true
}
```

### Protobuf 格式

使用与gRPC服务相同的Protobuf定义，通过WebSocket传输。

#### 请求消息 (StreamChatRequest)
```json
{
  "message": "你好，AI助手！",
  "stream_id": "stream-456",
  "context": {
    "session_id": "session-123",
    "history": ["用户之前的消息", "AI之前的回复"]
  }
}
```

#### 响应消息 (ChatResponse)
服务器发送一系列ChatResponse消息，最后一条的is_final为true。

#### 终止请求 (StopStreamRequest)
```json
{
  "stream_id": "stream-456"
}
```

> **说明**：当需要终止正在进行的流式推送时，客户端可发送终止请求。对于Protobuf格式，使用StopStreamRequest消息类型。

## 使用示例

### Vue 3应用中建立连接

```javascript
// 简化版连接示例
const connectWebSocket = () => {
  const token = localStorage.getItem('jwt_token');
  const socket = new WebSocket(`ws://localhost:9102/ws/chat?token=${encodeURIComponent(token)}`);
  
  socket.onopen = () => console.log('WebSocket连接已建立');
  
  socket.onmessage = (event) => {
    const response = JSON.parse(event.data);
    // 处理接收到的消息
  };
  
  socket.onerror = (error) => console.error('WebSocket错误:', error);
  socket.onclose = (event) => console.log(`连接已关闭，代码=${event.code}`);
  
  return socket;
};

// 发送消息示例
const sendMessage = (socket, message) => {
  const payload = {
    type: 'chat',
    action: 'message',
    message: message,
    sessionId: 'session-' + Date.now(),
    streamId: 'stream-' + Date.now(),
    context: { sessionId: 'session-' + Date.now() }
  };
  
  socket.send(JSON.stringify(payload));
};

// 发送终止消息示例
const stopStream = (socket, streamId) => {
  const stopMessage = {
    type: 'stream',
    action: 'stop',
    streamId: streamId
  };
  
  socket.send(JSON.stringify(stopMessage));
  console.log(`已请求终止流 ${streamId}`);
};
```

## 错误码说明

| 错误码 | 说明 |
|-------|------|
| STATUS_UNKNOWN | 未知状态 |
| STATUS_OK | 请求成功 |
| BAD_REQUEST | 请求参数错误 |
| UNAUTHORIZED | 认证失败 |
| FORBIDDEN | 权限不足 |
| INTERNAL_ERROR | 服务器内部错误 |
| SERVICE_UNAVAILABLE | 服务不可用 |

## 常见问题

### Q: 连接突然断开怎么办？
A: 实现指数退避重连策略，首次等待1秒，然后2秒，4秒等递增。

### Q: 如何处理大量历史消息？
A: 只发送最近的5-10条消息作为上下文，或使用API单独获取历史记录。

### Q: 如何终止正在进行的AI回复？
A: 当需要终止AI正在生成的回复（例如用户点击"停止生成"按钮）时，发送终止消息：
```javascript
// 使用之前记录的streamId发送终止消息
stopStream(socket, currentStreamId);
```
确保使用正确的streamId，这样后端才能准确识别并终止对应的流。终止后，后端通常会发送一个is_final=true的结束消息。

### Q: 为什么我的认证失败？
A: 可能原因：
1. HTTP头认证格式错误 - 必须为`Authorization: Bearer <token>`
2. URL参数名必须为`token` - 即`?token=<YOUR_JWT_TOKEN>`
3. 令牌已过期或签名无效
4. 令牌内容被截断或包含额外字符

### Q: 消息中的context字段是做什么用的？
A: context字段用于传递聊天上下文信息，而非认证信息。它包含：
1. sessionId - 当前会话的唯一标识符
2. history - 历史消息列表，帮助AI理解对话上下文

这些信息帮助AI服务理解对话的上下文，提供连贯的回复。认证信息仅在WebSocket连接建立阶段通过HTTP请求头或URL参数传递。