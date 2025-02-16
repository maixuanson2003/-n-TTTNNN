package middleware

import (
	"net/http"
	"strings"
	"ten_module/internal/service/authservice"
)

type UseMiddleware struct {
}
type Middleware func(http.HandlerFunc) http.HandlerFunc

var middleware *UseMiddleware

func InitMiddleWare() {
	middleware = &UseMiddleware{}
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
			checkRoleInToken := make(map[string]int)
			for _, role := range Role {
				checkRoleInToken[role]++
			}
			for _, roleRequire := range RoleRequire {
				_, checkExsits := checkRoleInToken[roleRequire]
				if !checkExsits {
					http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
					return
				}
			}
			next(w, r)
		}
	}
}
func (middle *UseMiddleware) chain(ApiFunc http.HandlerFunc, Middleware ...Middleware) http.HandlerFunc {
	for _, check := range Middleware {
		ApiFunc = check(ApiFunc)
	}
	return ApiFunc
}
