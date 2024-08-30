package cmd

import (
	"fmt"
	"github.com/Gavin152/pldap/internal/ldaputil"
	"github.com/spf13/cobra"
	"os"
)

var username string
var password string

var authenticateCmd = &cobra.Command{
	Use:   "authenticate",
	Short: "Authenticate a user",
	Long:  `Authenticate a user by providing necessary credentials.`,
	Run: func(cmd *cobra.Command, args []string) {
		authenticated, err := ldaputil.AuthenticateUser(username, password)
		if err != nil {
			fmt.Printf("Error authenticating user: %v\n", err)
			os.Exit(1)
		}

		if authenticated {
			fmt.Println("User authenticated successfully!")
		} else {
			fmt.Println("User authentication failed!")
		}
	},
}

func init() {
	authenticateCmd.Flags().StringVarP(&username, "username", "u", "", "Username for authentication")
	authenticateCmd.Flags().StringVarP(&password, "password", "p", "", "Password for authentication")
	authenticateCmd.MarkFlagRequired("username")
	authenticateCmd.MarkFlagRequired("password")
	rootCmd.AddCommand(authenticateCmd)
}
