package encryption

import (
	"crypto/hmac"
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"github.com/btcsuite/btcutil/base58"
	"time"
)

var secret = []byte("parchment_20250413") // 固定密钥

const (
	saltSize = 4
	tsSize   = 8
	signSize = 8
)

// GenerateUID 使用 Base58 编码生成 UID
func GenerateUID(prefix string) (string, error) {
	salt := make([]byte, saltSize)
	if _, err := crand.Read(salt); err != nil {
		return "", err
	}

	ts := make([]byte, tsSize)
	binary.BigEndian.PutUint64(ts, uint64(time.Now().UnixNano()))

	payload := append(salt, ts...)

	mac := hmac.New(sha256.New, secret)
	mac.Write(payload)
	signature := mac.Sum(nil)[:signSize]

	full := append(payload, signature...)
	uid := base58.Encode(full)
	return prefix + uid, nil
}

// ValidateUID 校验 Base58 编码的 UID
func ValidateUID(uid, prefix string) (bool, error) {
	if len(uid) <= len(prefix) || uid[:len(prefix)] != prefix {
		return false, errors.New("invalid prefix")
	}

	raw := base58.Decode(uid[len(prefix):])
	if len(raw) < saltSize+tsSize+signSize {
		return false, errors.New("uid too short")
	}

	payload := raw[:saltSize+tsSize]
	sig := raw[saltSize+tsSize:]

	mac := hmac.New(sha256.New, secret)
	mac.Write(payload)
	expectedSig := mac.Sum(nil)[:signSize]

	if hmac.Equal(sig, expectedSig) {
		return true, nil
	}
	return false, errors.New("signature mismatch")
}
