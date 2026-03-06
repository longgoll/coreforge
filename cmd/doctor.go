package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/longgoll/forge-cli/internal/registry"
	"github.com/spf13/cobra"
)

// Tool check configuration
type toolCheck struct {
	Name    string
	Command string
	Args    []string
	Help    string
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check your development environment",
	Long: fmt.Sprintf(`%s

Verify that required tools are installed on your system.
Checks for language runtimes, package managers, and other dependencies
needed to work with forge-cli.`, bold("🩺 Environment Doctor")),

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cyan("⚡"), bold("Forge CLI — Doctor"))
		fmt.Println(dimmed("─────────────────────────────────────────"))
		fmt.Println()

		checks := []toolCheck{
			{
				Name:    "Node.js",
				Command: "node",
				Args:    []string{"--version"},
				Help:    "Install from https://nodejs.org/",
			},
			{
				Name:    "npm",
				Command: "npm",
				Args:    []string{"--version"},
				Help:    "Comes with Node.js",
			},
			{
				Name:    ".NET SDK",
				Command: "dotnet",
				Args:    []string{"--version"},
				Help:    "Install from https://dotnet.microsoft.com/",
			},
			{
				Name:    "Go",
				Command: "go",
				Args:    []string{"version"},
				Help:    "Install from https://go.dev/",
			},
			{
				Name:    "Git",
				Command: "git",
				Args:    []string{"--version"},
				Help:    "Install from https://git-scm.com/",
			},
		}

		allOk := true
		for _, check := range checks {
			version, err := getToolVersion(check.Command, check.Args)
			if err != nil {
				fmt.Printf("  %s %-12s %s\n", red("✗"), check.Name, dimmed(check.Help))
				allOk = false
			} else {
				fmt.Printf("  %s %-12s %s\n", green("✓"), check.Name, dimmed(version))
			}
		}

		fmt.Println()
		fmt.Printf("  %s %s\n", dimmed("OS:"), bold(runtime.GOOS+"/"+runtime.GOARCH))

		// Registry info
		fmt.Println()
		fmt.Println(bold("  Registry:"))
		if registry.UseRemote {
			fmt.Printf("  %s %s %s\n", green("✓"), "Mode:", cyan("remote"))
			fmt.Printf("  %s %s %s\n", dimmed(" "), "URL:", dimmed(registry.RemoteRegistryURL))

			// Cache info
			if cacheInfo := registry.GetCacheInfo(); cacheInfo != nil {
				age := time.Since(cacheInfo.CachedAt).Round(time.Second)
				fmt.Printf("  %s %s %s %s\n", green("✓"), "Cache:", green("active"),
					dimmed(fmt.Sprintf("(age: %s)", age)))
			} else {
				fmt.Printf("  %s %s %s\n", dimmed("·"), "Cache:", dimmed("empty"))
			}
		} else {
			fmt.Printf("  %s %s %s\n", green("✓"), "Mode:", cyan("local (mock-registry)"))
		}

		fmt.Println()

		if allOk {
			fmt.Println(green("✓"), bold("All tools are installed! You're ready to go."))
		} else {
			fmt.Println(yellow("⚠"), "Some tools are missing. Install them for full functionality.")
		}
		fmt.Println()
	},
}

func getToolVersion(command string, args []string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	// Clean up output (remove newlines)
	version := string(output)
	if len(version) > 50 {
		version = version[:50] + "..."
	}
	return trimRight(version), nil
}

func trimRight(s string) string {
	for len(s) > 0 && (s[len(s)-1] == '\n' || s[len(s)-1] == '\r' || s[len(s)-1] == ' ') {
		s = s[:len(s)-1]
	}
	return s
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
