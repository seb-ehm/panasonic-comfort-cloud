package comfortcloud

import (
	"crypto/rand"
	"math/big"
)

type Authentication struct {
	username   string
	password   string
	token      *Token
	raw        bool
	appVersion string
}

func NewAuthentication(username, password string, token *Token, raw bool) *Authentication {
	return &Authentication{
		username:   username,
		password:   password,
		token:      token,
		raw:        raw,
		appVersion: XAppVersion,
	}
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}
	return string(result)
}
