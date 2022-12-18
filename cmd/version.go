/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/terag/kadok/internal/info"
)

// versionCmd represents the version command used to retrieve kadok version information
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Long:  `Display version information`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version: " + info.Version)
		fmt.Println("Url: " + info.URL)
		fmt.Println("BuildCommit: " + info.GitCommit)
		fmt.Println("BuildDate: " + info.BuildDate)
		fmt.Println("Go: " + info.GetInfo().GoVersion)
		fmt.Println("License: " + info.LicenseName)
		fmt.Println("LicenseUrl: " + info.LicenseURL)
		fmt.Println("Contributors: " + info.Contributors)
	},
}

// init flags for versionCmd
func init() {
	rootCmd.AddCommand(versionCmd)
}
