package middleware

import (
	"net/http"
	"strings"
	"ten_module/internal/service/authservice"
)

type UseMiddleware struct {
}
type Middleware func(http.HandlerFunc) http.HandlerFunc

var Middlewares *UseMiddleware

func InitMiddleWare() {
	Middlewares = &UseMiddleware{}
}
func (middle *UseMiddleware) CheckToken() Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			TokenCheck := authservice.TokenHelper{}
			Token := strings.Split(r.Header.Get("Authorization"), " ")[2]
			if Token == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			err := TokenCheck.VerifyToken(Token)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			next(w, r)
		}
	}

}
func (middle *UseMiddleware) VerifyRole(RoleRequire []string) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			TokenCheck := authservice.TokenHelper{}
			Token := strings.Split(r.Header.Get("Authorization"), " ")[2]
			Role, err := TokenCheck.GetRoleToken(Token)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			checkRoleInRequire := make(map[string]int)
			for _, roleRequire := range RoleRequire {
				checkRoleInRequire[roleRequire]++
			}
			for _, role := range Role {
				_, checkExsits := checkRoleInRequire[role]
				if !checkExsits {
					http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
					return
				}
			}
			next(w, r)
		}
	}
}
func (middle *UseMiddleware) Chain(ApiFunc http.HandlerFunc, Middleware ...Middleware) http.HandlerFunc {
	for _, check := range Middleware {
		ApiFunc = check(ApiFunc)
	}
	return ApiFunc
}
