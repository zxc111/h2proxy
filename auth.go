package h2proxy

import (
	"net/http"
	"strings"
)

// 检查请求中的 auth信息 和 用户信息是否一致
func CheckAuth(u *UserInfo, r *http.Request) bool {
	rightAuth := u.ToBase64()
	for _, auth := range getAuthFromHeader(r) {
		if strings.Replace(auth, "Basic ", "", 1) == rightAuth {
			return true
		}
	}
	return false
}

func getAuthFromHeader(r *http.Request) []string {
	result := make([]string, 0, 3)

	proxyAuth := r.Header.Get("Proxy-Authenticate")
	result = append(result, proxyAuth)
	normalAuth := r.Header.Get("Authorization")
	if normalAuth != proxyAuth {
		result = append(result, normalAuth)
	}
	wwwAuth := r.Header.Get("WWW-Authenticate")
	if wwwAuth != proxyAuth && wwwAuth != normalAuth {
		result = append(result, wwwAuth)
	}
	return result
}

// 在 hesder 中设置 auth
func SetAuthInHeader(u *UserInfo, req *http.Request) {
	req.Header.Set("Proxy-Authenticate", u.ToBase64())
}
