package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHashBytes(t *testing.T) {
	hash1 := HashBytes([]byte("hello world"))
	hash2 := HashBytes([]byte("hello world"))
	hash3 := HashBytes([]byte("different content"))

	if hash1 != hash2 {
		t.Error("Same content should produce same hash")
	}
	if hash1 == hash3 {
		t.Error("Different content should produce different hash")
	}
	if len(hash1) != 64 {
		t.Errorf("SHA256 hex string should be 64 chars, got %d", len(hash1))
	}
}

func TestHashFile(t *testing.T) {
	// Create temp file
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	content := []byte("test file content")
	os.WriteFile(filePath, content, 0644)

	hash, err := HashFile(filePath)
	if err != nil {
		t.Fatalf("HashFile failed: %v", err)
	}

	expectedHash := HashBytes(content)
	if hash != expectedHash {
		t.Errorf("HashFile should match HashBytes: got %s, want %s", hash, expectedHash)
	}
}

func TestHashFile_NotFound(t *testing.T) {
	_, err := HashFile("/nonexistent/file.txt")
	if err == nil {
		t.Error("HashFile should return error for nonexistent file")
	}
}

func TestAddComponentWithHashes(t *testing.T) {
	cfg := &ForgeConfig{
		Language:  "nodejs",
		Framework: "express",
	}

	hashes := map[string]string{
		"/middlewares/auth.js": "abc123",
		"/services/token.js":  "def456",
	}

	cfg.AddComponentWithHashes("jwt-auth", hashes)

	if len(cfg.InstalledComponents) != 1 {
		t.Fatal("Should have 1 installed component")
	}

	comp := cfg.InstalledComponents[0]
	if comp.Name != "jwt-auth" {
		t.Errorf("Name should be jwt-auth, got %s", comp.Name)
	}
	if len(comp.FileHashes) != 2 {
		t.Errorf("Should have 2 file hashes, got %d", len(comp.FileHashes))
	}
	if comp.FileHashes["/middlewares/auth.js"] != "abc123" {
		t.Error("Hash mismatch for auth.js")
	}
}

func TestGetInstalledComponent(t *testing.T) {
	cfg := &ForgeConfig{
		InstalledComponents: []InstalledComponent{
			{Name: "jwt-auth", Version: "1.0.0", FileHashes: map[string]string{"/a.js": "hash1"}},
			{Name: "cors", Version: "1.0.0"},
		},
	}

	comp := cfg.GetInstalledComponent("jwt-auth")
	if comp == nil {
		t.Fatal("Should find jwt-auth")
	}
	if comp.FileHashes["/a.js"] != "hash1" {
		t.Error("Hash should be preserved")
	}

	comp2 := cfg.GetInstalledComponent("cors")
	if comp2 == nil {
		t.Fatal("Should find cors")
	}

	comp3 := cfg.GetInstalledComponent("nonexistent")
	if comp3 != nil {
		t.Error("Should return nil for nonexistent component")
	}
}

func TestFileHashesSerialization(t *testing.T) {
	cfg := &ForgeConfig{
		Language:      "nodejs",
		Framework:     "express",
		SourceDir:     "./src",
		ConfigVersion: "2.0",
	}

	hashes := map[string]string{
		"/middlewares/errorHandler.js": "abc123def456",
	}
	cfg.AddComponentWithHashes("error-handler", hashes)

	// Save
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	if err := Save(cfg); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(loaded.InstalledComponents) != 1 {
		t.Fatal("Should have 1 component after load")
	}

	comp := loaded.InstalledComponents[0]
	if len(comp.FileHashes) != 1 {
		t.Errorf("FileHashes should be preserved after serialization, got %d", len(comp.FileHashes))
	}
	if comp.FileHashes["/middlewares/errorHandler.js"] != "abc123def456" {
		t.Error("Hash value should be preserved")
	}
}
