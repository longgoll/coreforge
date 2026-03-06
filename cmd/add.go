package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/longgoll/forge-cli/internal/config"
	"github.com/longgoll/forge-cli/internal/env"
	"github.com/longgoll/forge-cli/internal/registry"
	"github.com/longgoll/forge-cli/internal/tui"
	"github.com/spf13/cobra"
)

var (
	addForce bool // --force flag to skip conflict prompts
)

var addCmd = &cobra.Command{
	Use:   "add <component>",
	Short: "Add a component, schema, or blueprint to your project",
	Long: fmt.Sprintf(`%s

Download and install a production-ready component, schema, or blueprint
into your project. The code is copied directly into your source tree — you own it.

The CLI reads your %s to determine which stack you're using,
then fetches the correct implementation from the registry.

Supports all registry tiers:
  • Components  — Single-purpose middleware/utilities
  • Schemas     — Database models (forge add schema-user-auth)
  • Blueprints  — Complete workflows (forge add blueprint-auth)

%s
  forge add error-handler
  forge add jwt-auth
  forge add schema-user-auth
  forge add blueprint-auth
  forge add jwt-auth --force    # Skip conflict prompts`, bold("📦 Add a registry item"), cyan(".forge.json"), dimmed("Examples:")),

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		itemName := args[0]

		fmt.Println(cyan("⚡"), bold("Forge CLI — Add"))
		fmt.Println(dimmed("─────────────────────────────────────────"))
		fmt.Println()

		// Step 1: Load .forge.json
		if !config.Exists() {
			fmt.Println(red("✗"), "No", cyan(".forge.json"), "found in the current directory.")
			fmt.Println()
			fmt.Println("  Run", cyan("forge init"), "first to set up your project.")
			return fmt.Errorf("config not found")
		}

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		stackKey := cfg.GetStackKey()
		fmt.Printf("  %s %s + %s\n", dimmed("Stack:"), green(cfg.Language), green(cfg.Framework))
		fmt.Printf("  %s %s\n", dimmed("Item:"), bold(itemName))
		fmt.Println()

		// Step 2: Check if already installed
		if cfg.HasComponent(itemName) {
			fmt.Println(yellow("⚠"), "Item", cyan(itemName), "is already installed.")
			fmt.Println(dimmed("  If you want to reinstall, remove it from .forge.json first."))
			return nil
		}

		// Step 3: Load manifest
		manifest, err := registry.LoadManifest()
		if err != nil {
			return fmt.Errorf("failed to load registry: %w", err)
		}

		// Step 4: Resolve item across all tiers
		tier, found := manifest.ResolveItem(itemName)
		if !found {
			fmt.Println(red("✗"), "Item", cyan(itemName), "not found in registry.")
			fmt.Println()
			fmt.Println(dimmed("Available items:"))
			for name, comp := range manifest.Components {
				fmt.Printf("  %s %s — %s %s\n", dimmed("•"), cyan(name), comp.Description, dimmed("[component]"))
			}
			for name, schema := range manifest.Schemas {
				fmt.Printf("  %s %s — %s %s\n", dimmed("•"), cyan(name), schema.Description, dimmed("[schema]"))
			}
			for name, bp := range manifest.Blueprints {
				fmt.Printf("  %s %s — %s %s\n", dimmed("•"), cyan(name), bp.Description, dimmed("[blueprint]"))
			}
			return fmt.Errorf("item not found: %s", itemName)
		}

		// Dispatch based on tier
		switch tier {
		case "component":
			return addComponent(itemName, stackKey, cfg, manifest)
		case "schema":
			return addSchema(itemName, stackKey, cfg, manifest)
		case "blueprint":
			return addBlueprint(itemName, stackKey, cfg, manifest)
		default:
			return fmt.Errorf("unknown tier: %s", tier)
		}
	},
}

// addComponent handles installing a single component (Tier 2).
func addComponent(name, stackKey string, cfg *config.ForgeConfig, manifest *registry.Manifest) error {
	component := manifest.Components[name]
	impl, exists := component.Implementations[stackKey]
	if !exists {
		return printStackNotAvailable(name, "component", stackKey, component.Implementations)
	}

	return installImplementation(name, "component", stackKey, cfg, impl)
}

// addSchema handles installing a schema (Tier 3).
func addSchema(name, stackKey string, cfg *config.ForgeConfig, manifest *registry.Manifest) error {
	schema := manifest.Schemas[name]
	impl, exists := schema.Implementations[stackKey]
	if !exists {
		return printStackNotAvailable(name, "schema", stackKey, schema.Implementations)
	}

	return installImplementation(name, "schema", stackKey, cfg, impl)
}

