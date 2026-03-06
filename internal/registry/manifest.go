package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// ──────────────────────────────────────────────────────────────
// Registry URL Configuration
// ──────────────────────────────────────────────────────────────

const (
	// DEV: Local mock registry for development/testing
	LocalRegistryRelPath = "mock-registry/manifest.json"

	// PROD: GitHub raw content URL
	RemoteRegistryURL = "https://raw.githubusercontent.com/longgoll/forge-registry/main/manifest.json"

	// Cache settings
	CacheDirName      = ".forge"
	CacheFileName     = "manifest.json"
	CacheMetaFileName = "cache_meta.json"
	CacheTTL          = 1 * time.Hour // Re-validate cache after 1 hour
)

// UseRemote toggles between local and remote registry.
// Set to true when the GitHub registry is ready.
var UseRemote = true

// getExeDir returns the directory where the forge binary is located.
func getExeDir() string {
	exe, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Dir(exe)
}

// getCacheDir returns the path to ~/.forge/ cache directory.
func getCacheDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	cacheDir := filepath.Join(home, CacheDirName, "cache")
	return cacheDir, nil
}

// ──────────────────────────────────────────────────────────────
// Cache Metadata
// ──────────────────────────────────────────────────────────────

// CacheMeta stores metadata about the cached manifest.
type CacheMeta struct {
	ETag         string    `json:"etag"`
	LastModified string    `json:"lastModified"`
	CachedAt     time.Time `json:"cachedAt"`
	RegistryURL  string    `json:"registryUrl"`
}

func loadCacheMeta(cacheDir string) (*CacheMeta, error) {
	metaPath := filepath.Join(cacheDir, CacheMetaFileName)
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, err
	}
	var meta CacheMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}

func saveCacheMeta(cacheDir string, meta *CacheMeta) error {
	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(cacheDir, CacheMetaFileName), data, 0644)
}

// ──────────────────────────────────────────────────────────────
// Manifest Data Structures — V2 (backward compatible with V1)
// ──────────────────────────────────────────────────────────────

// Manifest represents the top-level registry manifest.json structure.
// V2 adds SchemaVersion, Foundations, Schemas, and Blueprints fields.
// V1 manifests (without schemaVersion) are still parsed correctly.
type Manifest struct {
	// V1 fields (always present)
	Version    string               `json:"version"`
	Languages  map[string]Language  `json:"languages"`
	Components map[string]Component `json:"components"`

	// V2 fields (optional, backward compatible)
	SchemaVersion string                `json:"schemaVersion,omitempty"` // "2.0.0" for v2
	Foundations   map[string]Foundation `json:"foundations,omitempty"`   // Tầng 1: Project templates
	Schemas       map[string]Schema     `json:"schemas,omitempty"`       // Tầng 3: DB models
	Blueprints    map[string]Blueprint  `json:"blueprints,omitempty"`    // Tầng 4: Composite workflows
}

// Language describes a supported programming language.
type Language struct {
	Name       string               `json:"name"`
	Frameworks map[string]Framework `json:"frameworks"`
}

// Framework describes a supported framework within a language.
type Framework struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	TemplateURL string `json:"templateUrl,omitempty"`
	MinVersion  string `json:"minVersion,omitempty"`
}

// Component describes a reusable component that can be added to a project.
// This is Tầng 2 in the registry architecture.
type Component struct {
	Description     string                    `json:"description"`
	Category        string                    `json:"category"`
	Tags            []string                  `json:"tags"`
	Implementations map[string]Implementation `json:"implementations"`
}

// Implementation describes how a component is implemented for a specific stack.
type Implementation struct {
	Files           []FileEntry `json:"files"`
	Dependencies    []string    `json:"dependencies"`
	DevDependencies []string    `json:"devDependencies,omitempty"` // V2: dev-only deps
	InstallCmd      string      `json:"installCmd,omitempty"`
	InstallDevCmd   string      `json:"installDevCmd,omitempty"` // V2: dev deps install cmd
	PostInstall     string      `json:"postInstall,omitempty"`
	Requires        []string    `json:"requires,omitempty"`
	Conflicts       []string    `json:"conflicts,omitempty"`
	EnvVars         []EnvVar    `json:"envVars,omitempty"` // V2: environment variables
}

