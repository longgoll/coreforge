package config

import (
	"encoding/json"
	"os"
	"time"
)

const configFileName = ".forge.json"

// ──────────────────────────────────────────────────────────────
// Data Structures
// ──────────────────────────────────────────────────────────────

// ForgeConfig represents the .forge.json project configuration file.
// Supports both v1 (legacy) and v2 (new fields).
// V2 fields use omitempty for backward compatibility — v1 configs still load fine.
type ForgeConfig struct {
	// V1 fields (always present)
	Language            string               `json:"language"`
	Framework           string               `json:"framework"`
	SourceDir           string               `json:"sourceDir"`
	InstalledComponents []InstalledComponent `json:"installedComponents"`

	// V2 fields (optional, backward compatible)
	ConfigVersion string `json:"configVersion,omitempty"` // "2.0" for v2 configs
	Database      string `json:"database,omitempty"`      // "postgresql", "mongodb", "mysql", "sqlite", "none"
	ORM           string `json:"orm,omitempty"`            // "prisma", "mongoose", "efcore", "gorm", "none"
	Architecture  string `json:"architecture,omitempty"`   // "mvc", "feature-based", "clean", "minimal"
}

// InstalledComponent tracks a component that has been added to the project.
type InstalledComponent struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	InstalledAt string            `json:"installedAt"`
	FileHashes  map[string]string `json:"fileHashes,omitempty"` // target path → SHA256 hash
}

// ──────────────────────────────────────────────────────────────
// Constants — Valid options for V2 fields
// ──────────────────────────────────────────────────────────────

// Valid database options
var ValidDatabases = []string{"postgresql", "mongodb", "mysql", "sqlite", "none"}

// Valid ORM options per language
var ValidORMs = map[string][]string{
	"nodejs": {"prisma", "mongoose", "sequelize", "none"},
	"csharp": {"efcore", "dapper", "none"},
	"golang": {"gorm", "sqlx", "none"},
}

// Valid architecture options
var ValidArchitectures = []string{"mvc", "feature-based", "clean", "minimal"}

// ──────────────────────────────────────────────────────────────
// Public API
// ──────────────────────────────────────────────────────────────

// Exists checks if a .forge.json file exists in the current directory.
func Exists() bool {
	_, err := os.Stat(configFileName)
	return err == nil
}

// Load reads and parses the .forge.json file from the current directory.
// Supports both v1 and v2 format. V1 configs are auto-migrated to v2 in memory.
func Load() (*ForgeConfig, error) {
	data, err := os.ReadFile(configFileName)
	if err != nil {
		return nil, err
	}

	var cfg ForgeConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Auto-migrate v1 → v2: fill defaults for missing v2 fields
	if cfg.ConfigVersion == "" {
		cfg.ConfigVersion = "1.0"
		// Don't set defaults for DB/ORM/Architecture — they were not configured
	}

	return &cfg, nil
}

// Save writes the ForgeConfig to .forge.json in the current directory.
func Save(cfg *ForgeConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configFileName, data, 0644)
}

// ──────────────────────────────────────────────────────────────
// Helper Methods
// ──────────────────────────────────────────────────────────────

// HasComponent checks if a component is already installed.
func (c *ForgeConfig) HasComponent(name string) bool {
	for _, comp := range c.InstalledComponents {
		if comp.Name == name {
			return true
		}
	}
	return false
}

// AddComponent adds a component to the installed list.
func (c *ForgeConfig) AddComponent(name string) {
	c.InstalledComponents = append(c.InstalledComponents, InstalledComponent{
		Name:        name,
		Version:     "1.0.0",
		InstalledAt: time.Now().Format(time.RFC3339),
	})
}

// AddComponentWithHashes adds a component with file hash tracking.
func (c *ForgeConfig) AddComponentWithHashes(name string, hashes map[string]string) {
	c.InstalledComponents = append(c.InstalledComponents, InstalledComponent{
		Name:        name,
		Version:     "1.0.0",
		InstalledAt: time.Now().Format(time.RFC3339),
		FileHashes:  hashes,
	})
}

// GetInstalledComponent returns the InstalledComponent by name, or nil if not found.
func (c *ForgeConfig) GetInstalledComponent(name string) *InstalledComponent {
	for i, comp := range c.InstalledComponents {
		if comp.Name == name {
			return &c.InstalledComponents[i]
		}
	}
	return nil
}

// RemoveComponent removes a component from the installed list.
// Returns true if the component was found and removed, false otherwise.
func (c *ForgeConfig) RemoveComponent(name string) bool {
	for i, comp := range c.InstalledComponents {
		if comp.Name == name {
			c.InstalledComponents = append(c.InstalledComponents[:i], c.InstalledComponents[i+1:]...)
			return true
		}
	}
	return false
}

// IsV2 returns true if this config was created with v2 fields.
func (c *ForgeConfig) IsV2() bool {
	return c.ConfigVersion == "2.0"
}

// GetStackKey returns the combined language_framework key used for registry lookups.
func (c *ForgeConfig) GetStackKey() string {
	return c.Language + "_" + c.Framework
}

// GetORMOptions returns valid ORM options for the current language.
func (c *ForgeConfig) GetORMOptions() []string {
	if orms, ok := ValidORMs[c.Language]; ok {
		return orms
	}
	return []string{"none"}
}