// addBlueprint handles installing a blueprint (Tier 4).
// Blueprints install their included components/schemas first, then their own files.
func addBlueprint(name, stackKey string, cfg *config.ForgeConfig, manifest *registry.Manifest) error {
	blueprint := manifest.Blueprints[name]

	fmt.Println(bold("  Blueprint:"), cyan(name), "—", blueprint.Description)
	fmt.Println()
	fmt.Println(dimmed("  This blueprint will install:"))
	for _, inc := range blueprint.Includes {
		status := green("  ○")
		if cfg.HasComponent(inc) {
			status = dimmed("  ✓") // Already installed
		}
		fmt.Printf("  %s %s\n", status, cyan(inc))
	}
	fmt.Println()

	// Install included components/schemas (skip if already installed)
	for _, incName := range blueprint.Includes {
		if cfg.HasComponent(incName) {
			fmt.Printf("  %s %s %s\n", dimmed("⊘"), cyan(incName), dimmed("(already installed, skipping)"))
			continue
		}

		incTier, incFound := manifest.ResolveItem(incName)
		if !incFound {
			fmt.Println(yellow("⚠"), "Included item", cyan(incName), "not found in registry, skipping")
			continue
		}

		switch incTier {
		case "component":
			comp := manifest.Components[incName]
			impl, exists := comp.Implementations[stackKey]
			if !exists {
				fmt.Println(yellow("⚠"), cyan(incName), "not available for", yellow(stackKey), ", skipping")
				continue
			}
			hashes, err := installFilesWithHashes(incName, cfg, impl)
			if err != nil {
				fmt.Println(yellow("⚠"), "Failed to install", cyan(incName), ":", err)
				continue
			}
			cfg.AddComponentWithHashes(incName, hashes)
			fmt.Printf("  %s %s installed\n", green("✓"), cyan(incName))

		case "schema":
			schema := manifest.Schemas[incName]
			impl, exists := schema.Implementations[stackKey]
			if !exists {
				fmt.Println(yellow("⚠"), cyan(incName), "not available for", yellow(stackKey), ", skipping")
				continue
			}
			hashes, err := installFilesWithHashes(incName, cfg, impl)
			if err != nil {
				fmt.Println(yellow("⚠"), "Failed to install", cyan(incName), ":", err)
				continue
			}
			cfg.AddComponentWithHashes(incName, hashes)
			fmt.Printf("  %s %s installed\n", green("✓"), cyan(incName))
		}
	}

	// Install blueprint-specific files (controllers, routes, etc.)
	impl, exists := blueprint.Implementations[stackKey]
	if exists && len(impl.Files) > 0 {
		fmt.Println()
		fmt.Println(dimmed("  Installing blueprint files..."))
		if _, err := installFilesWithHashes(name, cfg, impl); err != nil {
			return fmt.Errorf("failed to install blueprint files: %w", err)
		}
	}

	// Auto-configure .env (aggregate all envVars from includes + blueprint)
	allEnvVars := collectBlueprintEnvVars(blueprint, manifest, stackKey)
	if len(allEnvVars) > 0 {
		configureEnvVars(cfg, allEnvVars)
	}

	// Save config
	cfg.AddComponent(name)
	if err := config.Save(cfg); err != nil {
		fmt.Println(yellow("⚠"), "Failed to update .forge.json:", err)
	}

	// Post-install for blueprint
	fmt.Println()
	fmt.Println(green("✓"), bold("Blueprint"), cyan(name), bold("installed successfully!"))

	if exists && impl.PostInstall != "" {
		fmt.Println()
		fmt.Println(yellow("📋"), bold("Post-install steps:"))
		fmt.Println(dimmed("  " + impl.PostInstall))
	}

	fmt.Println()
	return nil
}

// collectBlueprintEnvVars gathers all env vars from blueprint includes + blueprint itself.
func collectBlueprintEnvVars(bp registry.Blueprint, manifest *registry.Manifest, stackKey string) []registry.EnvVar {
	var allVars []registry.EnvVar
	seen := make(map[string]bool)

	// Collect from includes
	for _, incName := range bp.Includes {
		tier, found := manifest.ResolveItem(incName)
		if !found {
			continue
		}

		var impl registry.Implementation
		var exists bool
		switch tier {
		case "component":
			comp := manifest.Components[incName]
			impl, exists = comp.Implementations[stackKey]
		case "schema":
			schema := manifest.Schemas[incName]
			impl, exists = schema.Implementations[stackKey]
		}

		if exists {
			for _, ev := range impl.EnvVars {
				if !seen[ev.Key] {
					allVars = append(allVars, ev)
					seen[ev.Key] = true
				}
			}
		}
	}

	// Collect from blueprint itself
	if impl, exists := bp.Implementations[stackKey]; exists {
		for _, ev := range impl.EnvVars {
			if !seen[ev.Key] {
				allVars = append(allVars, ev)
				seen[ev.Key] = true
			}
		}
	}

	return allVars
}

