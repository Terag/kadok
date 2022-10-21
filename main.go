// Kadok is a bot developed for the Discord Guild "Les petits pedestres". It aims to provide fun and useful functionnalities for the Guilde Members.
//
// The code of Kadok is open source, you can use it and modify it as long as you follow the GNU_v3 licence terms.
//
// The current available features are:
//
// 2. Quote characters from Kaamelott universe (and even more!)
//
// 3. Manage features' permission based on Discord roles
//
// 4. Manage user groups and clans
//
// This documentation is mainly to help developers to understand how the inside of Kadok works
// To find more on how to start and configure the bot, please visit the wiki page: https://github.com/terag/kadok/wiki

package main

import (
	"math/rand"
	"time"

	"github.com/terag/kadok/cmd"
)

func main() {
	// Generate a seed for the current instance of Kadok
	rand.Seed(time.Now().UnixNano())

	// Launch the CLI interface
	cmd.Execute()
}
