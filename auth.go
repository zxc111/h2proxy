package h2proxy

import (
	"net/http"
	"strings"
)

var authKeys = []string{
	"Proxy-Authorization",
}

// 检查请求中的 auth信息 和 用户信息是否一致
func CheckAuth(u *UserInfo, r *http.Request) bool {
	rightAuth := u.ToBase64()
	for auth, _ := range getAuthFromHeader(r) {
		//Log.Debug(auth)
		if auth == rightAuth {
			return true
		}
	}
	return false
}

func getAuthFromHeader(r *http.Request) map[string]struct{} {
	result := make(map[string]struct{}, 4)
	for _, k := range authKeys {
		auth := r.Header.Get(k)
		if auth == "" {
			continue
		}
		auth = strings.Replace(auth, "Basic ", "", 1)
		if auth == "" {
			continue
		}
		result[auth] = struct{}{}
	}

	return result
}

// 在 hesder 中设置 auth
func SetAuthInHeader(u *UserInfo, req *http.Request) {
	req.Header.Set("Proxy-Authorization", u.ToBase64())
}
