/*
Copyright © 2022 Victor Rouquette victor@rouquette.me

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
// This command does nothing and only return information regarding Kadok
var rootCmd = &cobra.Command{
	Use:   "kadok",
	Short: "Kadok is a Discord bot firstly developed for the Guild \"Les petits pedestres\". It aims to provide fun and useful functionalities for the Guild Members.",
	Long: `You like the French serie Kaamelott, own a Discord Guild and would like to add more
features to it ? Don't think twice, this bot is for you!

The current available features are:

- Quote characters from Kaamelott universe (and even more!)
- Manage features' permission based on Discord roles
- Allow users to join/leave groups and clans`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// init flags for rootCmd
func init() {
}
