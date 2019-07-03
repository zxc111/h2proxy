package h2proxy

import (
	"bytes"
	"encoding/base64"
	"fmt"
)

type userInfo struct {
	username string
	passwd   string
}

func (u *userInfo) String() string {
	return fmt.Sprintf("username: %s, passwd: %s", u.username, u.passwd)
}

func (u *userInfo) ToBase64() string {
	if u.username == "" && u.passwd == "" {
		return ""
	}
	b := bytes.NewBuffer([]byte(u.username))
	b.WriteByte(':')
	b.WriteString(u.passwd)

	return base64.RawURLEncoding.EncodeToString(b.Bytes())
}
