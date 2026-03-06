package config

import (
	"encoding/json"
	"os"
	"testing"
)

// ──────────────────────────────────────────────────────────────
// Test Helpers
// ──────────────────────────────────────────────────────────────

// setupTestDir creates a temporary directory and changes to it.
// Returns a cleanup function to restore the original directory.
func setupTestDir(t *testing.T) func() {
	t.Helper()

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal("failed to get working directory:", err)
	}

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal("failed to change to temp directory:", err)
	}

	return func() {
		os.Chdir(origDir)
	}
}

// writeTestConfig writes a ForgeConfig as .forge.json in the current directory.
func writeTestConfig(t *testing.T, cfg *ForgeConfig) {
	t.Helper()
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatal("failed to marshal config:", err)
	}
	if err := os.WriteFile(configFileName, data, 0644); err != nil {
		t.Fatal("failed to write config:", err)
	}
}

// writeRawJSON writes raw JSON bytes as .forge.json in the current directory.
func writeRawJSON(t *testing.T, jsonStr string) {
	t.Helper()
	if err := os.WriteFile(configFileName, []byte(jsonStr), 0644); err != nil {
		t.Fatal("failed to write raw JSON:", err)
	}
}

// ──────────────────────────────────────────────────────────────
// Tests: Exists
// ──────────────────────────────────────────────────────────────

func TestExists_NoFile(t *testing.T) {
	cleanup := setupTestDir(t)
	defer cleanup()

	if Exists() {
		t.Error("Exists() should return false when no .forge.json")
	}
}

func TestExists_WithFile(t *testing.T) {
	cleanup := setupTestDir(t)
	defer cleanup()

	writeRawJSON(t, `{"language":"nodejs"}`)

	if !Exists() {
		t.Error("Exists() should return true when .forge.json exists")
	}
}

// ──────────────────────────────────────────────────────────────
// Tests: Load V1 Config
// ──────────────────────────────────────────────────────────────

func TestLoadV1Config(t *testing.T) {
	cleanup := setupTestDir(t)
	defer cleanup()

	// V1 config: no configVersion, database, orm, architecture fields
	writeRawJSON(t, `{
		"language": "nodejs",
		"framework": "express",
		"sourceDir": "./src",
		"installedComponents": [
			{"name": "error-handler", "version": "1.0.0", "installedAt": "2026-01-01T00:00:00Z"}
		]
	}`)

	cfg, err := Load()
	if err != nil {
		t.Fatal("Load() failed:", err)
	}

	// V1 fields should be present
	if cfg.Language != "nodejs" {
		t.Errorf("Language = %q, want %q", cfg.Language, "nodejs")
	}
	if cfg.Framework != "express" {
		t.Errorf("Framework = %q, want %q", cfg.Framework, "express")
	}
	if cfg.SourceDir != "./src" {
		t.Errorf("SourceDir = %q, want %q", cfg.SourceDir, "./src")
	}
	if len(cfg.InstalledComponents) != 1 {
		t.Fatalf("InstalledComponents count = %d, want 1", len(cfg.InstalledComponents))
	}
	if cfg.InstalledComponents[0].Name != "error-handler" {
		t.Errorf("InstalledComponents[0].Name = %q, want %q", cfg.InstalledComponents[0].Name, "error-handler")
	}

	// Auto-migration: V1 config gets ConfigVersion = "1.0"
	if cfg.ConfigVersion != "1.0" {
		t.Errorf("ConfigVersion = %q, want %q (auto-migrated from v1)", cfg.ConfigVersion, "1.0")
	}

	// V2 fields should be empty
	if cfg.Database != "" {
		t.Errorf("Database = %q, want empty string", cfg.Database)
	}
	if cfg.ORM != "" {
		t.Errorf("ORM = %q, want empty string", cfg.ORM)
	}
	if cfg.Architecture != "" {
		t.Errorf("Architecture = %q, want empty string", cfg.Architecture)
	}
}

// ──────────────────────────────────────────────────────────────
// Tests: Load V2 Config
// ──────────────────────────────────────────────────────────────

func TestLoadV2Config(t *testing.T) {
	cleanup := setupTestDir(t)
	defer cleanup()

	writeRawJSON(t, `{
		"language": "golang",
		"framework": "gin",
		"sourceDir": "./",
		"installedComponents": [],
		"configVersion": "2.0",
		"database": "postgresql",
		"orm": "gorm",
		"architecture": "mvc"
	}`)

	cfg, err := Load()
	if err != nil {
		t.Fatal("Load() failed:", err)
	}

	if cfg.Language != "golang" {
		t.Errorf("Language = %q, want %q", cfg.Language, "golang")
	}
	if cfg.ConfigVersion != "2.0" {
		t.Errorf("ConfigVersion = %q, want %q", cfg.ConfigVersion, "2.0")
	}
	if cfg.Database != "postgresql" {
		t.Errorf("Database = %q, want %q", cfg.Database, "postgresql")
	}
	if cfg.ORM != "gorm" {
		t.Errorf("ORM = %q, want %q", cfg.ORM, "gorm")
	}
	if cfg.Architecture != "mvc" {
		t.Errorf("Architecture = %q, want %q", cfg.Architecture, "mvc")
	}
}

