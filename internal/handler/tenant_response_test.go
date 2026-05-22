package handler

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Tencent/WeKnora/internal/types"
)

func TestSafeTenantForResponseStripsSecretsAndConfigs(t *testing.T) {
	tenant := &types.Tenant{
		ID:     1,
		Name:   "shared",
		APIKey: "sk-secret",
		ParserEngineConfig: &types.ParserEngineConfig{
			MinerUAPIKey: "mineru-secret",
		},
		StorageEngineConfig: &types.StorageEngineConfig{
			MinIO: &types.MinIOEngineConfig{SecretAccessKey: "storage-secret"},
		},
		Credentials: &types.CredentialsConfig{
			WeKnoraCloud: &types.WeKnoraCloudCredentials{AppID: "app", AppSecret: "app-secret"},
		},
	}

	safe := safeTenantForResponse(tenant)
	raw, err := json.Marshal(safe)
	if err != nil {
		t.Fatalf("marshal safe tenant: %v", err)
	}
	body := string(raw)
	for _, secret := range []string{"sk-secret", "mineru-secret", "storage-secret", "app-secret"} {
		if strings.Contains(body, secret) {
			t.Fatalf("safe tenant response leaked %q in %s", secret, body)
		}
	}

	if tenant.APIKey == "" || tenant.ParserEngineConfig == nil || tenant.StorageEngineConfig == nil || tenant.Credentials == nil {
		t.Fatalf("safeTenantForResponse mutated the source tenant")
	}
}
