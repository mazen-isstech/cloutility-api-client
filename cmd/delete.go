package cmd

import (
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "The delete subcommand",
	Long:  `This subcommand bundles together the different delete operations that are available.`,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
