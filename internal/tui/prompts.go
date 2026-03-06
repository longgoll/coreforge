package tui

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/longgoll/forge-cli/internal/config"
	"github.com/longgoll/forge-cli/internal/registry"
)

// ──────────────────────────────────────────────────────────────
// Interactive Prompts using charmbracelet/huh
// ──────────────────────────────────────────────────────────────

// SelectLanguage shows an interactive menu to select a programming language.
func SelectLanguage(manifest *registry.Manifest) (string, error) {
	// Build options from manifest
	options := make([]huh.Option[string], 0, len(manifest.Languages))
	for key, lang := range manifest.Languages {
		label := fmt.Sprintf("%s", lang.Name)
		options = append(options, huh.NewOption(label, key))
	}

	var selected string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("🌐 Select a language").
				Description("Choose the programming language for your project").
				Options(options...).
				Value(&selected),
		),
	)

	err := form.Run()
	return selected, err
}

// SelectFramework shows an interactive menu to select a framework.
func SelectFramework(manifest *registry.Manifest, languageKey string) (string, error) {
	lang, exists := manifest.Languages[languageKey]
	if !exists {
		return "", fmt.Errorf("language %s not found in manifest", languageKey)
	}

	// Build options from frameworks
	options := make([]huh.Option[string], 0, len(lang.Frameworks))
	for key, fw := range lang.Frameworks {
		label := fmt.Sprintf("%s — %s", fw.Name, fw.Description)
		options = append(options, huh.NewOption(label, key))
	}

	// If only one framework, auto-select it
	if len(options) == 1 {
		return options[0].Value, nil
	}

	var selected string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("📦 Select a framework").
				Description(fmt.Sprintf("Choose a framework for %s", lang.Name)).
				Options(options...).
				Value(&selected),
		),
	)

	err := form.Run()
	return selected, err
}

// InputSourceDir asks the user for the source directory path.
func InputSourceDir(languageKey string) (string, error) {
	// Suggest default based on language
	defaultDir := "./"
	switch languageKey {
	case "nodejs":
		defaultDir = "./src"
	case "csharp":
		defaultDir = "./src"
	case "golang":
		defaultDir = "./"
	}

	var sourceDir string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("📁 Source directory").
				Description("Where is your main source code located?").
				Placeholder(defaultDir).
				Value(&sourceDir),
		),
	)

	err := form.Run()

	// Use default if empty
	if sourceDir == "" {
		sourceDir = defaultDir
	}

	return sourceDir, err
}

// ──────────────────────────────────────────────────────────────
// V2 Prompts — Database, ORM, Architecture
// ──────────────────────────────────────────────────────────────

// SelectDatabase shows an interactive menu to select a database.
func SelectDatabase() (string, error) {
	options := []huh.Option[string]{
		huh.NewOption("PostgreSQL", "postgresql"),
		huh.NewOption("MongoDB", "mongodb"),
		huh.NewOption("MySQL", "mysql"),
		huh.NewOption("SQLite", "sqlite"),
		huh.NewOption("None (I'll configure later)", "none"),
	}

	var selected string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("🗄️  Select a database").
				Description("Choose the database for your project").
				Options(options...).
				Value(&selected),
		),
	)

	err := form.Run()
	return selected, err
}

// SelectORM shows an interactive menu to select an ORM based on the language.
func SelectORM(languageKey string) (string, error) {
	var options []huh.Option[string]

	switch languageKey {
	case "nodejs":
		options = []huh.Option[string]{
			huh.NewOption("Prisma — Type-safe ORM with auto-generated client", "prisma"),
			huh.NewOption("Mongoose — MongoDB ODM for Node.js", "mongoose"),
			huh.NewOption("Sequelize — Multi-dialect SQL ORM", "sequelize"),
			huh.NewOption("None (I'll configure later)", "none"),
		}
	case "csharp":
		options = []huh.Option[string]{
			huh.NewOption("Entity Framework Core — Microsoft's recommended ORM", "efcore"),
			huh.NewOption("Dapper — Lightweight micro-ORM", "dapper"),
			huh.NewOption("None (I'll configure later)", "none"),
		}
	case "golang":
		options = []huh.Option[string]{
			huh.NewOption("GORM — Full-featured ORM for Go", "gorm"),
			huh.NewOption("sqlx — Extensions to database/sql", "sqlx"),
			huh.NewOption("None (I'll configure later)", "none"),
		}
	default:
		options = []huh.Option[string]{
			huh.NewOption("None", "none"),
		}
	}

	var selected string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("🔧 Select an ORM").
				Description("Choose an ORM/database library for your project").
				Options(options...).
				Value(&selected),
		),
	)

	err := form.Run()
	return selected, err
}