// ──────────────────────────────────────────────────────────────
// Tests: Save & Roundtrip
// ──────────────────────────────────────────────────────────────

func TestSaveV2Config(t *testing.T) {
	cleanup := setupTestDir(t)
	defer cleanup()

	cfg := &ForgeConfig{
		Language:            "csharp",
		Framework:           "dotnet-webapi",
		SourceDir:           "./src",
		InstalledComponents: []InstalledComponent{},
		ConfigVersion:       "2.0",
		Database:            "postgresql",
		ORM:                 "efcore",
		Architecture:        "clean",
	}

	if err := Save(cfg); err != nil {
		t.Fatal("Save() failed:", err)
	}

	// Read back
	loaded, err := Load()
	if err != nil {
		t.Fatal("Load() after Save() failed:", err)
	}

	if loaded.Language != cfg.Language {
		t.Errorf("Language = %q, want %q", loaded.Language, cfg.Language)
	}
	if loaded.Database != cfg.Database {
		t.Errorf("Database = %q, want %q", loaded.Database, cfg.Database)
	}
	if loaded.ORM != cfg.ORM {
		t.Errorf("ORM = %q, want %q", loaded.ORM, cfg.ORM)
	}
	if loaded.Architecture != cfg.Architecture {
		t.Errorf("Architecture = %q, want %q", loaded.Architecture, cfg.Architecture)
	}
}

func TestSaveV2OmitsEmpty(t *testing.T) {
	cleanup := setupTestDir(t)
	defer cleanup()

	cfg := &ForgeConfig{
		Language:            "nodejs",
		Framework:           "express",
		SourceDir:           "./src",
		InstalledComponents: []InstalledComponent{},
		// V2 fields intentionally empty
	}

	if err := Save(cfg); err != nil {
		t.Fatal("Save() failed:", err)
	}

	// Read raw JSON and verify V2 fields are omitted
	data, err := os.ReadFile(configFileName)
	if err != nil {
		t.Fatal("ReadFile failed:", err)
	}

	var raw map[string]interface{}
	json.Unmarshal(data, &raw)

	// V2 fields with omitempty should not appear
	if _, found := raw["database"]; found {
		t.Error("database field should be omitted when empty")
	}
	if _, found := raw["orm"]; found {
		t.Error("orm field should be omitted when empty")
	}
	if _, found := raw["architecture"]; found {
		t.Error("architecture field should be omitted when empty")
	}
}

// ──────────────────────────────────────────────────────────────
// Tests: V1 → V2 Auto-Migration
// ──────────────────────────────────────────────────────────────

func TestV1ToV2Migration(t *testing.T) {
	cleanup := setupTestDir(t)
	defer cleanup()

	// Write a pure V1 config (no configVersion field at all)
	writeRawJSON(t, `{
		"language": "nodejs",
		"framework": "express",
		"sourceDir": "./src",
		"installedComponents": []
	}`)

	cfg, err := Load()
	if err != nil {
		t.Fatal("Load() failed:", err)
	}

	// Should auto-set configVersion to "1.0"
	if cfg.ConfigVersion != "1.0" {
		t.Errorf("ConfigVersion = %q, want %q after migration", cfg.ConfigVersion, "1.0")
	}

	// V2 fields should be zero-valued (not filled with defaults)
	if cfg.Database != "" {
		t.Errorf("Database should be empty after migration, got %q", cfg.Database)
	}
}

// ──────────────────────────────────────────────────────────────
// Tests: HasComponent
// ──────────────────────────────────────────────────────────────

func TestHasComponent(t *testing.T) {
	cfg := &ForgeConfig{
		InstalledComponents: []InstalledComponent{
			{Name: "error-handler", Version: "1.0.0"},
			{Name: "jwt-auth", Version: "1.0.0"},
		},
	}

	tests := []struct {
		name string
		want bool
	}{
		{"error-handler", true},
		{"jwt-auth", true},
		{"cors", false},
		{"", false},
	}

	for _, tt := range tests {
		got := cfg.HasComponent(tt.name)
		if got != tt.want {
			t.Errorf("HasComponent(%q) = %v, want %v", tt.name, got, tt.want)
		}
	}
}

