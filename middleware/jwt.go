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
func VerifiToken(token string) (*Claims,int) {
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
//func JwtToken() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		tokenHerder := c.Request.Header.Get("Authorization")
//		var status int = utils.SUCCESS
//		if tokenHerder == "" {
//			status = utils.ERROR_TOKEN_NOT_EXIST
//			c.JSON(http.StatusOK, gin.H{
//				"status":status,
//				"desc":utils.GetErrorMessage(status),
//			})
//			c.Abort()
//			return
//		}
//		checkToken := strings.SplitN(tokenHerder, " ", 2)
//		if len(checkToken) != 2 && checkToken[0] != "Bearer" {
//			status = utils.ERROR_TOKEN_TYPE_WRONG
//			c.JSON(http.StatusOK, gin.H{
//				"status":status,
//				"desc":utils.GetErrorMessage(status),
//			})
//			c.Abort()
//			return
//		}
//		key, code := VerifiToken(checkToken[1])
//		if code == utils.ERROR {
//			status = utils.ERROR_TOKEN_WRONG
//			c.JSON(http.StatusOK, gin.H{
//				"status":status,
//				"desc":utils.GetErrorMessage(status),
//			})
//			c.Abort()
//			return
//		}
//		if time.Now().Unix() > key.ExpiresAt {
//			status = utils.ERROR_TOKEN_RUNTIME
//			c.JSON(http.StatusOK, gin.H{
//				"status":status,
//				"desc":utils.GetErrorMessage(status),
//			})
//			c.Abort()
//			return
//		}
//		c.Set("username", key.Username)
//		c.Next()
//	}
//}

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
		key, status := VerifiToken(checkToken[1])
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