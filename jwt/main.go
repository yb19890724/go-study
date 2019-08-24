package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

//
type UserClaim struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func main() {
	// ha256 对称加密
	sec := []byte("123abc")
	jwtTokn := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim{Username: "test"})
	token, err := jwtTokn.SignedString(sec)
	fmt.Println(token,err)
	
	uc:=UserClaim{}
	getToken,err:=jwt.ParseWithClaims(token,&uc,func(token *jwt.Token) (i interface{}, e error) {
		return sec,nil
	})
	
	if getToken.Valid {//验证通过
		fmt.Println(getToken.Claims,uc.Username)
	}
	
}
