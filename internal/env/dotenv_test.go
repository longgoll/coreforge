package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateRandomSecret(t *testing.T) {
	s1 := GenerateRandomSecret(32)
	s2 := GenerateRandomSecret(32)

	if len(s1) != 32 {
		t.Errorf("Expected length 32, got %d", len(s1))
	}
	if s1 == s2 {
		t.Error("Two random secrets should not be identical")
	}

	s64 := GenerateRandomSecret(64)
	if len(s64) != 64 {
		t.Errorf("Expected length 64, got %d", len(s64))
	}
}

func TestEnsureEnvFiles(t *testing.T) {
	tmpDir := t.TempDir()

	err := EnsureEnvFiles(tmpDir)
	if err != nil {
		t.Fatalf("EnsureEnvFiles failed: %v", err)
	}

	// Check .env exists
	envPath := filepath.Join(tmpDir, ".env")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		t.Error(".env should be created")
	}

	// Check .env.example exists
	examplePath := filepath.Join(tmpDir, ".env.example")
	if _, err := os.Stat(examplePath); os.IsNotExist(err) {
		t.Error(".env.example should be created")
	}

	// Calling again should not error (idempotent)
	err = EnsureEnvFiles(tmpDir)
	if err != nil {
		t.Fatalf("Second EnsureEnvFiles call failed: %v", err)
	}
}

func TestAppendEnvVars(t *testing.T) {
	tmpDir := t.TempDir()
	EnsureEnvFiles(tmpDir)

	vars := []EnvVarEntry{
		{Key: "JWT_SECRET", Value: "mysecret123", Description: "JWT signing key"},
		{Key: "PORT", Value: "3000", Description: "Server port"},
	}

	added, err := AppendEnvVars(tmpDir, vars)
	if err != nil {
		t.Fatalf("AppendEnvVars failed: %v", err)
	}
	if added != 2 {
		t.Errorf("Expected 2 vars added, got %d", added)
	}

	// Read .env and verify
	envContent, _ := os.ReadFile(filepath.Join(tmpDir, ".env"))
	envStr := string(envContent)

	if !strings.Contains(envStr, "JWT_SECRET=mysecret123") {
		t.Error(".env should contain JWT_SECRET")
	}
	if !strings.Contains(envStr, "PORT=3000") {
		t.Error(".env should contain PORT")
	}

	// Append again — should skip existing keys
	added2, err := AppendEnvVars(tmpDir, vars)
	if err != nil {
		t.Fatalf("Second AppendEnvVars failed: %v", err)
	}
	if added2 != 0 {
		t.Errorf("Expected 0 vars added (already exist), got %d", added2)
	}
}

func TestAppendEnvVars_NewKeysOnly(t *testing.T) {
	tmpDir := t.TempDir()
	EnsureEnvFiles(tmpDir)

	// Add first batch
	vars1 := []EnvVarEntry{
		{Key: "A", Value: "1"},
	}
	AppendEnvVars(tmpDir, vars1)

	// Add second batch with overlap
	vars2 := []EnvVarEntry{
		{Key: "A", Value: "changed"},  // should skip
		{Key: "B", Value: "2"},        // should add
	}
	added, _ := AppendEnvVars(tmpDir, vars2)
	if added != 1 {
		t.Errorf("Expected 1 new var added, got %d", added)
	}

	envContent, _ := os.ReadFile(filepath.Join(tmpDir, ".env"))
	envStr := string(envContent)

	// A should still have value "1" (not "changed")
	if strings.Contains(envStr, "A=changed") {
		t.Error("Existing key A should not be overwritten")
	}
	if !strings.Contains(envStr, "B=2") {
		t.Error("New key B should be added")
	}
}

func TestResolveValue(t *testing.T) {
	v1 := resolveValue("{{RANDOM_SECRET_32}}")
	if len(v1) != 32 {
		t.Errorf("RANDOM_SECRET_32 should produce 32 chars, got %d", len(v1))
	}

	v2 := resolveValue("{{RANDOM_SECRET_64}}")
	if len(v2) != 64 {
		t.Errorf("RANDOM_SECRET_64 should produce 64 chars, got %d", len(v2))
	}

	v3 := resolveValue("7d")
	if v3 != "7d" {
		t.Error("Non-placeholder should be returned as-is")
	}
}

func TestGetPlaceholder(t *testing.T) {
	p1 := getPlaceholder("JWT_SECRET", "{{RANDOM_SECRET_32}}")
	if p1 != "<your-secret-key-here>" {
		t.Errorf("Secret placeholder should be masked, got %s", p1)
	}

	p2 := getPlaceholder("PORT", "3000")
	if p2 != "3000" {
		t.Errorf("Non-secret should show default value, got %s", p2)
	}
}
