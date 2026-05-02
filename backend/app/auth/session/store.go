package session

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

type SessionData struct {
	UserID    uint
	Username  string
	IsAdmin   bool
	CreatedAt time.Time
}

var (
	store     = sync.Map{}
	cookieName = "gcmdb_session"
)

func generateID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func Create(userID uint, username string, isAdmin bool) string {
	sid := generateID()
	store.Store(sid, &SessionData{
		UserID:    userID,
		Username:  username,
		IsAdmin:   isAdmin,
		CreatedAt: time.Now(),
	})
	return sid
}

func Get(sid string) *SessionData {
	val, ok := store.Load(sid)
	if !ok {
		return nil
	}
	return val.(*SessionData)
}

func Delete(sid string) {
	store.Delete(sid)
}

func CookieName() string {
	return cookieName
}
