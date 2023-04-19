package middleware

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"strings"
)

type JWTMiddleware struct {
}

func NewJWTMiddleware() *JWTMiddleware {
	return &JWTMiddleware{}
}

func (m *JWTMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	logx.Info("jwt middleware")
	fmt.Println(11111)
	return func(w http.ResponseWriter, r *http.Request) {
		// JWTAuthMiddleware implementation
		fmt.Println(000000)
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println("一错")
			return
		}
		fmt.Println(22222)
		parts := strings.Split(authHeader, " ")
		if !(len(parts) == 3 && parts[0] == "Bearer") {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println("二错")

			return
		}
		fmt.Println(333333)
		parseToken, isUpd, err := ParseToken(parts[1], parts[2])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println("三错")

			return
		}
		if isUpd {
			parts[1], parts[2] = GetToken(parseToken.ID, parseToken.State)
			w.Header().Set("Authorization", fmt.Sprintf("Bearer %s %s", parts[1], parts[2]))
		}
		fmt.Println(44444)
		r = r.WithContext(context.WithValue(r.Context(), "userID", parseToken.ID))
		fmt.Println(parseToken.ID, "gggggg")
		fmt.Println("faefewfwqg3俄方无法")
		// Passthrough to next handler
		next(w, r)
	}
}
