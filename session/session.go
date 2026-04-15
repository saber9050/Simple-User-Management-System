package session

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"goweb/config"
	"goweb/myredis"
	"log"
	"net/http"
	"time"
)

type Session struct {
	ID       string    //会话id
	UserID   int       //用户id
	UserRole string    //用户角色
	ExpireAt time.Time //过期时间(24小时)
}

// 创建会话,同时返回会话id
func CreateSession(w http.ResponseWriter, userID int, role string) string {
	//	创建随机二进制字节数据
	b := make([]byte, 16)
	rand.Read(b)
	//	转换成32位16进制字符串，结果类似于镜像id，得到会话id
	sessionID := hex.EncodeToString(b)

	sess := &Session{
		ID:       sessionID,
		UserID:   userID,
		UserRole: role,
		ExpireAt: time.Now().Add(config.SessionExpiry),
	}

	//	缓存到redis
	saveSession, err := json.Marshal(*sess)
	if err != nil {
		log.Println("CreateSession会话序列化失败", err)
		return ""
	}
	//	先缓存用户会话映射关系	用户ID:sessionID
	userSession := fmt.Sprintf(config.PrefixUser, userID)
	err = myredis.RedisClient.Set(myredis.Ctx, userSession, sessionID, config.SessionExpiry).Err()
	if err != nil {
		log.Println("CreateSession用户会话映射关系设置失败", err)
		return ""
	}
	//	缓存session	sessionID对应session结构体
	err = myredis.RedisClient.Set(myredis.Ctx, sessionID, saveSession, config.SessionExpiry).Err()
	if err != nil {
		//	删除之前设置的映射关系
		myredis.RedisClient.Del(myredis.Ctx, userSession)
		log.Println("CreateSession会话设置失败", err)
		return ""
	}

	//	设置Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		HttpOnly: true,
		Path:     "/",
		Expires:  sess.ExpireAt,
	})
	return sessionID
}

// 从请求中获取会话内容
func GetSession(r *http.Request) *Session {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil
	}

	//	从redis读取会话
	var curSession Session
	res, err := myredis.RedisClient.Get(myredis.Ctx, cookie.Value).Result()
	if res == "" || err != nil {
		return nil
	}
	_ = json.Unmarshal([]byte(res), &curSession)

	return &curSession
}

// 获取对应用户的sessionID
func GetSessionByID(userID int) string {
	userSession := fmt.Sprintf(config.PrefixUser, userID)
	sessionID, err := myredis.RedisClient.Get(myredis.Ctx, userSession).Result()
	if err != nil || sessionID == "" {
		return ""
	}
	return sessionID
}

// 销毁会话，当登出时使用
func DestroySession(w http.ResponseWriter, r *http.Request) {
	sess := GetSession(r)
	if sess != nil {
		userSession := fmt.Sprintf(config.PrefixUser, sess.UserID)
		myredis.RedisClient.Del(myredis.Ctx, userSession)
		myredis.RedisClient.Del(myredis.Ctx, sess.ID)
	}
	//	删除cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
}

// 删除指定用户的会话（管理员专用）
func DestroyUserSession(userID int) error {
	userSession := fmt.Sprintf(config.PrefixUser, userID)
	sessionID, err := myredis.RedisClient.Get(myredis.Ctx, userSession).Result()
	//	删除映射键
	err = myredis.RedisClient.Del(myredis.Ctx, userSession).Err()
	if err != nil {
		return err
	}
	//	删除session
	if sessionID != "" {
		err = myredis.RedisClient.Del(myredis.Ctx, sessionID).Err()
	}
	return err
}

// 根据id更新用户角色(更新redis中缓存的session内容)
func UpdateUserRole(userID int, newRole string) error {
	userSession := fmt.Sprintf(config.PrefixUser, userID)
	sesssionID, _ := myredis.RedisClient.Get(myredis.Ctx, userSession).Result()
	if sesssionID != "" {
		//	如果会话存在
		tempsess, _ := myredis.RedisClient.Get(myredis.Ctx, sesssionID).Result()
		var sess Session
		_ = json.Unmarshal([]byte(tempsess), &sess)
		sess.UserRole = newRole
		data, _ := json.Marshal(sess)
		err := myredis.RedisClient.Set(myredis.Ctx, userSession, sesssionID, config.SessionExpiry).Err()
		if err != nil {
			return err
		}
		err = myredis.RedisClient.Set(myredis.Ctx, sesssionID, data, config.SessionExpiry).Err()
		return err
	}
	return nil
}
