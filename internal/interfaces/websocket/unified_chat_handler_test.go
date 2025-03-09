package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// 测试用的服务器
type TestWebSocketServer struct {
	Server *httptest.Server
	URL    string
}

func NewTestWebSocketServer(t *testing.T, handler http.HandlerFunc) *TestWebSocketServer {
	server := httptest.NewServer(handler)
	url := "ws" + strings.TrimPrefix(server.URL, "http")

	return &TestWebSocketServer{
		Server: server,
		URL:    url,
	}
}

func (s *TestWebSocketServer) Connect(t *testing.T) *websocket.Conn {
	conn, _, err := websocket.DefaultDialer.Dial(s.URL, nil)
	require.NoError(t, err)
	return conn
}

func (s *TestWebSocketServer) Close() {
	s.Server.Close()
}

// Mock services for testing
type MockTokenService struct {
	UserID uint
	Roles  []string
	Error  error
}

func (m *MockTokenService) ValidateTokenFromContext(ctx context.Context) (entity.UID, []string, error) {
	return entity.UID(m.UserID), m.Roles, m.Error
}

type MockAIService struct {
	ChatResponse    *service.ChatResponse
	ChatError       error
	StreamResponses []*service.ChatResponse
	StreamError     error
}

func (m *MockAIService) Chat(ctx context.Context, userID string, message string, context *service.ChatContext) (*service.ChatResponse, error) {
	return m.ChatResponse, m.ChatError
}

func (m *MockAIService) StreamChat(ctx context.Context, userID string, message string, context *service.ChatContext) (<-chan *service.ChatResponse, error) {
	if m.StreamError != nil {
		return nil, m.StreamError
	}

	responseChan := make(chan *service.ChatResponse, len(m.StreamResponses))
	go func() {
		defer close(responseChan)
		for _, resp := range m.StreamResponses {
			responseChan <- resp
		}
	}()
	return responseChan, nil
}

func TestUnifiedChatHandler_HandleJSON(t *testing.T) {
	// Setup test logger
	zap.ReplaceGlobals(zap.NewExample())

	// Create mock services
	mockAI := &MockAIService{
		StreamResponses: []*service.ChatResponse{
			{
				Message:   "你好，我是AI助手",
				CreatedAt: time.Now().Unix(),
			},
		},
	}

	mockToken := &MockTokenService{
		UserID: 123,
		Roles:  []string{"user"},
	}

	// Create handler with mocks
	handler := &UnifiedChatHandler{
		aiService:    mockAI,
		tokenService: mockToken,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	// Create test server
	server := NewTestWebSocketServer(t, handler.HandleWebSocket)
	defer server.Close()

	// Connect to server
	conn := server.Connect(t)
	defer conn.Close()

	// Send test message
	testMsg := ChatMessage{
		Type:      "chat",
		Action:    "message",
		Message:   "你好，这是测试消息",
		SessionID: "test-session",
		StreamID:  "test-stream",
		Context:   &service.ChatContext{SessionID: "test-session"},
	}

	msgBytes, err := json.Marshal(testMsg)
	require.NoError(t, err)
	require.NoError(t, conn.WriteMessage(websocket.TextMessage, msgBytes))

	// Verify start message
	var startMsg ChatMessage
	err = conn.ReadJSON(&startMsg)
	require.NoError(t, err)
	assert.Equal(t, "stream", startMsg.Type)
	assert.Equal(t, "start", startMsg.Action)

	// Verify content message
	var contentMsg ChatMessage
	err = conn.ReadJSON(&contentMsg)
	require.NoError(t, err)
	assert.Equal(t, "stream", contentMsg.Type)
	assert.Equal(t, "message", contentMsg.Action)
	assert.Equal(t, "你好，我是AI助手", contentMsg.Message)

	// Verify end message
	var endMsg ChatMessage
	err = conn.ReadJSON(&endMsg)
	require.NoError(t, err)
	assert.Equal(t, "stream", endMsg.Type)
	assert.Equal(t, "end", endMsg.Action)
	assert.True(t, endMsg.IsFinal)
}
