package session

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
	"userManagement/model"
)

// Session 结构
type Session struct {
	User model.User
}

// Session 存储
var (
	store = make(map[string]Session)
	lock  sync.RWMutex
)

// 生成随机sessionID
func generateSessionID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

// CreateSession 创建Session
func CreateSession(w http.ResponseWriter, user model.User) {
	sessionID := generateSessionID()
	lock.Lock()
	store[sessionID] = Session{User: user}
	lock.Unlock()
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   3600,
	}
	http.SetCookie(w, &cookie)
}

// GetSession 获取Session
func GetSession(r *http.Request) (model.User, bool) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return model.User{}, false
	}
	lock.RLock()
	s, exists := store[cookie.Value]
	lock.RUnlock()
	if !exists {
		return model.User{}, false
	}
	return s.User, true
}

// DestroySession 销毁Session(退出登录)
func DestroySession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return
	}
	lock.Lock()
	delete(store, cookie.Value)
	lock.Unlock()
	//删除浏览器的cookie
	expired := http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, &expired)
}

// 更新当前回话中的用户信息
// UpdateSessionUser 更新当前会话中的用户信息

func UpdateSessionUser(w http.ResponseWriter, r *http.Request, user *model.User) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return
	}
	sessionID := cookie.Value

	// 更新全局 store 中的用户信息
	// 注意：这里需要加锁，确保并发安全
	lock.Lock()
	defer lock.Unlock()

	if _, ok := store[sessionID]; ok {
		store[sessionID] = Session{User: *user}
	}
}