// FileEntry describes a single file to download and place in the project.
type FileEntry struct {
	URL    string `json:"url"`
	Target string `json:"target"`
}

// EnvVar describes an environment variable needed by a component.
type EnvVar struct {
	Key         string `json:"key"`
	Default     string `json:"default"`
	Description string `json:"description"`
}

// ──────────────────────────────────────────────────────────────
// V2 New Tier Structures
// ──────────────────────────────────────────────────────────────

// Foundation represents a Tầng 1 project skeleton/template.
// Used with `forge init` to set up complete project structure.
type Foundation struct {
	Description  string   `json:"description"`
	Language     string   `json:"language"`
	Framework    string   `json:"framework"`
	Architecture string   `json:"architecture,omitempty"` // "mvc", "feature-based", "clean", "minimal"
	TemplateURL  string   `json:"templateUrl"`
	Includes     []string `json:"includes,omitempty"` // What's included: "tsconfig", "eslint", "dockerfile", etc.
}

// Schema represents a Tầng 3 database model/entity.
// Used with `forge add schema-<name>`.
type Schema struct {
	Description     string                    `json:"description"`
	Category        string                    `json:"category"`
	Tags            []string                  `json:"tags"`
	Implementations map[string]Implementation `json:"implementations"`
}

// Blueprint represents a Tầng 4 composite workflow.
// Bundles multiple components + schemas + custom files.
// Used with `forge add blueprint-<name>`.
type Blueprint struct {
	Description     string                    `json:"description"`
	Category        string                    `json:"category,omitempty"`
	Tags            []string                  `json:"tags,omitempty"`
	Includes        []string                  `json:"includes"`        // Components + schemas to install
	Implementations map[string]Implementation `json:"implementations"` // Blueprint-specific files
}

// ──────────────────────────────────────────────────────────────
// Loading Functions
// ──────────────────────────────────────────────────────────────

// LoadManifest loads the manifest from the configured source (local or remote).
func LoadManifest() (*Manifest, error) {
	if UseRemote {
		return loadRemoteManifest()
	}
	return loadLocalManifest()
}

// loadLocalManifest reads the manifest from a local JSON file.
// Path is resolved relative to the forge binary location.
func loadLocalManifest() (*Manifest, error) {
	manifestPath := filepath.Join(getExeDir(), LocalRegistryRelPath)
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read local manifest at %s: %w\n"+
			"Make sure the mock-registry directory exists with a manifest.json file", manifestPath, err)
	}

	return parseManifest(data)
}

// loadRemoteManifest fetches the manifest from the remote registry URL.
// Uses a local cache with ETag-based revalidation to minimize network requests.
//
// Flow:
//  1. Check if cached manifest exists and is fresh (< CacheTTL old)
//  2. If fresh → return cached version
//  3. If stale/missing → fetch from remote with conditional headers
//  4. If 304 Not Modified → return cached version
//  5. If 200 OK → save to cache and return
//  6. If fetch fails but cache exists → return stale cache (offline fallback)
func loadRemoteManifest() (*Manifest, error) {
	cacheDir, err := getCacheDir()
	if err != nil {
		// Can't cache, fetch directly
		return fetchManifestFromURL(RemoteRegistryURL)
	}

	cachePath := filepath.Join(cacheDir, CacheFileName)

	// Try to load cache metadata
	meta, metaErr := loadCacheMeta(cacheDir)

	// If cache is fresh (within TTL), return it directly
	if metaErr == nil && time.Since(meta.CachedAt) < CacheTTL {
		manifest, err := loadManifestFromFile(cachePath)
		if err == nil {
			return manifest, nil
		}
		// Cache file corrupted, fall through to re-fetch
	}

	// Fetch from remote with conditional headers
	manifest, err := fetchAndCacheManifest(cacheDir, cachePath, meta)
	if err != nil {
		// If fetch failed but we have a cached version, use it (offline fallback)
		if metaErr == nil {
			manifest, cacheErr := loadManifestFromFile(cachePath)
			if cacheErr == nil {
				return manifest, nil
			}
		}
		return nil, fmt.Errorf("failed to fetch remote registry: %w", err)
	}

	return manifest, nil
}

