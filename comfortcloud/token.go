package comfortcloud

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Token struct {
	AccessToken          string `json:"access_token"`
	AccessTokenIssuedAt  int64  `json:"access_token_issued_at"`
	AccessTokenExpiresAt int64  `json:"access_token_expires_at"`
	RefreshToken         string `json:"refresh_token"`
	IDToken              string `json:"id_token"`
	ExpiresInSec         int    `json:"expires_in"`
	AccClientID          string `json:"acc_client_id"`
	Scope                string `json:"scope"`
}

func (t *Token) isValid() bool {
	if t == nil {
		return false
	}
	if t.AccessToken == "" {
		return false
	}
	parts := strings.Split(t.AccessToken, ".")
	if len(parts) != 3 {
		return false
	}

	expired, err := t.isAccessTokenExpired()
	if err != nil {
		return false
	}
	return !expired
}

func (t *Token) isAccessTokenExpired() (bool, error) {
	if t.AccessTokenExpiresAt == 0 || t.AccessTokenIssuedAt == 0 {
		err := t.setIATAndEXP()
		if err != nil {
			return false, fmt.Errorf("failed to set IAT and EXP: %s", err)
		}
	}
	now := time.Now().Unix()

	if now > t.AccessTokenExpiresAt {
		return true, nil
	}

	return false, nil
}

func (t *Token) setIATAndEXP() error {
	iat, exp, err := extractIATAndEXPFromJWT(t.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to extract IAT and EXP from access token: %w", err)
	}
	t.AccessTokenIssuedAt, t.AccessTokenExpiresAt = iat, exp
	return nil
}

func (t *Token) getAPIKey(timestamp time.Time) string {

	normalizedTime := time.Date(
		timestamp.Year(), timestamp.Month(), timestamp.Day(),
		timestamp.Hour(), timestamp.Minute(), timestamp.Second(),
		0, time.UTC) // Force UTC by leaving out timezone information

	// Convert to Unix timestamp in milliseconds
	timestampMs := fmt.Sprintf("%d", normalizedTime.UnixNano()/int64(time.Millisecond))

	components := []string{
		"Comfort Cloud",
		"521325fb2dd486bf4831b47644317fca",
		timestampMs,
		"Bearer ",
		t.AccessToken,
	}

	inputBuffer := strings.Join(components, "")
	hash := sha256.Sum256([]byte(inputBuffer))
	hashStr := hex.EncodeToString(hash[:])
	result := hashStr[:9] + "cfc" + hashStr[9:]
	return result
}

func extractIATAndEXPFromJWT(token string) (int64, int64, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return 0, 0, fmt.Errorf("invalid JWT token")
	}

	// Decode the payload (second part of the JWT)
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("error decoding JWT payload: %w", err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return 0, 0, fmt.Errorf("error parsing JWT JSON: %w", err)
	}

	iat, ok := payload["iat"].(float64)
	if !ok {
		return 0, 0, fmt.Errorf("iat not found or invalid")
	}
	exp, ok := payload["exp"].(float64)
	if !ok {
		return 0, 0, fmt.Errorf("exp not found or invalid")
	}

	return int64(iat), int64(exp), nil
}
