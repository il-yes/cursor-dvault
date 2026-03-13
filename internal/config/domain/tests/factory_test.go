package config_domain_tests

import (
	"testing"

	app_config_domain "vault-app/internal/config/domain"
)

func TestInitConfig(t *testing.T) {
	userID := "test-user"

	cfg, err := app_config_domain.InitConfig(userID)

	if err != nil {
		t.Fatalf("InitConfig returned error: %v", err)
	}

	if cfg == nil {
		t.Fatal("expected config, got nil")
	}

	if cfg.App.UserID != userID {
		t.Errorf("expected userID %s, got %s", userID, cfg.App.UserID)
	}

	if cfg.App.Branch != "main" {
		t.Errorf("expected branch main, got %s", cfg.App.Branch)
	}

	if cfg.App.VaultSettings.EncryptionScheme != "AES-256-GCM" {
		t.Errorf("unexpected encryption scheme: %s", cfg.App.VaultSettings.EncryptionScheme)
	}

	if cfg.User.ID != userID {
		t.Errorf("expected user ID %s, got %s", userID, cfg.User.ID)
	}

	if cfg.User.Role != "user" {
		t.Errorf("expected role user, got %s", cfg.User.Role)
	}
}

func TestInitConfigFromVault(t *testing.T) {
	userID := "test-user"
	vaultName := "test-vault"

	cfg, err := app_config_domain.InitConfigFromVault(userID, vaultName)

	if err != nil {
		t.Fatalf("InitConfigFromVault returned error: %v", err)
	}

	if cfg == nil {
		t.Fatal("expected config, got nil")
	}

	vault := cfg.Vaults

	if vault.VaultName != vaultName {
		t.Errorf("expected vault name %s, got %s", vaultName, vault.VaultName)
	}

	if vault.UserID != userID {
		t.Errorf("expected userID %s, got %s", userID, vault.UserID)
	}
}

func TestVaultFeatureFlags(t *testing.T) {
	cfg, _ := app_config_domain.InitConfigFromVault("user1", "vault1")

	flags := cfg.Vaults.Features

	if !flags.TracecoreEnabled {
		t.Error("tracecore should be enabled")
	}

	if !flags.CloudBackupEnabled {
		t.Error("cloud backup should be enabled")
	}

	if !flags.ThreatDetectionEnabled {
		t.Error("threat detection should be enabled")
	}

	if !flags.BrowserExtensionEnabled {
		t.Error("browser extension should be enabled")
	}

	if !flags.GitCLIEnabled {
		t.Error("git CLI should be enabled")
	}
}

func TestVaultSyncDefaults(t *testing.T) {
	cfg, _ := app_config_domain.InitConfigFromVault("user1", "vault1")

	sync := cfg.Vaults.Sync

	if !sync.AutoSync {
		t.Error("AutoSync should be enabled")
	}

	if sync.SyncIntervalSeconds != app_config_domain.DefaultSyncInterval {
		t.Errorf("unexpected sync interval: %d", sync.SyncIntervalSeconds)
	}

	if sync.MaxRetries != app_config_domain.DefaultMaxRetries {
		t.Errorf("unexpected max retries: %d", sync.MaxRetries)
	}
}

func TestVaultSecurityDefaults(t *testing.T) {
	cfg, _ := app_config_domain.InitConfigFromVault("user1", "vault1")

	sec := cfg.Vaults.Security

	if sec.AutoLockSeconds != app_config_domain.DefaultAutoLockTimeout {
		t.Errorf("unexpected auto lock timeout: %d", sec.AutoLockSeconds)
	}

	if sec.SessionTimeout != app_config_domain.DefaultSessionTimeout {
		t.Errorf("unexpected session timeout: %d", sec.SessionTimeout)
	}

	if !sec.RequireBiometric {
		t.Error("biometric should be required")
	}
}

func TestDeviceInitialization(t *testing.T) {
	cfg, _ := app_config_domain.InitConfigFromVault("user1", "vault1")

	if len(cfg.Devices) == 0 {
		t.Fatal("expected at least one device")
	}

	device := cfg.Devices[0]

	if device.DeviceName != app_config_domain.DeviceName {
		t.Errorf("expected device name %s, got %s", app_config_domain.DeviceName, device.DeviceName)
	}

	if device.DeviceID == "" {
		t.Error("device ID should not be empty")
	}
}




func TestVaultCreationCases(t *testing.T) {

	tests := []struct {
		userID    string
		vaultName string
	}{
		{"user1", "vaultA"},
		{"user2", "vaultB"},
		{"user3", "vaultC"},
	}

	for _, tt := range tests {

		cfg, err := app_config_domain.InitConfigFromVault(tt.userID, tt.vaultName)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Vaults.VaultName != tt.vaultName {
			t.Errorf("expected %s got %s", tt.vaultName, cfg.Vaults.VaultName)
		}
	}
}