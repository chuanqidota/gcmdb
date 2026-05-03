package cache

import (
	"gcmdb/config"
	"sync"
	"time"
)

type tokenEntry struct {
	UserID    uint
	Username  string
	IsAdmin   bool
	expiresAt time.Time
}

var tokenCache sync.Map

func getTTL() time.Duration {
	minutes := config.Conf.TokenCacheTTL
	if minutes <= 0 {
		minutes = 5
	}
	return time.Duration(minutes) * time.Minute
}

// GetTokenUser 从缓存获取用户信息，未命中返回 nil
func GetTokenUser(token string) *tokenEntry {
	if v, ok := tokenCache.Load(token); ok {
		entry := v.(*tokenEntry)
		if time.Now().Before(entry.expiresAt) {
			return entry
		}
		tokenCache.Delete(token)
	}
	return nil
}

// SetTokenUser 写入缓存
func SetTokenUser(token string, userID uint, username string, isAdmin bool) {
	tokenCache.Store(token, &tokenEntry{
		UserID:    userID,
		Username:  username,
		IsAdmin:   isAdmin,
		expiresAt: time.Now().Add(getTTL()),
	})
}

// InvalidateToken 使指定 token 缓存失效（用户登出/禁用时调用）
func InvalidateToken(token string) {
	tokenCache.Delete(token)
}
