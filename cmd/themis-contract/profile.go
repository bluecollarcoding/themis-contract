package main

import (
	"fmt"
	"os"
	"strings"

	contract "github.com/informalsystems/themis-contract/pkg/themis-contract"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func profileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Profile management",
	}
	cmd.AddCommand(
		profileListCmd(),
		profileAddCmd(),
		profileRemoveCmd(),
		profileRenameCmd(),
		profileSetCmd(),
	)
	return cmd
}

func profileListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List existing profiles",
		Run: func(cmd *cobra.Command, args []string) {
			profiles, err := globalCtx.Profiles()
			if err != nil {
				log.Error().Err(err).Msg("Failed to load profiles")
				os.Exit(1)
			}
			if len(profiles) == 0 {
				log.Info().Msgf("No profiles configured yet")
				return
			}
			log.Info().Msgf("%d profile(s) available:", len(profiles))
			for _, profile := range profiles {
				log.Info().Msgf("- %s", profile.Display())
			}
		},
	}
}

func profileAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add [name] [signatureID]",
		Short: "Add a new profile",
		Long:  "Add a new profile with the given name and optionally specify a signature ID to use for the profile",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			sigID := ""
			if len(args) > 1 {
				sigID = args[1]
			}
			profile, err := globalCtx.AddProfile(args[0], sigID)
			if err != nil {
				log.Error().Msgf("Failed to add new profile: %s", err)
				os.Exit(1)
			}
			log.Info().Msgf("Added profile: %s", profile.Display())
		},
	}
}

func profileRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rm [id]",
		Short: "Remove a profile",
		Long:  "Remove the profile with the given ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := globalCtx.RemoveProfile(args[0]); err != nil {
				log.Error().Msgf("Failed to remove profile \"%s\": %s", args[0], err)
				os.Exit(1)
			}
			log.Info().Msgf("Successfully removed profile with ID \"%s\"", args[0])
		},
	}
}

func profileRenameCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mv [src-id] [dest-name]",
		Short: "Rename a profile",
		Long: `Rename the profile with the given ID to have the specified name (the new ID
will automatically be derived from the name)`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if err := globalCtx.RenameProfile(args[0], args[1]); err != nil {
				log.Error().Msgf("Failed to rename profile \"%s\": %s", args[0], err)
				os.Exit(1)
			}
			log.Info().Msgf("Successfully renamed profile with ID \"%s\" to \"%s\"", args[0], args[1])
		},
	}
}

func profileSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set [id] [param] [value]",
		Short: "Set a profile parameter value",
		Long: fmt.Sprintf(`Set a specific profile parameter to the given value

Valid profile parameter names include: %s`, strings.Join(contract.ValidProfileParamNames(), ",")),
		Args: cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			profile, err := globalCtx.GetProfileByID(args[0])
			if err != nil {
				log.Error().Msgf("Failed to load profile \"%s\": %s", args[0], err)
				os.Exit(1)
			}
			if err := profile.SetParam(args[1], args[2], globalCtx); err != nil {
				log.Error().Msgf("Failed to set parameter \"%s\" for profile \"%s\": %s", args[1], args[0], err)
				os.Exit(1)
			}
			if err := profile.Save(); err != nil {
				log.Error().Msgf("Failed to save profile \"%s\": %s", args[0], err)
				os.Exit(1)
			}
		},
	}
}
