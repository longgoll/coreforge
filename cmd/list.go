package cmd

import (
	"fmt"
	"sort"

	"github.com/longgoll/forge-cli/internal/config"
	"github.com/longgoll/forge-cli/internal/registry"
	"github.com/spf13/cobra"
)

var (
	listInstalled bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available or installed components",
	Long: fmt.Sprintf(`%s

Show all items available in the registry (components, schemas, blueprints),
or list the ones already installed in your project.

%s
  forge list                  # List all available items (grouped by tier)
  forge list --installed      # List installed components`, bold("📋 List Registry Items"), dimmed("Examples:")),

	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(cyan("⚡"), bold("Forge CLI — Registry"))
		fmt.Println(dimmed("─────────────────────────────────────────"))
		fmt.Println()

		if listInstalled {
			return listInstalledComponents()
		}
		return listAvailableItems()
	},
}

func listAvailableItems() error {
	manifest, err := registry.LoadManifest()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	// Determine current stack (if config exists)
	stackKey := ""
	if config.Exists() {
		cfg, err := config.Load()
		if err == nil {
			stackKey = cfg.GetStackKey()
			fmt.Printf("  %s %s + %s\n", dimmed("Stack:"), green(cfg.Language), green(cfg.Framework))
			if cfg.IsV2() && cfg.Database != "" && cfg.Database != "none" {
				fmt.Printf("  %s %s (%s)\n", dimmed("Database:"), green(cfg.Database), green(cfg.ORM))
			}
			fmt.Println()
		}
	}

	// Show schema version
	if manifest.IsV2() {
		fmt.Printf("  %s %s\n", dimmed("Registry schema:"), dimmed("v"+manifest.SchemaVersion))
	} else {
		fmt.Printf("  %s %s\n", dimmed("Registry schema:"), dimmed("v1 (legacy)"))
	}
	fmt.Println()

	// ── Tier 1: Foundations ────────────────────────
	if manifest.HasFoundations() {
		fmt.Println(bold("  🏗️  Foundations"), dimmed(fmt.Sprintf("(%d available)", manifest.GetFoundationCount())))
		fmt.Println()

		for name, f := range manifest.Foundations {
			langMatch := ""
			if stackKey != "" {
				fStack := f.Language + "_" + f.Framework
				if fStack == stackKey {
					langMatch = green("✓ ")
				} else {
					langMatch = dimmed("  ")
				}
			}
			fmt.Printf("  %s %s — %s\n", langMatch, cyan(name), f.Description)
			fmt.Printf("      %s %s + %s\n", dimmed("stack:"), dimmed(f.Language), dimmed(f.Framework))
			if f.Architecture != "" {
				fmt.Printf("      %s %s\n", dimmed("arch:"), dimmed(f.Architecture))
			}
		}
		fmt.Println()
	}

	// ── Tier 2: Components ────────────────────────
	fmt.Println(bold("  📦 Components"), dimmed(fmt.Sprintf("(%d available)", manifest.GetComponentCount())))
	fmt.Println()

	// Sort component names for consistent output
	compNames := make([]string, 0, len(manifest.Components))
	for name := range manifest.Components {
		compNames = append(compNames, name)
	}
	sort.Strings(compNames)

	for _, name := range compNames {
		comp := manifest.Components[name]
		available := "  "
		if stackKey != "" {
			if _, exists := comp.Implementations[stackKey]; exists {
				available = green("✓ ")
			} else {
				available = red("✗ ")
			}
		}

		category := ""
		if comp.Category != "" {
			category = dimmed(fmt.Sprintf(" [%s]", comp.Category))
		}

		fmt.Printf("  %s %s — %s%s\n", available, cyan(name), comp.Description, category)

		// Show available stacks for this component
		stacks := []string{}
		for key := range comp.Implementations {
			stacks = append(stacks, key)
		}
		sort.Strings(stacks)
		fmt.Printf("      %s %s\n", dimmed("stacks:"), dimmed(fmt.Sprintf("%v", stacks)))
	}
	fmt.Println()

	// ── Tier 3: Schemas ────────────────────────
	if manifest.HasSchemas() {
		fmt.Println(bold("  🗄️  Schemas"), dimmed(fmt.Sprintf("(%d available)", manifest.GetSchemaCount())))
		fmt.Println()

		for name, schema := range manifest.Schemas {
			available := "  "
			if stackKey != "" {
				if _, exists := schema.Implementations[stackKey]; exists {
					available = green("✓ ")
				} else {
					available = red("✗ ")
				}
			}

			category := ""
			if schema.Category != "" {
				category = dimmed(fmt.Sprintf(" [%s]", schema.Category))
			}

			fmt.Printf("  %s %s — %s%s\n", available, cyan(name), schema.Description, category)

			stacks := []string{}
			for key := range schema.Implementations {
				stacks = append(stacks, key)
			}
			if len(stacks) > 0 {
				sort.Strings(stacks)
				fmt.Printf("      %s %s\n", dimmed("stacks:"), dimmed(fmt.Sprintf("%v", stacks)))
			}
		}
		fmt.Println()
	}

	// ── Tier 4: Blueprints ────────────────────────
	if manifest.HasBlueprints() {
		fmt.Println(bold("  🔮 Blueprints"), dimmed(fmt.Sprintf("(%d available)", manifest.GetBlueprintCount())))
		fmt.Println()

		for name, bp := range manifest.Blueprints {
			available := "  "
			if stackKey != "" {
				if _, exists := bp.Implementations[stackKey]; exists {
					available = green("✓ ")
				} else {
					available = red("✗ ")
				}
			}

			fmt.Printf("  %s %s — %s\n", available, cyan(name), bp.Description)
			fmt.Printf("      %s %s\n", dimmed("includes:"), dimmed(fmt.Sprintf("%v", bp.Includes)))

			stacks := []string{}
			for key := range bp.Implementations {
				stacks = append(stacks, key)
			}
			if len(stacks) > 0 {
				sort.Strings(stacks)
				fmt.Printf("      %s %s\n", dimmed("stacks:"), dimmed(fmt.Sprintf("%v", stacks)))
			}
		}
		fmt.Println()
	}

	// Footer
	fmt.Println(dimmed("  Use"), cyan("forge add <name>"), dimmed("to install a component, schema, or blueprint"))
	fmt.Println()

	return nil
}

