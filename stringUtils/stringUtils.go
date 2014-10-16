package stringUtils

import (
	"crypto/md5"
	"fmt"
	"strings"
	"encoding/base64"
)


var coder = base64.StdEncoding

func Base64Encode(str string) string {
	return coder.EncodeToString([]byte(str))
}

func Base64Decode(src []byte) (string, error) {
	b,err:=coder.DecodeString(string(src))
	return string(b),err
}

func TrimRightSpace(s string) string {
	return strings.TrimRight(string(s), "\r\n\t ")
}

func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}
