package logic

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type MyClaims struct {
	ID    int64
	State string `json:"state"`
	jwt.StandardClaims
}

func (l *LoginLogic) GetToken(id int64, state string) (string, string) {
	fmt.Println(l.svcCtx.Config.Auth.AccessSecret)
	fmt.Println(l.svcCtx.Config.Auth.RefreshSecret)
	var accessSecret = []byte(l.svcCtx.Config.Auth.AccessSecret)
	var refreshSecret = []byte(l.svcCtx.Config.Auth.RefreshSecret)
	// accessToken 的数据
	aT := MyClaims{
		id,
		state,
		jwt.StandardClaims{
			Issuer:    "AR",
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(3 * time.Minute).Unix(),
		},
	}
	// refreshToken 的数据
	rT := MyClaims{
		id,
		state,
		jwt.StandardClaims{
			Issuer:    "AR",
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
	}
	fmt.Println(333333333)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, aT)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rT)
	fmt.Println(4444444)

	accessTokenSigned, err := accessToken.SignedString(accessSecret)
	fmt.Println(55555555)

	if err != nil {
		fmt.Println("获取Token失败，Secret错误")
		return "", ""
	}
	fmt.Println(7)
	refreshTokenSigned, err := refreshToken.SignedString(refreshSecret)
	if err != nil {
		fmt.Println("获取Token失败，Secret错误")
		return "", ""
	}
	fmt.Println(accessTokenSigned)
	fmt.Println(9999)
	fmt.Println(refreshTokenSigned)
	return accessTokenSigned, refreshTokenSigned
}

func (l *LoginLogic) ParseToken(accessTokenString, refreshTokenString string) (*MyClaims, bool, error) {
	logx.Debug("In ParseToken")
	var accessSecret = []byte(l.svcCtx.Config.Auth.AccessSecret)
	var refreshSecret = []byte(l.svcCtx.Config.Auth.RefreshSecret)
	accessToken, err := jwt.ParseWithClaims(accessTokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return accessSecret, nil
	})
	if claims, ok := accessToken.Claims.(*MyClaims); ok && accessToken.Valid {
		return claims, false, nil
	}

	logx.Debug("RefreshToken")
	refreshToken, err := jwt.ParseWithClaims(refreshTokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return refreshSecret, nil
	})
	if err != nil {
		return nil, false, err
	}
	if claims, ok := refreshToken.Claims.(*MyClaims); ok && refreshToken.Valid {
		return claims, true, nil
	}

	return nil, false, errors.New("invalid token")
}