func listInstalledComponents() error {
	if !config.Exists() {
		fmt.Println(yellow("⚠"), "No", cyan(".forge.json"), "found. Run", cyan("forge init"), "first.")
		return nil
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Printf("  %s %s + %s\n", dimmed("Stack:"), green(cfg.Language), green(cfg.Framework))
	if cfg.IsV2() {
		if cfg.Database != "" && cfg.Database != "none" {
			fmt.Printf("  %s %s\n", dimmed("Database:"), green(cfg.Database))
		}
		if cfg.ORM != "" && cfg.ORM != "none" {
			fmt.Printf("  %s %s\n", dimmed("ORM:"), green(cfg.ORM))
		}
		if cfg.Architecture != "" {
			fmt.Printf("  %s %s\n", dimmed("Architecture:"), green(cfg.Architecture))
		}
	}
	fmt.Println()

	if len(cfg.InstalledComponents) == 0 {
		fmt.Println(dimmed("  No components installed yet."))
		fmt.Println()
		fmt.Println(dimmed("  Use"), cyan("forge add <name>"), dimmed("to install a component"))
	} else {
		fmt.Println(bold("  Installed components:"))
		fmt.Println()
		for _, comp := range cfg.InstalledComponents {
			fmt.Printf("  %s %s %s\n",
				green("✓"),
				cyan(comp.Name),
				dimmed(fmt.Sprintf("(installed: %s)", comp.InstalledAt)),
			)
		}
	}

	fmt.Println()
	return nil
}

func init() {
	listCmd.Flags().BoolVar(&listInstalled, "installed", false, "Show only installed components")
	rootCmd.AddCommand(listCmd)
}
