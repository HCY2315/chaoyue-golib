// Package utils 提供常用小功能
package utils

import (
	"encoding/base64"
	"fmt"
)

// Base64 加密
func EncodeByBase64(data string) string {
	sEnc := base64.StdEncoding.EncodeToString([]byte(data))
	return sEnc
}

// Base64 解密
func DecodeByBase64(sEncDaata string) (string, error) {
	sDec, err := base64.StdEncoding.DecodeString(sEncDaata)
	if err != nil {
		fmt.Printf("Error decoding string: %s ", err.Error())
		return "", err
	}
	return string(sDec), nil
}
