package cmd

import "github.com/spf13/cobra"

func CreateRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "lofictl",
		Short: `lofictl is a CLI client of transforming a track to a low fidelity track`,
		Long:  `lofictl is a CLI built by D.K for transforming music to low fidelity ones, currently supports mp3s `,
	}

	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "apply",
		RunE:  applyCommand,
	}
	applyCmd.Flags().StringP("file", "f", "", "Track name to modify")
	applyCmd.Flags().Bool("boost", false, "Enable bass boost")
	applyCmd.Flags().Float64P("ratio", "r", 1, "Speed ratio for track")
	applyCmd.MarkFlagRequired("file")
	rootCmd.AddCommand(applyCmd)

	return rootCmd
}