// fetchAndCacheManifest fetches the manifest from remote, using conditional
// headers if we have cache metadata, and saves the result to cache.
func fetchAndCacheManifest(cacheDir, cachePath string, meta *CacheMeta) (*Manifest, error) {
	req, err := http.NewRequest("GET", RemoteRegistryURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "forge-cli/1.0")
	req.Header.Set("Accept", "application/json")

	// Add conditional headers if we have cache metadata
	if meta != nil {
		if meta.ETag != "" {
			req.Header.Set("If-None-Match", meta.ETag)
		}
		if meta.LastModified != "" {
			req.Header.Set("If-Modified-Since", meta.LastModified)
		}
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network request failed: %w", err)
	}
	defer resp.Body.Close()

	// 304 Not Modified — cache is still valid
	if resp.StatusCode == http.StatusNotModified && meta != nil {
		// Update cache timestamp
		meta.CachedAt = time.Now()
		_ = saveCacheMeta(cacheDir, meta)

		return loadManifestFromFile(cachePath)
	}

	// Non-200 response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("registry returned HTTP %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse manifest
	manifest, err := parseManifest(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse remote manifest JSON: %w", err)
	}

	// Save to cache
	if err := os.MkdirAll(cacheDir, 0755); err == nil {
		_ = os.WriteFile(cachePath, body, 0644)

		newMeta := &CacheMeta{
			ETag:         resp.Header.Get("ETag"),
			LastModified: resp.Header.Get("Last-Modified"),
			CachedAt:     time.Now(),
			RegistryURL:  RemoteRegistryURL,
		}
		_ = saveCacheMeta(cacheDir, newMeta)
	}

	return manifest, nil
}

// fetchManifestFromURL is a simple fallback that fetches without caching.
func fetchManifestFromURL(url string) (*Manifest, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return parseManifest(body)
}

// loadManifestFromFile reads and parses a manifest from a local file path.
func loadManifestFromFile(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return parseManifest(data)
}

// parseManifest parses JSON data into a Manifest struct.
// Handles both v1 and v2 schema formats.
// V1 manifests (without schemaVersion) are normalized with empty v2 collections.
func parseManifest(data []byte) (*Manifest, error) {
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest JSON: %w", err)
	}

	// Normalize: ensure maps are never nil (safe for iteration)
	if manifest.Components == nil {
		manifest.Components = make(map[string]Component)
	}
	if manifest.Foundations == nil {
		manifest.Foundations = make(map[string]Foundation)
	}
	if manifest.Schemas == nil {
		manifest.Schemas = make(map[string]Schema)
	}
	if manifest.Blueprints == nil {
		manifest.Blueprints = make(map[string]Blueprint)
	}

	return &manifest, nil
}

// ClearCache removes the cached manifest files.
func ClearCache() error {
	cacheDir, err := getCacheDir()
	if err != nil {
		return err
	}

	cachePath := filepath.Join(cacheDir, CacheFileName)
	metaPath := filepath.Join(cacheDir, CacheMetaFileName)

	os.Remove(cachePath)
	os.Remove(metaPath)

	return nil
}

// GetCacheInfo returns info about the current cache state.
// Returns nil if no cache exists.
func GetCacheInfo() *CacheMeta {
	cacheDir, err := getCacheDir()
	if err != nil {
		return nil
	}
	meta, err := loadCacheMeta(cacheDir)
	if err != nil {
		return nil
	}
	return meta
}

