package user

import (
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
func authenticate(username string, password string) bool {
	for _, user := range userList {
		if user.Username == username {
			currPd, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			// first param should be hash password, second param should be original password
			err := bcrypt.CompareHashAndPassword(currPd, []byte(user.Password))
			return err == nil
		}
	}
	return false
}

// http://localhost:8081/login?username=hadoop&password=hadoop
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	session, err := sessionStore.Get(r, sessionCookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	isSuccess := authenticate(username, password)
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
