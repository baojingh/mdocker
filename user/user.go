package user

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string
	Password string
}

var userList = []User{
	{Username: "hadoop", Password: "hadoop"},
}

var sessionStore *sessions.CookieStore
var sessionCookieName = "mdocker-user-session"

const (
	cookieStoreAuthKey    = "key"
	cookieStoreEncryptKey = "1234567891234567"
)

// Session存储的初始化工作
func init() {
	sessionStore = sessions.NewCookieStore(
		[]byte(cookieStoreAuthKey),
		[]byte(cookieStoreEncryptKey),
	)

	sessionStore.Options = &sessions.Options{
		HttpOnly: true,
		MaxAge:   60 * 10,
	}
}

// authenticate the user and password
func authenticate(username string, password string) (bool, error) {
	var err error
	for _, user := range userList {
		if user.Username == username {
			currPd, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			// first param should be hash password, second param should be original password
			err := bcrypt.CompareHashAndPassword(currPd, []byte(user.Password))
			return err == nil, err
		}
	}
	return false, err
}

// http://localhost:8081/login?username=hadoop&password=hadoop
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// process cross origin issues for frontend
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	req := &struct {
		Username string `json: "username"`
		Password string `json: "password"`
	}{}
	json.NewDecoder(r.Body).Decode(req)
	username := req.Username
	password := req.Password

	session, err := sessionStore.Get(r, sessionCookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	isSuccess, _ := authenticate(username, password)
	log.Infof("username: %s, password: %s, isSuccess: %v", username, password, isSuccess)
	if isSuccess {
		log.Infof("User %s login success", username)
	} else {
		log.Infof("User %s login failure", username)
	}

	// 在session中标记用户已经通过登录验证
	session.Values["authenticated"] = true

	_ = session.Save(r, w)
}

// 用户登出，会在 Session 中标记用户是未认证的
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessionStore.Get(r, sessionCookieName)
	session.Values["authenticated"] = false
	w.Write([]byte("logout success\n"))
	session.Save(r, w)
	log.Info("User logout success")
}

// 通过用户 Session 判断用户是否已认证，未认证返回 403 Forbidden 错误。
func SecretHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessionStore.Get(r, sessionCookieName)
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	log.Info(w)
}
