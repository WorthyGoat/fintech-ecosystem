package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Display the current logged in user",
	Run: func(cmd *cobra.Command, args []string) {
		email := viper.GetString("email")
		key := viper.GetString("api_key")
		if email == "" {
			fmt.Println("Not logged in.")
			return
		}
		fmt.Printf("Logged in as: %s\n", email)
		if key != "" {
			fmt.Printf("API Key: %s...%s\n", key[:7], key[len(key)-4:])
		}
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
