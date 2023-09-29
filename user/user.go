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
