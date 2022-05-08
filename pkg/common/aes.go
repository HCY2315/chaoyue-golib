package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"git.cestong.com.cn/cecf/cecf-golib/pkg/errors"
	"io"
)

type GCMCipher struct {
	gcm cipher.AEAD
}

func NewGCMCipher(passwd []byte) (*GCMCipher, error) {
	c, err := aes.NewCipher(passwd)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	return &GCMCipher{
		gcm: gcm,
	}, nil
}

func (g *GCMCipher) Encrypt(plainData []byte) ([]byte, error) {
	nonce := make([]byte, g.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errors.Wrap(err, "read rand")
	}

	return g.gcm.Seal(nonce, nonce, plainData, nil), nil
}

func (g *GCMCipher) Decrypt(encryptData []byte) ([]byte, error) {
	nonceSize := g.gcm.NonceSize()
	if len(encryptData) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short(%d)", len(encryptData))
	}

	nonce, ciphertext := encryptData[:nonceSize], encryptData[nonceSize:]
	return g.gcm.Open(nil, nonce, ciphertext, nil)
}