// ──────────────────────────────────────────────────────────────
// Tests: AddComponent
// ──────────────────────────────────────────────────────────────

func TestAddComponent(t *testing.T) {
	cfg := &ForgeConfig{
		InstalledComponents: []InstalledComponent{},
	}

	cfg.AddComponent("cors")

	if len(cfg.InstalledComponents) != 1 {
		t.Fatalf("InstalledComponents count = %d, want 1", len(cfg.InstalledComponents))
	}
	if cfg.InstalledComponents[0].Name != "cors" {
		t.Errorf("InstalledComponents[0].Name = %q, want %q", cfg.InstalledComponents[0].Name, "cors")
	}
	if cfg.InstalledComponents[0].Version != "1.0.0" {
		t.Errorf("InstalledComponents[0].Version = %q, want %q", cfg.InstalledComponents[0].Version, "1.0.0")
	}
	if cfg.InstalledComponents[0].InstalledAt == "" {
		t.Error("InstalledComponents[0].InstalledAt should not be empty")
	}
}

// ──────────────────────────────────────────────────────────────
// Tests: RemoveComponent
// ──────────────────────────────────────────────────────────────

func TestRemoveComponent(t *testing.T) {
	cfg := &ForgeConfig{
		InstalledComponents: []InstalledComponent{
			{Name: "error-handler", Version: "1.0.0"},
			{Name: "jwt-auth", Version: "1.0.0"},
			{Name: "cors", Version: "1.0.0"},
		},
	}

	// Remove middle element
	removed := cfg.RemoveComponent("jwt-auth")
	if !removed {
		t.Error("RemoveComponent('jwt-auth') should return true")
	}
	if len(cfg.InstalledComponents) != 2 {
		t.Errorf("InstalledComponents count = %d, want 2", len(cfg.InstalledComponents))
	}
	if cfg.HasComponent("jwt-auth") {
		t.Error("jwt-auth should no longer be in InstalledComponents")
	}

	// Remove non-existent
	removed = cfg.RemoveComponent("nonexistent")
	if removed {
		t.Error("RemoveComponent('nonexistent') should return false")
	}
}

// ──────────────────────────────────────────────────────────────
// Tests: IsV2
// ──────────────────────────────────────────────────────────────

func TestIsV2(t *testing.T) {
	tests := []struct {
		version string
		want    bool
	}{
		{"2.0", true},
		{"1.0", false},
		{"", false},
	}

	for _, tt := range tests {
		cfg := &ForgeConfig{ConfigVersion: tt.version}
		got := cfg.IsV2()
		if got != tt.want {
			t.Errorf("IsV2() with version %q = %v, want %v", tt.version, got, tt.want)
		}
	}
}

// ──────────────────────────────────────────────────────────────
// Tests: GetStackKey
// ──────────────────────────────────────────────────────────────

func TestGetStackKey(t *testing.T) {
	tests := []struct {
		language  string
		framework string
		want      string
	}{
		{"nodejs", "express", "nodejs_express"},
		{"csharp", "dotnet-webapi", "csharp_dotnet-webapi"},
		{"golang", "gin", "golang_gin"},
	}

	for _, tt := range tests {
		cfg := &ForgeConfig{Language: tt.language, Framework: tt.framework}
		got := cfg.GetStackKey()
		if got != tt.want {
			t.Errorf("GetStackKey() = %q, want %q", got, tt.want)
		}
	}
}

// ──────────────────────────────────────────────────────────────
// Tests: GetORMOptions
// ──────────────────────────────────────────────────────────────

func TestGetORMOptions(t *testing.T) {
	tests := []struct {
		language string
		minLen   int
	}{
		{"nodejs", 3},  // prisma, mongoose, sequelize, none
		{"csharp", 2},  // efcore, dapper, none
		{"golang", 2},  // gorm, sqlx, none
		{"unknown", 1}, // fallback to "none"
	}

	for _, tt := range tests {
		cfg := &ForgeConfig{Language: tt.language}
		options := cfg.GetORMOptions()
		if len(options) < tt.minLen {
			t.Errorf("GetORMOptions(%q) returned %d options, want at least %d", tt.language, len(options), tt.minLen)
		}
	}
}

// ──────────────────────────────────────────────────────────────
// Tests: Load Error Cases
// ──────────────────────────────────────────────────────────────

func TestLoad_FileNotFound(t *testing.T) {
	cleanup := setupTestDir(t)
	defer cleanup()

	_, err := Load()
	if err == nil {
		t.Error("Load() should return error when file doesn't exist")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	cleanup := setupTestDir(t)
	defer cleanup()

	writeRawJSON(t, `{invalid json}`)

	_, err := Load()
	if err == nil {
		t.Error("Load() should return error for invalid JSON")
	}
}
