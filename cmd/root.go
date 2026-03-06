package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/longgoll/forge-cli/internal/registry"
	"github.com/spf13/cobra"
)

// ──────────────────────────────────────────────────────────────
// ASCII Art Banner
// ──────────────────────────────────────────────────────────────
const banner = `
   ______                        ____ _     ___ 
  / ____/___  _________ ____    / __ \ |   /  _/
 / /_  / __ \/ ___/ __ '/ _ \  / /  \/ |   / /  
/ __/ / /_/ / /  / /_/ /  __/ / /___/ |__/ /   
/_/    \____/_/   \__, /\___/  \____/|____/___/  
                 /____/                          
`

var (
	// Version info — injected at build time via ldflags
	Version   = "dev"
	BuildDate = "unknown"

	// Global flags
	useRemoteFlag bool

	// Color helpers
	cyan   = color.New(color.FgCyan, color.Bold).SprintFunc()
	green  = color.New(color.FgGreen, color.Bold).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	red    = color.New(color.FgRed, color.Bold).SprintFunc()
	bold   = color.New(color.Bold).SprintFunc()
	dimmed = color.New(color.Faint).SprintFunc()
)

// rootCmd is the base command when called without subcommands.
var rootCmd = &cobra.Command{
	Use:   "forge",
	Short: "Universal Backend CLI — scaffolding for any stack",
	Long: fmt.Sprintf(`%s
%s is a universal backend scaffolding tool inspired by shadcn/ui.

Generate projects, add production-ready components (error handlers,
JWT auth, loggers…) to any backend stack with a single command.
Code is copied into YOUR project — no vendor lock-in.

Supported stacks: %s, %s, %s`,
		cyan(banner),
		bold("forge"),
		green("Node.js (Express)"),
		green("C# (.NET Web API)"),
		green("Golang (Gin)"),
	),

	// Apply global flags before any subcommand runs
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		registry.UseRemote = useRemoteFlag
	},

	// If no subcommand is provided, print help
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute is the entry point called from main.go.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, red("Error:"), err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVar(&useRemoteFlag, "remote", false, "Use remote GitHub registry instead of local")

	// Version template
	rootCmd.Version = Version
	rootCmd.SetVersionTemplate(fmt.Sprintf(
		"%s forge-cli %s (built %s)\n",
		cyan("⚡"),
		bold(Version),
		dimmed(BuildDate),
	))
}
