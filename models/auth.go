package models

import (
	"time"
)

type AuthInfo struct {
	PublicKey         string    `json:"public_key"`
	ExpiresAt         time.Time `json:"expires_at"`
	IssuedAt          time.Time `json:"issued_at"`
	DeviceID          string    `json:"device_id"`
	DeviceFingerprint string    `json:"device_fingerprint"`
	LastValidate      time.Time `json:"last_validate"`
}

func (ai *AuthInfo) IsValid() bool {
	return time.Now().Before(ai.ExpiresAt)
}

// IsValidOnThisDevice 檢查認證是否在當前設備有效
func (ai *AuthInfo) IsValidOnThisDevice() bool {
	if !ai.IsValid() {
		return false
	}

	// 需要匯入 utils 包來檢查設備指紋
	return true // 暫時返回 true，稍後在 auth_service 中實現
}

func (ai *AuthInfo) IsExpiringSoon() bool {
	return time.Until(ai.ExpiresAt) < 7*24*time.Hour
}

type AuthRequest struct {
	LicenseKey string `json:"license_key"`
	DeviceID   string `json:"device_id"`
	AppVersion string `json:"app_version"`
}

type AuthResponse struct {
	Success   bool   `json:"success"`
	PublicKey string `json:"public_key"`
	ExpiresAt int64  `json:"expires_at"` // Unix timestamp
	Message   string `json:"message"`
	ErrorCode string `json:"error_code,omitempty"`
}
