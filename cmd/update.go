package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/longgoll/forge-cli/internal/config"
	"github.com/longgoll/forge-cli/internal/registry"
	"github.com/spf13/cobra"
)

var (
	updateAll   bool
	updateCheck bool
)

var updateCmd = &cobra.Command{
	Use:   "update [component]",
	Short: "Update installed components to the latest version",
	Long: fmt.Sprintf(`%s

Check for and install updates to your installed components.
The CLI compares your installed versions against the latest registry.

%s
  forge update jwt-auth       # Update a specific component
  forge update --all          # Update all installed components
  forge update --check        # Check for updates without installing`, bold("🔄 Update Components"), dimmed("Examples:")),

	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(cyan("⚡"), bold("Forge CLI — Update"))
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
		fmt.Println()

		if len(cfg.InstalledComponents) == 0 {
			fmt.Println(dimmed("  No components installed yet."))
			fmt.Println(dimmed("  Use"), cyan("forge add <name>"), dimmed("to install a component"))
			fmt.Println()
			return nil
		}

		// Step 2: Load manifest (force remote if available)
		manifest, err := registry.LoadManifest()
		if err != nil {
			return fmt.Errorf("failed to load registry: %w", err)
		}

		// Step 3: Determine which components to update
		var componentsToUpdate []string

		if updateAll {
			for _, comp := range cfg.InstalledComponents {
				componentsToUpdate = append(componentsToUpdate, comp.Name)
			}
		} else if len(args) > 0 {
			componentName := args[0]
			if !cfg.HasComponent(componentName) {
				fmt.Println(red("✗"), "Component", cyan(componentName), "is not installed.")
				fmt.Println()
				fmt.Println(dimmed("  Installed components:"))
				for _, comp := range cfg.InstalledComponents {
					fmt.Printf("    %s %s\n", dimmed("•"), cyan(comp.Name))
				}
				return nil
			}
			componentsToUpdate = append(componentsToUpdate, componentName)
		} else {
			// No args and no --all: show usage
			fmt.Println(dimmed("  Specify a component to update, or use --all:"))
			fmt.Println()
			fmt.Printf("    %s\n", cyan("forge update <component>"))
			fmt.Printf("    %s\n", cyan("forge update --all"))
			fmt.Printf("    %s\n", cyan("forge update --check"))
			fmt.Println()
			return nil
		}

		// Step 4: Check each component
		updatesAvailable := 0
		updatesApplied := 0

		for _, compName := range componentsToUpdate {
			// Resolve across tiers
			tier, found := manifest.ResolveItem(compName)
			if !found {
				fmt.Printf("  %s %s — %s\n", yellow("⚠"), cyan(compName), dimmed("not found in registry (may have been removed)"))
				continue
			}

			// Get the implementation for current stack
			var impl registry.Implementation
			var exists bool

			switch tier {
			case "component":
				comp := manifest.Components[compName]
				impl, exists = comp.Implementations[stackKey]
			case "schema":
				schema := manifest.Schemas[compName]
				impl, exists = schema.Implementations[stackKey]
			case "blueprint":
				bp := manifest.Blueprints[compName]
				impl, exists = bp.Implementations[stackKey]
			}

			if !exists {
				fmt.Printf("  %s %s — %s\n", dimmed("·"), cyan(compName), dimmed("no implementation for "+stackKey))
				continue
			}

			// Check if files have changed (compare hashes)
			hasChanges := false
			installedComp := cfg.GetInstalledComponent(compName)

			if installedComp != nil && len(installedComp.FileHashes) > 0 {
				// Compare stored hashes with current file contents
				for _, file := range impl.Files {
					targetPath := filepath.Join(cfg.SourceDir, file.Target)
					currentHash, err := config.HashFile(targetPath)
					if err != nil {
						hasChanges = true // File missing = needs update
						break
					}
					if storedHash, ok := installedComp.FileHashes[file.Target]; ok {
						if currentHash == storedHash {
							// File hasn't been modified by user, safe to update
							hasChanges = true // We want to redownload
						}
					}
				}
			} else {
				// No hashes stored — treat as updateable
				hasChanges = true
			}

			if !hasChanges {
				fmt.Printf("  %s %s — %s\n", green("✓"), cyan(compName), dimmed("up to date"))
				continue
			}

			updatesAvailable++

			if updateCheck {
				fmt.Printf("  %s %s — %s %s\n", yellow("↑"), cyan(compName), yellow("update available"), dimmed("["+tier+"]"))
				continue
			}

			// Apply update
			fmt.Printf("  %s Updating %s...\n", cyan("↻"), cyan(compName))

			// Remove from installed list first (will be re-added by installImplementation)
			cfg.RemoveComponent(compName)

			if err := installImplementation(compName, tier, stackKey, cfg, impl); err != nil {
				fmt.Printf("  %s Failed to update %s: %s\n", red("✗"), cyan(compName), err)
				continue
			}

			updatesApplied++
		}

		// Summary
		fmt.Println()
		if updateCheck {
			if updatesAvailable == 0 {
				fmt.Println(green("✓"), bold("All components are up to date!"))
			} else {
				fmt.Printf("%s %s %d %s\n", yellow("↑"), bold("Updates available:"), updatesAvailable, dimmed("component(s)"))
				fmt.Println(dimmed("  Run"), cyan("forge update --all"), dimmed("to apply all updates"))
			}
		} else {
			if updatesApplied > 0 {
				fmt.Printf("%s %s %d %s\n", green("✓"), bold("Updated"), updatesApplied, bold("component(s) successfully!"))
			} else {
				fmt.Println(green("✓"), bold("All components are up to date!"))
			}
		}

		fmt.Println()
		return nil
	},
}

func init() {
	updateCmd.Flags().BoolVar(&updateAll, "all", false, "Update all installed components")
	updateCmd.Flags().BoolVar(&updateCheck, "check", false, "Check for updates without installing")
	rootCmd.AddCommand(updateCmd)
}
