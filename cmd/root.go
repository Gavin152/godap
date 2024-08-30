package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{Use: "pldap"}

//rootCmd.AddCommand(authenticateCmd)
//rootCmd.AddCommand(groupsCmd)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
