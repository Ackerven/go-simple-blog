package middleware

import "C"
import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"simple-blog/utils"
	"strings"
	"time"
)

var JwtKey = []byte(utils.JwtKey)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

//生成token
func SetToken(username string) (string,int) {
	expireTime := time.Now().Add(24*time.Hour)
	SetClaims := Claims{
		username,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer: "Ackerven",
		},
	}
	requestClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, SetClaims)
	token, err := requestClaim.SignedString(JwtKey)
	if err != nil {
		return "", utils.ERROR
	}
	return token, utils.SUCCESS
}

//验证token
func VerifyToken(token string) (*Claims,int) {
	setToken, _ := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if key, _ := setToken.Claims.(*Claims); setToken.Valid {
		return key, utils.SUCCESS
	} else {
		return nil, utils.ERROR
	}
}

//jwt
func JwtToken(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err == nil {
				return
			}
			switch result := err.(type) {
			case int :
				//status := err.(int)
				w.Write(utils.MapToBody(utils.Map{
					"status":result,
					"desc": utils.GetErrorMessage(result),
				}))
			default:
				fmt.Printf("系统错误：%v\n", result)
			}
		}()
		tokenHerder := r.Header.Get("Authorization")
		var status int = utils.SUCCESS
		if tokenHerder == ""{
			panic(utils.ERROR_TOKEN_TYPE_WRONG)
		}
		checkToken := strings.SplitN(tokenHerder, " ", 2)
		if len(checkToken) != 2 && checkToken[0] != "Bearer" {
			panic(utils.ERROR_TOKEN_TYPE_WRONG)
		}
		key, status := VerifyToken(checkToken[1])
		if status != utils.SUCCESS {
			panic(utils.ERROR_TOKEN_WRONG)
		}
		if time.Now().Unix() > key.ExpiresAt {
			panic(utils.ERROR_TOKEN_RUNTIME)
		}
		r.Header.Set("username", key.Username)
		h(w, r)
	}
}