// SelectArchitecture shows an interactive menu to select a project architecture.
func SelectArchitecture(languageKey string) (string, error) {
	var options []huh.Option[string]

	switch languageKey {
	case "nodejs":
		options = []huh.Option[string]{
			huh.NewOption("MVC — Model-View-Controller (routes/controllers/services)", "mvc"),
			huh.NewOption("Feature-Based — Each feature in its own folder", "feature-based"),
			huh.NewOption("Minimal — Simple flat structure", "minimal"),
		}
	case "csharp":
		options = []huh.Option[string]{
			huh.NewOption("Clean Architecture — Domain/Application/Infrastructure layers", "clean"),
			huh.NewOption("MVC — Model-View-Controller", "mvc"),
			huh.NewOption("Minimal API — Lightweight single-file endpoints", "minimal"),
		}
	case "golang":
		options = []huh.Option[string]{
			huh.NewOption("Standard Layout — /cmd, /internal, /pkg", "mvc"),
			huh.NewOption("Feature-Based — Each feature in its own package", "feature-based"),
			huh.NewOption("Minimal — Simple flat structure", "minimal"),
		}
	default:
		options = []huh.Option[string]{
			huh.NewOption("MVC", "mvc"),
			huh.NewOption("Minimal", "minimal"),
		}
	}

	var selected string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("🏗️  Select project architecture").
				Description("Choose the architecture pattern for your project").
				Options(options...).
				Value(&selected),
		),
	)

	err := form.Run()
	return selected, err
}

// ──────────────────────────────────────────────────────────────
// Confirmation Prompts
// ──────────────────────────────────────────────────────────────

// ConfirmOverwrite asks if the user wants to overwrite existing .forge.json.
func ConfirmOverwrite() (bool, error) {
	var confirmed bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Overwrite existing .forge.json?").
				Description("This will replace your current configuration").
				Affirmative("Yes, overwrite").
				Negative("No, keep it").
				Value(&confirmed),
		),
	)

	err := form.Run()
	return confirmed, err
}

// ConfirmAction shows a generic confirmation prompt with custom title/description.
func ConfirmAction(title, description string) (bool, error) {
	var confirmed bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(title).
				Description(description).
				Affirmative("Yes").
				Negative("No").
				Value(&confirmed),
		),
	)

	err := form.Run()
	return confirmed, err
}

// ──────────────────────────────────────────────────────────────
// Conflict Resolution Prompts
// ──────────────────────────────────────────────────────────────

// SelectConflictAction shows a prompt when a file exists and has been modified.
// Returns one of: "overwrite", "skip", "backup"
func SelectConflictAction(filePath string) (string, error) {
	options := []huh.Option[string]{
		huh.NewOption("Overwrite — Replace with new version", "overwrite"),
		huh.NewOption("Skip     — Keep your current file", "skip"),
		huh.NewOption("Backup   — Save current as .bak, then overwrite", "backup"),
	}

	var selected string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(fmt.Sprintf("⚠ File conflict: %s", filePath)).
				Description("This file has been modified since installation. What would you like to do?").
				Options(options...).
				Value(&selected),
		),
	)

	err := form.Run()
	return selected, err
}

// ──────────────────────────────────────────────────────────────
// Non-Interactive Defaults (for --non-interactive flag)
// ──────────────────────────────────────────────────────────────

// GetDefaultDatabase returns the default database for non-interactive mode.
func GetDefaultDatabase() string {
	return "none"
}

// GetDefaultORM returns the default ORM for the given language in non-interactive mode.
func GetDefaultORM(languageKey string) string {
	return "none"
}

// GetDefaultArchitecture returns the default architecture in non-interactive mode.
func GetDefaultArchitecture() string {
	return "mvc"
}

// GetDefaultSourceDir returns the default source directory for the given language.
func GetDefaultSourceDir(languageKey string) string {
	switch languageKey {
	case "nodejs":
		return "./src"
	case "csharp":
		return "./src"
	case "golang":
		return "./"
	default:
		return "./"
	}
}

// ──────────────────────────────────────────────────────────────
// Blueprint Prompts
// ──────────────────────────────────────────────────────────────

// ConfirmBlueprintInstall shows the blueprint details and asks for confirmation.
func ConfirmBlueprintInstall(name string, includes []string) (bool, error) {
	description := fmt.Sprintf("This blueprint will install %d components:\n", len(includes))
	for _, inc := range includes {
		description += fmt.Sprintf("  • %s\n", inc)
	}

	var confirmed bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("Install blueprint: %s?", name)).
				Description(description).
				Affirmative("Yes, install all").
				Negative("No, cancel").
				Value(&confirmed),
		),
	)

	err := form.Run()
	return confirmed, err
}

// ──────────────────────────────────────────────────────────────
// Helper: ORM Compatibility Filter
// ──────────────────────────────────────────────────────────────

// FilterORMByDatabase returns ORM options compatible with the selected database.
func FilterORMByDatabase(languageKey, database string) []huh.Option[string] {
	allORMs := config.ValidORMs[languageKey]
	options := make([]huh.Option[string], 0)

	for _, orm := range allORMs {
		// Mongoose only works with MongoDB
		if orm == "mongoose" && database != "mongodb" {
			continue
		}
		// Prisma/Sequelize/EFCore/Dapper/GORM/sqlx don't work with MongoDB
		if database == "mongodb" && (orm == "prisma" || orm == "sequelize" || orm == "efcore" || orm == "dapper" || orm == "gorm" || orm == "sqlx") {
			continue
		}
		options = append(options, huh.NewOption(orm, orm))
	}

	// Always include "none"
	hasNone := false
	for _, o := range options {
		if o.Value == "none" {
			hasNone = true
			break
		}
	}
	if !hasNone {
		options = append(options, huh.NewOption("None", "none"))
	}

	return options
}
