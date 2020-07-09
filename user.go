package h2proxy

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"
)

type UserInfo struct {
	Username string
	Passwd   string
}

func (u *UserInfo) String() string {
	return fmt.Sprintf("Username: %s, Passwd: %s", u.Username, u.Passwd)
}

func (u *UserInfo) ToBase64() string {
	if u.Username == "" && u.Passwd == "" {
		return ""
	}
	b := bytes.NewBuffer([]byte(u.Username))
	b.WriteByte(':')
	b.WriteString(u.Passwd)
	return strings.ReplaceAll(base64.URLEncoding.EncodeToString(b.Bytes()), "=", "")
}
