package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	// "time"

	"github.com/dontizi/rlama/internal/domain"
	"github.com/dontizi/rlama/internal/repository"
	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage API profiles",
	Long:  `Create, list, and manage API profiles for different providers.`,
}

var profileAddCmd = &cobra.Command{
	Use:   "add [name] [provider] [api-key]",
	Short: "Add a new API profile",
	Long: `Add a new API profile for a specific provider.
Examples: 
  rlama profile add openai-work openai sk-...your-api-key...
  rlama profile add lmstudio openai-api none --base-url http://localhost:1234/v1`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		provider := args[1]
		apiKey := args[2]

		// Get base URL flag
		baseURL, _ := cmd.Flags().GetString("base-url")

		// Validate the provider
		switch provider {
		case "openai":
			// Official OpenAI API
		case "openai-api":
			// Generic OpenAI-compatible API
			if baseURL == "" {
				return fmt.Errorf("base-url is required for openai-api provider")
			}
		default:
			return fmt.Errorf("unsupported provider: %s. Supported providers: openai, openai-api", provider)
		}

		// Create the repository
		profileRepo := repository.NewProfileRepository()

		// Check if the profile already exists
		if profileRepo.Exists(name) {
			return fmt.Errorf("profile '%s' already exists", name)
		}

		// Create and save the profile
		profile := domain.NewAPIProfile(name, provider, apiKey)
		profile.BaseURL = baseURL
		if err := profileRepo.Save(profile); err != nil {
			return err
		}

		if baseURL != "" {
			fmt.Printf("Profile '%s' for '%s' (base URL: %s) added successfully.\n", name, provider, baseURL)
		} else {
			fmt.Printf("Profile '%s' for '%s' added successfully.\n", name, provider)
		}
		return nil
	},
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all API profiles",
	Long:  `Display a list of all configured API profiles.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		profileRepo := repository.NewProfileRepository()

		profiles, err := profileRepo.ListAll()
		if err != nil {
			return err
		}

		if len(profiles) == 0 {
			fmt.Println("No API profiles found.")
			return nil
		}

		fmt.Printf("Available API profiles (%d found):\n\n", len(profiles))

		// Use tabwriter to align the display
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tPROVIDER\tBASE URL\tCREATED ON\tLAST USED")

		for _, name := range profiles {
			profile, err := profileRepo.Load(name)
			if err != nil {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", name, "error", "error", "error", "error")
				continue
			}

			// Format dates
			createdAt := profile.CreatedAt.Format("2006-01-02 15:04:05")
			lastUsed := "never"
			if !profile.LastUsedAt.IsZero() {
				lastUsed = profile.LastUsedAt.Format("2006-01-02 15:04:05")
			}

			baseURL := profile.BaseURL
			if baseURL == "" {
				baseURL = "default"
			}

			// Hide the API key
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				profile.Name, profile.Provider, baseURL, createdAt, lastUsed)
		}
		w.Flush()

		return nil
	},
}

var profileDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete an API profile",
	Long:  `Delete an API profile by name.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		profileRepo := repository.NewProfileRepository()

		// Check if the profile exists
		if !profileRepo.Exists(name) {
			return fmt.Errorf("profile '%s' does not exist", name)
		}

		// Ask for confirmation
		fmt.Printf("Are you sure you want to delete profile '%s'? (y/n): ", name)
		var response string
		fmt.Scanln(&response)

		if response != "y" && response != "Y" {
			fmt.Println("Deletion cancelled.")
			return nil
		}

		// Delete the profile
		if err := profileRepo.Delete(name); err != nil {
			return err
		}

		fmt.Printf("Profile '%s' deleted successfully.\n", name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(profileCmd)
	profileCmd.AddCommand(profileAddCmd)
	profileCmd.AddCommand(profileListCmd)
	profileCmd.AddCommand(profileDeleteCmd)

	// Add flags for profile add command
	profileAddCmd.Flags().String("base-url", "", "Base URL for OpenAI-compatible endpoints (required for openai-api provider)")
}