// ──────────────────────────────────────────────────────────────
// Shared Installation Logic
// ──────────────────────────────────────────────────────────────

// installImplementation handles the full flow of installing an implementation:
// download files, install dependencies, update config, show post-install.
func installImplementation(name, tier, stackKey string, cfg *config.ForgeConfig, impl registry.Implementation) error {
	// Download and write files (with conflict resolution + hash tracking)
	hashes, err := installFilesWithHashes(name, cfg, impl)
	if err != nil {
		return err
	}

	// Install dependencies
	if impl.InstallCmd != "" {
		fmt.Println()
		fmt.Println(dimmed("  Installing dependencies..."))
		fmt.Printf("  %s %s\n", dimmed("$"), impl.InstallCmd)

		if err := runShellCommand(impl.InstallCmd, cfg.SourceDir); err != nil {
			fmt.Println(yellow("⚠"), "Dependency install failed:", err)
			fmt.Println(dimmed("  You may need to run manually:"), cyan(impl.InstallCmd))
		} else {
			fmt.Println(green("  ✓"), "Dependencies installed")
		}
	}

	// Install dev dependencies (V2)
	if impl.InstallDevCmd != "" {
		fmt.Printf("  %s %s\n", dimmed("$"), impl.InstallDevCmd)
		if err := runShellCommand(impl.InstallDevCmd, cfg.SourceDir); err != nil {
			fmt.Println(yellow("⚠"), "Dev dependency install failed:", err)
			fmt.Println(dimmed("  You may need to run manually:"), cyan(impl.InstallDevCmd))
		} else {
			fmt.Println(green("  ✓"), "Dev dependencies installed")
		}
	}

	// Auto-configure .env
	if len(impl.EnvVars) > 0 {
		configureEnvVars(cfg, impl.EnvVars)
	}

	// Update .forge.json (with hashes)
	cfg.AddComponentWithHashes(name, hashes)
	if err := config.Save(cfg); err != nil {
		fmt.Println(yellow("⚠"), "Failed to update .forge.json:", err)
	}

	// Success message
	fmt.Println()
	fmt.Printf("%s %s %s %s %s\n", green("✓"), bold(capitalize(tier)), cyan(name), bold("added successfully!"), dimmed("["+tier+"]"))

	if impl.PostInstall != "" {
		fmt.Println()
		fmt.Println(yellow("📋"), bold("Post-install steps:"))
		fmt.Println(dimmed("  " + impl.PostInstall))
	}

	// Show required components if any
	if len(impl.Requires) > 0 {
		fmt.Println()
		fmt.Println(yellow("💡"), "This", tier, "works best with:")
		for _, req := range impl.Requires {
			if !cfg.HasComponent(req) {
				fmt.Printf("  %s Run: %s\n", dimmed("•"), cyan("forge add "+req))
			}
		}
	}

	fmt.Println()
	return nil
}

// installFilesWithHashes downloads, writes files (with conflict resolution),
// and returns a map of target path → SHA256 hash for tracking.
func installFilesWithHashes(name string, cfg *config.ForgeConfig, impl registry.Implementation) (map[string]string, error) {
	hashes := make(map[string]string)

	fmt.Println(dimmed("  Downloading files..."))
	for _, file := range impl.Files {
		targetPath := filepath.Join(cfg.SourceDir, file.Target)

		// Create parent directories
		dir := filepath.Dir(targetPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return hashes, fmt.Errorf("failed to create directory %s: %w", dir, err)
		}

		// Download file content
		content, err := downloadFile(file.URL)
		if err != nil {
			return hashes, fmt.Errorf("failed to download %s: %w", file.URL, err)
		}

		// Process content through Go text/template engine
		content = processTemplate(content)

		// Compute hash of new content
		newHash := config.HashBytes(content)

		// ── Conflict Resolution ──────────────────────────
		if _, err := os.Stat(targetPath); err == nil {
			// File already exists — check if user has modified it
			currentHash, _ := config.HashFile(targetPath)

			// Find the original hash (stored when component was first installed)
			originalHash := ""
			if installed := cfg.GetInstalledComponent(name); installed != nil {
				if h, ok := installed.FileHashes[file.Target]; ok {
					originalHash = h
				}
			}

			// If the file hasn't changed from original install, safe to overwrite
			if originalHash != "" && currentHash == originalHash {
				// File unmodified by user → overwrite silently
			} else if currentHash == newHash {
				// Same content → skip
				fmt.Printf("  %s %s %s\n", dimmed("·"), targetPath, dimmed("(unchanged)"))
				hashes[file.Target] = newHash
				continue
			} else if !addForce {
				// File was modified by user → ask what to do
				fmt.Printf("  %s %s %s\n", yellow("⚠"), targetPath, yellow("(modified by user)"))

				action, err := tui.SelectConflictAction(targetPath)
				if err != nil {
					return hashes, fmt.Errorf("conflict prompt failed: %w", err)
				}

				switch action {
				case "skip":
					fmt.Printf("  %s %s %s\n", dimmed("⊘"), targetPath, dimmed("(skipped)"))
					// Keep existing hash
					hashes[file.Target] = currentHash
					continue
				case "backup":
					backupPath := targetPath + ".bak"
					if err := copyFile(targetPath, backupPath); err != nil {
						fmt.Println(yellow("⚠"), "Failed to backup:", err)
					} else {
						fmt.Printf("  %s %s %s\n", green("✓"), dimmed("Backup:"), backupPath)
					}
					// Fall through to overwrite
				case "overwrite":
					// Fall through to overwrite
				}
			}
		}

		// Write file
		if err := os.WriteFile(targetPath, content, 0644); err != nil {
			return hashes, fmt.Errorf("failed to write %s: %w", targetPath, err)
		}

		hashes[file.Target] = newHash
		fmt.Printf("  %s %s\n", green("✓"), targetPath)
	}
	return hashes, nil
}

