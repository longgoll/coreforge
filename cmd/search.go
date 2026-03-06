package cmd

import (
	"fmt"
	"strings"

	"github.com/longgoll/forge-cli/internal/config"
	"github.com/longgoll/forge-cli/internal/registry"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <keyword>",
	Short: "Search for components, schemas, and blueprints by keyword",
	Long: fmt.Sprintf(`%s

Search across all registry tiers (components, schemas, blueprints) by name,
description, tag, or category.

%s
  forge search auth
  forge search middleware
  forge search database
  forge search blueprint`, bold("🔍 Search Registry"), dimmed("Examples:")),

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		keyword := strings.ToLower(args[0])

		fmt.Println(cyan("⚡"), bold("Forge CLI — Search"))
		fmt.Println(dimmed("─────────────────────────────────────────"))
		fmt.Println()

		// Load manifest
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
			}
		}

		fmt.Printf("  %s %s\n", dimmed("Query:"), bold(keyword))
		fmt.Println()

		// Search across all tiers
		matches := []searchResult{}

		// Search components (Tier 2)
		for name, comp := range manifest.Components {
			matchReason := matchComponent(name, comp, keyword)
			if matchReason != "" {
				matches = append(matches, searchResult{
					Name:      name,
					Component: comp,
					Reason:    matchReason,
					Tier:      "component",
				})
			}
		}

		// Search schemas (Tier 3)
		for name, schema := range manifest.Schemas {
			matchReason := matchSchema(name, schema, keyword)
			if matchReason != "" {
				matches = append(matches, searchResult{
					Name: name,
					Component: registry.Component{
						Description:     schema.Description,
						Category:        schema.Category,
						Tags:            schema.Tags,
						Implementations: schema.Implementations,
					},
					Reason: matchReason,
					Tier:   "schema",
				})
			}
		}

		// Search blueprints (Tier 4)
		for name, bp := range manifest.Blueprints {
			matchReason := matchBlueprint(name, bp, keyword)
			if matchReason != "" {
				matches = append(matches, searchResult{
					Name: name,
					Component: registry.Component{
						Description:     bp.Description,
						Category:        bp.Category,
						Tags:            bp.Tags,
						Implementations: bp.Implementations,
					},
					Reason:    matchReason,
					Tier:      "blueprint",
					Includes:  bp.Includes,
				})
			}
		}

		if len(matches) == 0 {
			fmt.Println(yellow("⚠"), "No items found matching", cyan(keyword))
			fmt.Println()
			fmt.Println(dimmed("  Try a different keyword, or run"), cyan("forge list"), dimmed("to see all items."))
			fmt.Println()
			return nil
		}

		fmt.Printf(bold("  Found %d item(s):\n"), len(matches))
		fmt.Println()

		for _, match := range matches {
			// Status indicator based on current stack
			status := "  "
			if stackKey != "" {
				if _, exists := match.Component.Implementations[stackKey]; exists {
					status = green("✓ ")
				} else {
					status = red("✗ ")
				}
			}

			// Tier badge
			tierBadge := ""
			switch match.Tier {
			case "component":
				tierBadge = dimmed("[component]")
			case "schema":
				tierBadge = dimmed("[schema]")
			case "blueprint":
				tierBadge = dimmed("[blueprint]")
			}

			category := ""
			if match.Component.Category != "" {
				category = dimmed(fmt.Sprintf(" [%s]", match.Component.Category))
			}

			fmt.Printf("  %s %s — %s %s%s\n", status, cyan(match.Name), match.Component.Description, tierBadge, category)
			fmt.Printf("      %s %s\n", dimmed("match:"), dimmed(match.Reason))

			// Show tags
			if len(match.Component.Tags) > 0 {
				fmt.Printf("      %s %s\n", dimmed("tags:"), dimmed(strings.Join(match.Component.Tags, ", ")))
			}

			// Show included items for blueprints
			if len(match.Includes) > 0 {
				fmt.Printf("      %s %s\n", dimmed("includes:"), dimmed(strings.Join(match.Includes, ", ")))
			}

			// Show available stacks
			stacks := []string{}
			for key := range match.Component.Implementations {
				stacks = append(stacks, key)
			}
			fmt.Printf("      %s %s\n", dimmed("stacks:"), dimmed(strings.Join(stacks, ", ")))
			fmt.Println()
		}

		fmt.Println(dimmed("  Use"), cyan("forge add <name>"), dimmed("to install an item"))
		fmt.Println()

		return nil
	},
}

type searchResult struct {
	Name      string
	Component registry.Component
	Reason    string
	Tier      string   // "component", "schema", "blueprint"
	Includes  []string // For blueprints
}

// matchComponent checks if a component matches the search keyword.
// Returns the match reason or empty string if no match.
func matchComponent(name string, comp registry.Component, keyword string) string {
	// Match by name
	if strings.Contains(strings.ToLower(name), keyword) {
		return "name"
	}

	// Match by description
	if strings.Contains(strings.ToLower(comp.Description), keyword) {
		return "description"
	}

	// Match by category
	if strings.Contains(strings.ToLower(comp.Category), keyword) {
		return "category"
	}

	// Match by tags
	for _, tag := range comp.Tags {
		if strings.Contains(strings.ToLower(tag), keyword) {
			return fmt.Sprintf("tag: %s", tag)
		}
	}

	return ""
}

// matchSchema checks if a schema matches the search keyword.
func matchSchema(name string, schema registry.Schema, keyword string) string {
	if strings.Contains(strings.ToLower(name), keyword) {
		return "name"
	}
	if strings.Contains(strings.ToLower(schema.Description), keyword) {
		return "description"
	}
	if strings.Contains(strings.ToLower(schema.Category), keyword) {
		return "category"
	}
	for _, tag := range schema.Tags {
		if strings.Contains(strings.ToLower(tag), keyword) {
			return fmt.Sprintf("tag: %s", tag)
		}
	}
	return ""
}

// matchBlueprint checks if a blueprint matches the search keyword.
func matchBlueprint(name string, bp registry.Blueprint, keyword string) string {
	if strings.Contains(strings.ToLower(name), keyword) {
		return "name"
	}
	if strings.Contains(strings.ToLower(bp.Description), keyword) {
		return "description"
	}
	if strings.Contains(strings.ToLower(bp.Category), keyword) {
		return "category"
	}
	for _, tag := range bp.Tags {
		if strings.Contains(strings.ToLower(tag), keyword) {
			return fmt.Sprintf("tag: %s", tag)
		}
	}
	// Also search in included component names
	for _, inc := range bp.Includes {
		if strings.Contains(strings.ToLower(inc), keyword) {
			return fmt.Sprintf("includes: %s", inc)
		}
	}
	return ""
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
