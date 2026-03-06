package cmd

import (
	"fmt"

	"github.com/longgoll/forge-cli/internal/config"
	"github.com/longgoll/forge-cli/internal/registry"
	"github.com/longgoll/forge-cli/internal/tui"
	"github.com/spf13/cobra"
)

var (
	// Non-interactive flags for CI/scripting
	initLanguage     string
	initFramework    string
	initSourceDir    string
	initDatabase     string
	initORM          string
	initArchitecture string
	nonInteractive   bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new project or recognize an existing one",
	Long: fmt.Sprintf(`%s

Interactively set up a new backend project or configure an existing one
for use with forge-cli. This command will:

  1. Ask you to choose a language and framework
  2. Select a database and ORM
  3. Choose a project architecture pattern
  4. Configure the project source directory
  5. Generate a %s file in the current directory

If the current directory already has code, forge will help you create
a config file to enable the "add" command.

%s
  forge init
  forge init --non-interactive --language nodejs --framework express
  forge init --language golang --framework gin --database postgresql --orm gorm`,
		bold("🚀 Initialize a forge project"), cyan(".forge.json"), dimmed("Examples:")),

	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(cyan("⚡"), bold("Forge CLI — Project Setup"))
		fmt.Println(dimmed("─────────────────────────────────────────"))
		fmt.Println()

		// Step 1: Check if .forge.json already exists
		if config.Exists() {
			cfg, err := config.Load()
			if err == nil {
				fmt.Println(yellow("⚠"), "A", cyan(".forge.json"), "already exists in this directory:")
				fmt.Printf("   Language:     %s\n", green(cfg.Language))
				fmt.Printf("   Framework:    %s\n", green(cfg.Framework))
				fmt.Printf("   Source:       %s\n", green(cfg.SourceDir))
				if cfg.Database != "" {
					fmt.Printf("   Database:     %s\n", green(cfg.Database))
				}
				if cfg.ORM != "" {
					fmt.Printf("   ORM:          %s\n", green(cfg.ORM))
				}
				if cfg.Architecture != "" {
					fmt.Printf("   Architecture: %s\n", green(cfg.Architecture))
				}
				fmt.Println()

				if nonInteractive {
					fmt.Println(dimmed("Non-interactive mode: overwriting existing config."))
				} else {
					overwrite, err := tui.ConfirmOverwrite()
					if err != nil || !overwrite {
						fmt.Println(dimmed("Aborted. Existing config unchanged."))
						return nil
					}
				}
				fmt.Println()
			}
		}

		// Step 2: Load manifest (registry)
		manifest, err := registry.LoadManifest()
		if err != nil {
			return fmt.Errorf("failed to load registry: %w", err)
		}

		// Step 3: Language selection
		languageKey := initLanguage
		if languageKey == "" {
			if nonInteractive {
				return fmt.Errorf("--language is required in non-interactive mode")
			}
			languageKey, err = tui.SelectLanguage(manifest)
			if err != nil {
				return fmt.Errorf("language selection cancelled: %w", err)
			}
		}

		// Step 4: Framework selection
		frameworkKey := initFramework
		if frameworkKey == "" {
			if nonInteractive {
				// Auto-select first available framework
				keys := manifest.GetFrameworkKeys(languageKey)
				if len(keys) == 0 {
					return fmt.Errorf("no frameworks found for language: %s", languageKey)
				}
				frameworkKey = keys[0]
			} else {
				frameworkKey, err = tui.SelectFramework(manifest, languageKey)
				if err != nil {
					return fmt.Errorf("framework selection cancelled: %w", err)
				}
			}
		}

		// Step 5: Database selection (V2)
		database := initDatabase
		if database == "" {
			if nonInteractive {
				database = tui.GetDefaultDatabase()
			} else {
				database, err = tui.SelectDatabase()
				if err != nil {
					return fmt.Errorf("database selection cancelled: %w", err)
				}
			}
		}

		// Step 6: ORM selection (V2)
		orm := initORM
		if orm == "" {
			if nonInteractive {
				orm = tui.GetDefaultORM(languageKey)
			} else {
				if database != "none" {
					orm, err = tui.SelectORM(languageKey)
					if err != nil {
						return fmt.Errorf("ORM selection cancelled: %w", err)
					}
				} else {
					orm = "none"
				}
			}
		}

		// Step 7: Architecture selection (V2)
		architecture := initArchitecture
		if architecture == "" {
			if nonInteractive {
				architecture = tui.GetDefaultArchitecture()
			} else {
				architecture, err = tui.SelectArchitecture(languageKey)
				if err != nil {
					return fmt.Errorf("architecture selection cancelled: %w", err)
				}
			}
		}

		// Step 8: Source directory
		sourceDir := initSourceDir
		if sourceDir == "" {
			if nonInteractive {
				sourceDir = tui.GetDefaultSourceDir(languageKey)
			} else {
				sourceDir, err = tui.InputSourceDir(languageKey)
				if err != nil {
					return fmt.Errorf("source directory input cancelled: %w", err)
				}
			}
		}

		// Download template if available
		lang, ok := manifest.Languages[languageKey]
		if ok {
			framework, ok := lang.Frameworks[frameworkKey]
			if ok && framework.TemplateURL != "" {
				fmt.Println()
				fmt.Println(dimmed("  Downloading project template..."))
				err := downloadAndExtractZip(framework.TemplateURL, ".")
				if err != nil {
					fmt.Println(yellow("⚠"), "Failed to download template:", err)
				} else {
					fmt.Println(green("  ✓"), "Template extracted successfully")
				}
			}
		}

		// Step 9: Create .forge.json (V2)
		cfg := &config.ForgeConfig{
			// V1 fields
			Language:            languageKey,
			Framework:           frameworkKey,
			SourceDir:           sourceDir,
			InstalledComponents: []config.InstalledComponent{},

			// V2 fields
			ConfigVersion: "2.0",
			Database:      database,
			ORM:           orm,
			Architecture:  architecture,
		}

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		// Step 10: Success message
		fmt.Println()
		fmt.Println(green("✓"), bold("Project initialized successfully!"))
		fmt.Println()
		fmt.Printf("  %s %s + %s\n", dimmed("Stack:"), green(languageKey), green(frameworkKey))
		fmt.Printf("  %s %s\n", dimmed("Source:"), green(sourceDir))
		if database != "none" {
			fmt.Printf("  %s %s\n", dimmed("Database:"), green(database))
		}
		if orm != "none" {
			fmt.Printf("  %s %s\n", dimmed("ORM:"), green(orm))
		}
		fmt.Printf("  %s %s\n", dimmed("Architecture:"), green(architecture))
		fmt.Printf("  %s %s\n", dimmed("Config:"), green(".forge.json"))
		fmt.Println()
		fmt.Println(dimmed("Next step:"), "run", cyan("forge add <component>"), "to add features")
		fmt.Println(dimmed("Example: "), cyan("forge add error-handler"))
		fmt.Println()

		return nil
	},
}

func init() {
	// Non-interactive flags
	initCmd.Flags().BoolVar(&nonInteractive, "non-interactive", false, "Run without interactive prompts (requires --language)")
	initCmd.Flags().StringVar(&initLanguage, "language", "", "Programming language (nodejs, csharp, golang)")
	initCmd.Flags().StringVar(&initFramework, "framework", "", "Framework (express, dotnet-webapi, gin)")
	initCmd.Flags().StringVar(&initSourceDir, "source", "", "Source directory path")
	initCmd.Flags().StringVar(&initDatabase, "database", "", "Database (postgresql, mongodb, mysql, sqlite, none)")
	initCmd.Flags().StringVar(&initORM, "orm", "", "ORM (prisma, mongoose, efcore, gorm, none)")
	initCmd.Flags().StringVar(&initArchitecture, "architecture", "", "Architecture (mvc, feature-based, clean, minimal)")

	rootCmd.AddCommand(initCmd)
}
