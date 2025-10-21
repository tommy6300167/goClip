package services

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"clipmini/models"
	"clipmini/utils"
)

type AuthService struct {
	config *models.AppConfig
}

func NewAuthService(config *models.AppConfig) *AuthService {
	return &AuthService{
		config: config,
	}
}

func (as *AuthService) IsFirstLaunch() bool {
	_, err := os.Stat(as.config.AuthFilePath)
	return os.IsNotExist(err)
}

func (as *AuthService) getOrGenerateDeviceID() string {
	if authInfo, err := as.loadAuthInfo(); err == nil {
		return authInfo.DeviceID
	}

	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (as *AuthService) loadAuthInfo() (*models.AuthInfo, error) {
	data, err := os.ReadFile(as.config.AuthFilePath)
	if err != nil {
		return nil, err
	}

	var authInfo models.AuthInfo
	if err := json.Unmarshal(data, &authInfo); err != nil {
		return nil, err
	}

	return &authInfo, nil
}

func (as *AuthService) saveAuthInfo(authInfo *models.AuthInfo) error {
	if err := os.MkdirAll(as.config.LogDirPath, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(authInfo, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(as.config.AuthFilePath, data, 0600)
}

func (as *AuthService) callAuthAPI(licenseKey string) (*models.AuthResponse, error) {
	deviceID := as.getOrGenerateDeviceID()

	request := models.AuthRequest{
		LicenseKey: licenseKey,
		DeviceID:   deviceID,
		AppVersion: "1.0.0",
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 30 * time.Second}

	resp, err := client.Post(as.config.AuthApiEndpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var authResponse models.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		return nil, err
	}

	return &authResponse, nil
}

func (as *AuthService) ValidateKey(inputKey string) bool {
	println("DEBUG: ValidateKey 被調用，輸入金鑰:", inputKey)

	inputKey = strings.TrimSpace(inputKey)
	if inputKey == "" {
		println("DEBUG: 金鑰為空")
		return false
	}

	isFirstLaunch := as.IsFirstLaunch()
	println("DEBUG: IsFirstLaunch:", isFirstLaunch)

	if !isFirstLaunch {
		println("DEBUG: 不是初次啟動，檢查已存認證")
		return as.validateStoredAuth()
	}

	println("DEBUG: 初次啟動，調用API驗證")
	return as.validateWithAPI(inputKey)
}

func (as *AuthService) validateStoredAuth() bool {
	println("DEBUG: validateStoredAuth 被調用")

	authInfo, err := as.loadAuthInfo()
	if err != nil {
		println("DEBUG: 載入認證資訊失敗:", err.Error())
		return false
	}

	if !authInfo.IsValid() {
		println("DEBUG: 認證已過期")
		return false
	}

	if !utils.IsValidDeviceFingerprint(authInfo.DeviceFingerprint) {
		println("DEBUG: 設備指紋不匹配")
		return false
	}

	authInfo.LastValidate = time.Now()
	as.saveAuthInfo(authInfo)

	println("DEBUG: 已存認證驗證成功")
	return true
}

func (as *AuthService) validateWithAPI(licenseKey string) bool {
	return as.fallbackValidation(licenseKey)

	// 原本的API調用邏輯（暫時註解）
	/*
		response, err := as.callAuthAPI(licenseKey)
		if err != nil {
			// API無法連接時的備用驗證
			return as.fallbackValidation(licenseKey)
		}

		if !response.Success {
			return false
		}
	*/

	// 這部分邏輯已移到 fallbackValidation 中
	/*
		// 獲取設備指紋
		deviceFingerprint, err := utils.DeviceFingerprint()
		if err != nil {
			return false
		}

		// 保存認證資訊
		authInfo := &models.AuthInfo{
			PublicKey:         response.PublicKey,
			ExpiresAt:         time.Unix(response.ExpiresAt, 0),
			IssuedAt:          time.Now(),
			DeviceID:          as.getOrGenerateDeviceID(),
			DeviceFingerprint: deviceFingerprint,
			LastValidate:      time.Now(),
		}

		return as.saveAuthInfo(authInfo) == nil
	*/
}

func (as *AuthService) fallbackValidation(licenseKey string) bool {
	validTestKeys := []string{
		"test-key-2024",
		"superclip-test",
		"dev-license-001",
	}

	println("DEBUG: 輸入的金鑰:", licenseKey)

	for _, validKey := range validTestKeys {
		println("DEBUG: 比較金鑰:", validKey)
		if licenseKey == validKey {
			println("DEBUG: 金鑰匹配成功")

			deviceFingerprint, err := utils.DeviceFingerprint()
			if err != nil {
				println("DEBUG: 設備指紋生成失敗:", err.Error())
				return false
			}
			println("DEBUG: 設備指紋:", deviceFingerprint)

			authInfo := &models.AuthInfo{
				PublicKey:         "test-public-key-" + licenseKey,
				ExpiresAt:         time.Now().Add(30 * 24 * time.Hour),
				IssuedAt:          time.Now(),
				DeviceID:          as.getOrGenerateDeviceID(),
				DeviceFingerprint: deviceFingerprint,
				LastValidate:      time.Now(),
			}

			err = as.saveAuthInfo(authInfo)
			if err != nil {
				println("DEBUG: 保存認證失敗:", err.Error())
				return false
			}
			println("DEBUG: 認證保存成功")
			return true
		}
	}

	println("DEBUG: 沒有匹配的金鑰")
	return false
}

func (as *AuthService) SaveAuthConfig() error {
	if err := os.MkdirAll(as.config.LogDirPath, 0755); err != nil {
		return err
	}

	content := "initialized=true\n"
	return os.WriteFile(as.config.ConfigFilePath, []byte(content), 0644)
}

func (as *AuthService) GetAuthStatus() (*models.AuthInfo, error) {
	return as.loadAuthInfo()
}

func (as *AuthService) IsAuthValid() bool {
	authInfo, err := as.loadAuthInfo()
	if err != nil {
		return false
	}

	if !authInfo.IsValid() {
		return false
	}

	return utils.IsValidDeviceFingerprint(authInfo.DeviceFingerprint)
}

func (as *AuthService) GetTimeUntilExpiry() (time.Duration, error) {
	authInfo, err := as.loadAuthInfo()
	if err != nil {
		return 0, err
	}
	return time.Until(authInfo.ExpiresAt), nil
}

func (as *AuthService) ClearAuth() error {
	return os.Remove(as.config.AuthFilePath)
}

// TODO: 未來可以添加這個方法來透過API驗證金鑰
// func (as *AuthService) ValidateKeyWithAPI(inputKey string) (bool, error) {
//     // 實作API呼叫邏輯
//     return false, nil
// }
