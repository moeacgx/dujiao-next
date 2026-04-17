package service

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/dujiao-next/internal/config"
	"github.com/dujiao-next/internal/constants"
	"github.com/dujiao-next/internal/models"
	"github.com/dujiao-next/internal/repository"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupRegistrationBlacklistTestService(t *testing.T) (*UserAuthService, *gorm.DB) {
	t.Helper()

	dsn := fmt.Sprintf("file:user_auth_service_registration_blacklist_%d?mode=memory&cache=shared", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.UserOAuthIdentity{}, &models.EmailVerifyCode{}, &models.Setting{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}

	cfg := &config.Config{
		UserJWT: config.JWTConfig{
			SecretKey:   "user-jwt-test-secret",
			ExpireHours: 24,
		},
	}
	settingSvc := NewSettingService(repository.NewSettingRepository(db))
	svc := NewUserAuthService(
		cfg,
		repository.NewUserRepository(db),
		repository.NewUserOAuthIdentityRepository(db),
		repository.NewEmailVerifyCodeRepository(db),
		settingSvc,
		nil,
		nil,
	)
	return svc, db
}

func TestRegisterRejectsBlacklistedEmailDomain(t *testing.T) {
	svc, db := setupRegistrationBlacklistTestService(t)
	if _, err := svc.settingService.Update(constants.SettingKeyOrderRiskControlConfig, map[string]interface{}{
		"enabled":                true,
		"email_domain_blacklist": []string{"example.com"},
	}); err != nil {
		t.Fatalf("update order risk control config failed: %v", err)
	}

	user, token, expiresAt, err := svc.Register("Bad@Example.com", "Abcd123!", "", true, false)
	if !errors.Is(err, ErrRiskEmailDomainBlacklisted) {
		t.Fatalf("expected ErrRiskEmailDomainBlacklisted, got user=%v token=%q expiresAt=%v err=%v", user, token, expiresAt, err)
	}

	var count int64
	if err := db.Model(&models.User{}).Count(&count).Error; err != nil {
		t.Fatalf("count users failed: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected no users created, got %d", count)
	}
}

func TestSendVerifyCodeRejectsBlacklistedEmailDomainOnRegister(t *testing.T) {
	svc, _ := setupRegistrationBlacklistTestService(t)
	if _, err := svc.settingService.Update(constants.SettingKeyOrderRiskControlConfig, map[string]interface{}{
		"enabled":                true,
		"email_domain_blacklist": []string{"@example.com"},
	}); err != nil {
		t.Fatalf("update order risk control config failed: %v", err)
	}

	err := svc.SendVerifyCode("spam@example.com", constants.VerifyPurposeRegister, constants.LocaleZhCN)
	if !errors.Is(err, ErrRiskEmailDomainBlacklisted) {
		t.Fatalf("expected ErrRiskEmailDomainBlacklisted, got %v", err)
	}
}
