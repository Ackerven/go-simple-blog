package api

import (
	"net/http"
	. "simple-blog/model"
	. "simple-blog/utils"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var user Account
	params := RequestJsonInterface(r)
	//断言
	if username, ok := params["username"].(string); !ok {
		panic(ERROR_USERNAME_TYPE_WRONG)
	} else {
		if username == "" {
			panic(ERROR_USERNAME_NOT_NULL)
		}
		user.Username = username
	}
	if password, ok := params["password"].(string); !ok {
		panic(ERROR_PASSWORD_TYPE_WRONG)
	} else {
		if password == "" {
			panic(ERROR_PASSWORD_NOT_NULL)
		}
		user.Password = password
	}

	status := CheckLogin(user.Username, user.Password)
	cookie := http.Cookie{
		Name: "login",
		Value: user.Username,
		Path: "/",
		Expires: time.Now().Add(24*time.Hour),
	}
	http.SetCookie(w, &cookie)
	w.Write(MapToBody(Map{
		"status":status,
		"desc":GetErrorMessage(status),
	}))
}