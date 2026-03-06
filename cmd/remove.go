package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/longgoll/forge-cli/internal/config"
	"github.com/longgoll/forge-cli/internal/registry"
	"github.com/spf13/cobra"
)

var (
	removeForce bool
)

var removeCmd = &cobra.Command{
	Use:   "remove <component>",
	Short: "Remove an installed component from your project",
	Long: fmt.Sprintf(`%s

Remove an installed component by deleting its files and updating %s.

The CLI uses the registry manifest to determine which files belong to
the component, then removes them from your project.

%s
  forge remove error-handler
  forge remove jwt-auth
  forge remove logger --force`, bold("🗑️  Remove a component"), cyan(".forge.json"), dimmed("Examples:")),

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		componentName := args[0]

		fmt.Println(cyan("⚡"), bold("Forge CLI — Remove Component"))
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

		fmt.Printf("  %s %s + %s\n", dimmed("Stack:"), green(cfg.Language), green(cfg.Framework))
		fmt.Printf("  %s %s\n", dimmed("Component:"), bold(componentName))
		fmt.Println()

		// Step 2: Check if component is installed
		if !cfg.HasComponent(componentName) {
			fmt.Println(yellow("⚠"), "Component", cyan(componentName), "is not installed.")
			fmt.Println()
			fmt.Println(dimmed("  Installed components:"))
			if len(cfg.InstalledComponents) == 0 {
				fmt.Println(dimmed("    (none)"))
			} else {
				for _, comp := range cfg.InstalledComponents {
					fmt.Printf("    %s %s\n", dimmed("•"), cyan(comp.Name))
				}
			}
			return nil
		}

		// Step 3: Load manifest to find the files to remove
		manifest, err := registry.LoadManifest()
		if err != nil {
			return fmt.Errorf("failed to load registry: %w", err)
		}

		stackKey := cfg.Language + "_" + cfg.Framework
		component, exists := manifest.Components[componentName]
		if !exists {
			// Component exists in config but not in manifest — just remove from config
			fmt.Println(yellow("⚠"), "Component", cyan(componentName), "not found in registry.")
			fmt.Println(dimmed("  Removing from .forge.json only (cannot determine files to delete)."))
			cfg.RemoveComponent(componentName)
			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}
			fmt.Println()
			fmt.Println(green("✓"), "Removed", cyan(componentName), "from .forge.json")
			fmt.Println()
			return nil
		}

		impl, exists := component.Implementations[stackKey]
		if !exists {
			// No implementation for this stack — just remove from config
			cfg.RemoveComponent(componentName)
			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}
			fmt.Println()
			fmt.Println(green("✓"), "Removed", cyan(componentName), "from .forge.json")
			fmt.Println()
			return nil
		}

		// Step 4: Show files that will be deleted
		fmt.Println(bold("  Files to be removed:"))
		for _, file := range impl.Files {
			targetPath := filepath.Join(cfg.SourceDir, file.Target)
			if _, err := os.Stat(targetPath); err == nil {
				fmt.Printf("    %s %s\n", red("✗"), targetPath)
			} else {
				fmt.Printf("    %s %s %s\n", dimmed("·"), targetPath, dimmed("(not found)"))
			}
		}
		fmt.Println()

		// Step 5: Confirm removal (unless --force)
		if !removeForce {
			var confirm bool
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewConfirm().
						Title(fmt.Sprintf("Remove %s? This will delete the files above.", componentName)).
						Value(&confirm),
				),
			)

			if err := form.Run(); err != nil {
				return fmt.Errorf("prompt failed: %w", err)
			}

			if !confirm {
				fmt.Println(yellow("⚠"), "Cancelled. No files were removed.")
				return nil
			}
		}

		// Step 6: Delete files
		fmt.Println(dimmed("  Removing files..."))
		removedCount := 0
		for _, file := range impl.Files {
			targetPath := filepath.Join(cfg.SourceDir, file.Target)

			if _, err := os.Stat(targetPath); os.IsNotExist(err) {
				fmt.Printf("  %s %s %s\n", dimmed("·"), targetPath, dimmed("(already gone)"))
				continue
			}

			if err := os.Remove(targetPath); err != nil {
				fmt.Printf("  %s %s — %s\n", yellow("⚠"), targetPath, err)
			} else {
				fmt.Printf("  %s %s\n", green("✓"), targetPath)
				removedCount++
			}

			// Try to remove empty parent directory
			dir := filepath.Dir(targetPath)
			cleanEmptyDir(dir, cfg.SourceDir)
		}

		// Step 7: Update .forge.json
		cfg.RemoveComponent(componentName)
		if err := config.Save(cfg); err != nil {
			fmt.Println(yellow("⚠"), "Failed to update .forge.json:", err)
		}

		// Step 8: Summary
		fmt.Println()
		fmt.Printf("%s Component %s removed successfully! (%d files deleted)\n",
			green("✓"), cyan(componentName), removedCount)

		// Warn about dependencies that may have been installed
		if impl.InstallCmd != "" {
			fmt.Println()
			fmt.Println(yellow("💡"), bold("Note:"), "Dependencies installed by this component were NOT removed:")
			fmt.Printf("    %s %s\n", dimmed("packages:"), dimmed(fmt.Sprintf("%v", impl.Dependencies)))
			fmt.Println(dimmed("    You may want to remove them manually if no longer needed."))
		}

		fmt.Println()
		return nil
	},
}

// cleanEmptyDir removes a directory if it's empty.
// Walks up the tree until it hits the stopAt directory.
func cleanEmptyDir(dir string, stopAt string) {
	absDir, _ := filepath.Abs(dir)
	absStop, _ := filepath.Abs(stopAt)

	for absDir != absStop && absDir != filepath.Dir(absDir) {
		entries, err := os.ReadDir(absDir)
		if err != nil || len(entries) > 0 {
			return
		}
		os.Remove(absDir)
		absDir = filepath.Dir(absDir)
	}
}

func init() {
	removeCmd.Flags().BoolVarP(&removeForce, "force", "f", false, "Skip confirmation prompt")
	rootCmd.AddCommand(removeCmd)
}
