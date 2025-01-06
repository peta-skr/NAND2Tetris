/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// assemblerCmd represents the assembler command
var assemblerCmd = &cobra.Command{
	Use:   "assembler",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
Args: func(cmd *cobra.Command, args []string) error {
	if len(args) != 1{
		return errors.New("require arg is only one")
	}

	return nil
},
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("assembler called")
		fmt.Println(args[0])

	},
}

func init() {
	rootCmd.AddCommand(assemblerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// assemblerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// assemblerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
