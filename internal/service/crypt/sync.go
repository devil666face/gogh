package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"math/rand"
	"unsafe"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

func randString(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

func AES32RandomKey() []byte {
	return []byte(randString(32))
}

type Sync struct {
	AesKey []byte
}

func New(key ...string) (*Sync, error) {
	s := &Sync{
		AesKey: AES32RandomKey(),
	}
	if len(key) > 0 {
		if err := s.WithB64Key(key[0]); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (s *Sync) WithB64Key(b64key string) error {
	key, err := base64.StdEncoding.DecodeString(b64key)
	if err != nil {
		return fmt.Errorf("failed to set key: %w", err)
	}
	s.AesKey = key
	return nil
}

func (s *Sync) B64Key() string {
	return base64.StdEncoding.EncodeToString(s.AesKey)
}

func (s *Sync) Encrypt(data []byte) ([]byte, error) {
	aes, err := aes.NewCipher(s.AesKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = crand.Read(nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, data, nil), nil
}

func (s *Sync) Decrypt(data []byte) ([]byte, error) {
	aes, err := aes.NewCipher(s.AesKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return nil, err
	}

	// Since we know the ciphertext is actually nonce+ciphertext
	// And len(nonce) == NonceSize(). We can separate the two.
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plain, nil
}
