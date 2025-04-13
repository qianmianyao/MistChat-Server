package encryption

import (
	"crypto/hmac"
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"

	"time"
)

var secret = []byte("parchment_20250413") // 固定密钥

const (
	saltSize = 4
	tsSize   = 8
	signSize = 8
)

// GenerateUID 生成带签名的 UID，不含用户名信息
func GenerateUID(prefix string) (string, error) {
	salt := make([]byte, saltSize)
	if _, err := crand.Read(salt); err != nil {
		return "", err
	}

	ts := make([]byte, tsSize)
	binary.BigEndian.PutUint64(ts, uint64(time.Now().UnixNano()))

	// payload = salt + timestamp
	payload := append(salt, ts...)

	// 签名
	mac := hmac.New(sha256.New, secret)
	mac.Write(payload)
	signature := mac.Sum(nil)[:signSize]

	// 完整数据 = payload + signature
	full := append(payload, signature...)

	// 编码成 UID 字符串
	uid := base64.RawURLEncoding.EncodeToString(full)
	return prefix + uid, nil
}

func ValidateUID(uid, prefix string) (bool, error) {
	if len(uid) <= len(prefix) || uid[:len(prefix)] != prefix {
		return false, errors.New("invalid prefix")
	}

	raw, err := base64.RawURLEncoding.DecodeString(uid[len(prefix):])
	if err != nil {
		return false, err
	}

	if len(raw) < saltSize+tsSize+signSize {
		return false, errors.New("uid too short")
	}

	payload := raw[:saltSize+tsSize]
	sig := raw[saltSize+tsSize:]

	// 重新签名
	mac := hmac.New(sha256.New, secret)
	mac.Write(payload)
	expectedSig := mac.Sum(nil)[:signSize]

	if hmac.Equal(sig, expectedSig) {
		return true, nil
	}

	return false, errors.New("signature mismatch")
}
