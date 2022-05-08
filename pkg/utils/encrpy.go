package utils

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/wenzhenxi/gorsa"
)

func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// 公钥加密
func PublicEncrypt(data string, pubKey string) (string, error) {
	return rsaPubEncrypt(data, pubKey)
}

// 私钥解密
func PrivateDecrypt(cipherData, priKey string) (string, error) {
	return rsaPriDecrypt(cipherData, priKey)
}

// 私钥加密
func PrivateEncrypt(data string, priKey string) (string, error) {
	return rsaPriEncrypt(data, priKey)
}

// 公钥解密
func PublicDecrypt(cipherData, pubKey string) (string, error) {
	return rsaPubDecrypt(cipherData, pubKey)
}

// --------------RSA start
// 公钥加密
func rsaPubEncrypt(data string, pobulicKey string) (string, error) {
	pubenctypt, err := gorsa.PublicEncrypt(data, pobulicKey)
	if err != nil {
		return "", err
	} else {
		return pubenctypt, nil
	}
}

// 私钥解密
func rsaPriDecrypt(dataMd, privatekey string) (string, error) {
	pridecrypt, err := gorsa.PriKeyDecrypt(dataMd, privatekey)
	if err != nil {
		return "", err
	}
	return pridecrypt, nil
}

// 私钥加密
func rsaPriEncrypt(data string, priKey string) (string, error) {
	prienctypt, err := gorsa.PriKeyEncrypt(data, priKey)
	if err != nil {
		return "", err
	} else {
		return prienctypt, nil
	}
}

// 公钥解密
func rsaPubDecrypt(dataMd, pubKey string) (string, error) {
	pubdecrypt, err := gorsa.PublicDecrypt(dataMd, pubKey)
	if err != nil {
		return "", err
	} else {
		return pubdecrypt, nil
	}
}

// ---------------RSA end