// configureEnvVars handles auto .env configuration for a component.
func configureEnvVars(cfg *config.ForgeConfig, envVars []registry.EnvVar) {
	if len(envVars) == 0 {
		return
	}

	// Determine project root (parent of sourceDir or current dir)
	projectDir := "."
	if cfg.SourceDir != "" && cfg.SourceDir != "./" && cfg.SourceDir != "." {
		// Go up from sourceDir to project root
		projectDir = filepath.Dir(filepath.Clean(cfg.SourceDir))
		if projectDir == "." || projectDir == cfg.SourceDir {
			projectDir = "."
		}
	}

	// Ensure .env files exist
	if err := env.EnsureEnvFiles(projectDir); err != nil {
		fmt.Println(yellow("⚠"), "Failed to create .env files:", err)
		return
	}

	// Convert registry.EnvVar to env.EnvVarEntry
	entries := make([]env.EnvVarEntry, len(envVars))
	for i, ev := range envVars {
		entries[i] = env.EnvVarEntry{
			Key:         ev.Key,
			Value:       ev.Default,
			Description: ev.Description,
		}
	}

	// Append env vars
	added, err := env.AppendEnvVars(projectDir, entries)
	if err != nil {
		fmt.Println(yellow("⚠"), "Failed to configure .env:", err)
		return
	}

	if added > 0 {
		fmt.Println()
		fmt.Printf("  %s %d environment variable(s) added to %s\n", green("✓"), added, cyan(".env"))
	}
}

// ──────────────────────────────────────────────────────────────
// File Helpers
// ──────────────────────────────────────────────────────────────

// printStackNotAvailable shows an error when a stack implementation isn't found.
func printStackNotAvailable(name, tier, stackKey string, implementations map[string]registry.Implementation) error {
	fmt.Println(red("✗"), capitalize(tier), cyan(name), "is not available for stack", yellow(stackKey))
	fmt.Println()
	fmt.Println(dimmed("Available implementations:"))
	for key := range implementations {
		fmt.Printf("  %s %s\n", dimmed("•"), green(key))
	}
	return fmt.Errorf("no implementation for stack: %s", stackKey)
}

// copyFile copies a file from src to dst.
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

// ──────────────────────────────────────────────────────────────
// File Download & Shell Helpers
// ──────────────────────────────────────────────────────────────

// getExeDir returns the directory where the forge binary is located.
func getExeDir() string {
	exe, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Dir(exe)
}

// downloadFile fetches a file from a URL and returns its content.
// Supports both HTTP URLs and local file paths (for development).
func downloadFile(url string) ([]byte, error) {
	// If it's a local file path (for development/testing)
	if strings.HasPrefix(url, "file://") {
		localPath := strings.TrimPrefix(url, "file://")
		return os.ReadFile(localPath)
	}

	// If it's a relative path (for local mock registry),
	// resolve relative to the binary location
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		// Strip leading "./" if present
		cleanPath := strings.TrimPrefix(url, "./")
		absPath := filepath.Join(getExeDir(), cleanPath)
		return os.ReadFile(absPath)
	}

	// HTTP download
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	return io.ReadAll(resp.Body)
}

// runShellCommand executes a shell command in the given directory.
func runShellCommand(command string, dir string) error {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// capitalize returns the string with the first letter uppercased.
func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func init() {
	addCmd.Flags().BoolVarP(&addForce, "force", "f", false, "Skip conflict resolution prompts (overwrite all)")
	rootCmd.AddCommand(addCmd)
}
