package h2proxy

import (
	"net/http"
	"testing"
)

func TestAuth(t *testing.T) {
	testCase := []struct {
		u      userInfo
		r      http.Request
		result bool
	}{

		{
			result: true,
		},
		{
			userInfo{
				"a",
				"b",
			},
			http.Request{},
			false,
		},

		{
			userInfo{
				"a",
				"b",
			},
			http.Request{},
			false,
		},
		{
			userInfo{
				"a",
				"b",
			},
			http.Request{Header: http.Header{"Proxy-Authenticate": []string{"123"}}},
			false,
		},
		{
			userInfo{
				"a",
				"b",
			},
			http.Request{Header: http.Header{"Proxy-Authenticate": []string{"YTpi"}}},
			true,
		},
		{
			userInfo{
				"a",
				"b",
			},
			http.Request{Header: http.Header{"Authorization": []string{"YTpi"}}},
			true,
		},
	}
	for _, Case := range testCase {
		if Case.result != CheckAuth(&Case.u, &Case.r) {
			t.Fatal(Case)
		}
	}
}