// ──────────────────────────────────────────────────────────────
// Helper Functions — Language / Framework
// ──────────────────────────────────────────────────────────────

// GetLanguageNames returns a list of language display names for the TUI.
func (m *Manifest) GetLanguageNames() []string {
	names := make([]string, 0, len(m.Languages))
	for _, lang := range m.Languages {
		names = append(names, lang.Name)
	}
	return names
}

// GetLanguageKeys returns a list of language keys.
func (m *Manifest) GetLanguageKeys() []string {
	keys := make([]string, 0, len(m.Languages))
	for key := range m.Languages {
		keys = append(keys, key)
	}
	return keys
}

// GetFrameworkKeys returns the framework keys for a given language.
func (m *Manifest) GetFrameworkKeys(languageKey string) []string {
	lang, exists := m.Languages[languageKey]
	if !exists {
		return nil
	}
	keys := make([]string, 0, len(lang.Frameworks))
	for key := range lang.Frameworks {
		keys = append(keys, key)
	}
	return keys
}

// GetFrameworkNames returns the framework display names for a given language.
func (m *Manifest) GetFrameworkNames(languageKey string) []string {
	lang, exists := m.Languages[languageKey]
	if !exists {
		return nil
	}
	names := make([]string, 0, len(lang.Frameworks))
	for _, fw := range lang.Frameworks {
		names = append(names, fw.Name)
	}
	return names
}

// ──────────────────────────────────────────────────────────────
// Helper Functions — Schema Version Detection
// ──────────────────────────────────────────────────────────────

// IsV2 returns true if the manifest uses schema version 2.0.
func (m *Manifest) IsV2() bool {
	return m.SchemaVersion != ""
}

// HasFoundations returns true if foundations are defined in the manifest.
func (m *Manifest) HasFoundations() bool {
	return len(m.Foundations) > 0
}

// HasSchemas returns true if schemas are defined in the manifest.
func (m *Manifest) HasSchemas() bool {
	return len(m.Schemas) > 0
}

// HasBlueprints returns true if blueprints are defined in the manifest.
func (m *Manifest) HasBlueprints() bool {
	return len(m.Blueprints) > 0
}

// ──────────────────────────────────────────────────────────────
// Helper Functions — Cross-Tier Resolution
// ──────────────────────────────────────────────────────────────

// ResolveItem looks up an item name across all tiers (components, schemas, blueprints).
// Returns the tier name ("component", "schema", "blueprint") and whether it was found.
func (m *Manifest) ResolveItem(name string) (tier string, found bool) {
	if _, ok := m.Components[name]; ok {
		return "component", true
	}
	if _, ok := m.Schemas[name]; ok {
		return "schema", true
	}
	if _, ok := m.Blueprints[name]; ok {
		return "blueprint", true
	}
	return "", false
}

// GetAllItemNames returns all available item names across all tiers.
func (m *Manifest) GetAllItemNames() []string {
	names := make([]string, 0)
	for name := range m.Components {
		names = append(names, name)
	}
	for name := range m.Schemas {
		names = append(names, name)
	}
	for name := range m.Blueprints {
		names = append(names, name)
	}
	return names
}

// GetComponentCount returns the total number of components.
func (m *Manifest) GetComponentCount() int {
	return len(m.Components)
}

// GetSchemaCount returns the total number of schemas.
func (m *Manifest) GetSchemaCount() int {
	return len(m.Schemas)
}

// GetBlueprintCount returns the total number of blueprints.
func (m *Manifest) GetBlueprintCount() int {
	return len(m.Blueprints)
}

// GetFoundationCount returns the total number of foundations.
func (m *Manifest) GetFoundationCount() int {
	return len(m.Foundations)
}

// GetFoundationsForStack returns foundations matching a specific language+framework.
func (m *Manifest) GetFoundationsForStack(language, framework string) map[string]Foundation {
	result := make(map[string]Foundation)
	for name, f := range m.Foundations {
		if f.Language == language && f.Framework == framework {
			result[name] = f
		}
	}
	return result
}
