package util

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/base64"
)

const (
	base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

// EncodeMD5 md5 encryption
func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}

var coder = base64.NewEncoding(base64Table)
 
func Base64Encode(src string) string {
	return coder.EncodeToString([]byte(src))
}
 
func Base64Decode(src string) (string, error) {
	code, err := coder.DecodeString(src)
	return string(code), err
}
