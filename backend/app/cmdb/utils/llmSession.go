package utils

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type llmSession struct{}

var LLMSession = new(llmSession)

// Session 会话数据
type Session struct {
	Messages  []ChatMessage
	LastUsed  time.Time
	SystemMsg string
}

var (
	sessionStore = make(map[string]*Session)
	sessionMu    sync.RWMutex
)

const (
	sessionTTL      = 30 * time.Minute
	maxHistoryTurns = 10
)

// GetOrCreate 获取或创建会话
func (ls *llmSession) GetOrCreate(sessionId, systemPrompt string) *Session {
	sessionMu.Lock()
	defer sessionMu.Unlock()

	if s, ok := sessionStore[sessionId]; ok {
		s.LastUsed = time.Now()
		return s
	}

	s := &Session{
		Messages: []ChatMessage{
			{Role: "system", Content: systemPrompt},
		},
		LastUsed:  time.Now(),
		SystemMsg: systemPrompt,
	}
	sessionStore[sessionId] = s
	return s
}

// AddMessage 添加消息到会话
func (ls *llmSession) AddMessage(sessionId string, role, content string) {
	sessionMu.Lock()
	defer sessionMu.Unlock()

	s, ok := sessionStore[sessionId]
	if !ok {
		return
	}

	s.Messages = append(s.Messages, ChatMessage{Role: role, Content: content})
	s.LastUsed = time.Now()

	// 限制历史轮数：system(1) + user/assistant pairs(maxHistoryTurns*2)
	maxMessages := 1 + maxHistoryTurns*2
	if len(s.Messages) > maxMessages {
		// 保留 system message + 最近的对话
		truncated := make([]ChatMessage, 0, maxMessages)
		truncated = append(truncated, s.Messages[0]) // system message
		truncated = append(truncated, s.Messages[len(s.Messages)-maxMessages+1:]...)
		s.Messages = truncated
	}
}

// ClearSession 清空会话
func (ls *llmSession) ClearSession(sessionId string) {
	sessionMu.Lock()
	defer sessionMu.Unlock()
	delete(sessionStore, sessionId)
}

// CleanupExpiredSessions 清理过期会话（可定期调用）
func (ls *llmSession) CleanupExpiredSessions() {
	sessionMu.Lock()
	defer sessionMu.Unlock()

	now := time.Now()
	for id, s := range sessionStore {
		if now.Sub(s.LastUsed) > sessionTTL {
			delete(sessionStore, id)
		}
	}
}

// GenerateSessionId 生成会话 ID
func (ls *llmSession) GenerateSessionId() string {
	return uuid.New().String()
}
