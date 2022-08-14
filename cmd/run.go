/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/Terag/kadok/bot"
	"github.com/spf13/cobra"
)

var (
	token      string
	properties string
)

// runCmd represents the run command used to start Kadok bot
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run Kadok bot",
	Long: `Run Kadok bot. You must provide a valid Discord token with -t for the bot to access your Discord guild
    You can generate such token from Discord developer portal: https://discord.com/developers/applications

    You can also provide the parth to a properties file with -p. By default, kadok will look at "config/properties.yaml". Below is the default value:

--- yaml
prefix: "kadok"
guild:
  name: "Les petits pedestres"
security:
  roles: "config/roles.yaml"
characters:
  folder: "assets/characters"
templates: "assets/templates"
---

    Where:

      - prefix:         string  is the prefix that will be used by kadok to detect when it is called. ex: "kadok groupe liste"
        (REQUIRED)

      - guild.name:     string  Kadok is a mono guild bot. It means that kadok is designed to interact with only one Discord guild.
        (REQUIRED)              This value is mandatory and used by the bot to filter evens if kadok is invited to multiple guilds.

      - security.roles: string  path to the file container roles configuration for Kadok bot. Those roles are directly correlated
        (REQUIRED)              to the roles defined in the Discord guild. The format for this file is:

--- yaml
roles:
- name: member
  parent: null
  permissions:
    - CallCharacter
    - GetHelp
    - GetCharacterList
- name: moderator
  parent: member
  permissions: []
- name: administrator
  parent: moderator
  permissions: []
- name: Pédèstres seniors
  parent: member
  permissions: []
  type: Clan
  description: Ils sont vieux et ils vont en marchant !
- name: Semi Croustillants
  parent: member
  permissions: []
  type: Clan
  description: Pas complètement mais a moitie quand même !
- name: PoliticalLeader
  parent: member
  permissions: []
  type: Group
  description: (🌏-grande-strategie) Moins langue de bois, et aucune scrupule politique
- name: Roliste
  parent: member
  permissions: []
  type: Group
  description: (🌲La foret 🌲) Du papier et des dés, que demande le peuple ?
---

        The structure of a Role is:

          - name:           string    the name of the role as defined on Discord. Kadok uses the name to fin the role on Discord.
            (REQUIRED)                The role does not have to exist on the guild, this is useful if you want to do permissions inheritance.

          - parent:         string    the parent role. The children role will inherit of the parent permissions. You can define a role member
            (REQUIRED)                and all the roles with member as parent will get the parent's permissions.

          - permissions:    string[]  list of permissions. The permissions to List/Join/Leave groups and clans are by design accessible to everyone.
            (REQUIRED)                The available permissions are:

                                        - CallCharacter: for Kadok to quote characters on user's message
                                        - GetHelp: to ask kadok about its available commands
                                        - GetCharacterList: get the list of characters that can be quoted

          - description:    string    description about the role. Displayed by kadok when listing available groups and clans.

          - type:           enum      Type of roles. Available values: Default (default), Group and Clan. A Default role is not displayed
            (OPTIONAL)                to the user, therefore a user can't interact with it through Kadok. Groups and Clans can be joined and leaved
                                      by users at any moment. A user can be in multiple group but can belong to only one clan.

      - charaters.folder: string  path to the folder containing characters that kadok can quotes.
        (REQUIRED)                The structure of a character file is:

--- json
{
    "name": "kadoc",
    "sentences": [
        "À Kadoc ! À Kadoc !",
        "Tatan, elle fait des flans.",
        "Les pattes de canaaaaaaaaaaaaaaaaaaaaaaaaaaaaard !",
        "Elle est où la poulette ?",
        "Le caca des pigeons c'est caca, faut pas manger.",
        "Il est où Kadoc ? Il est bien caché ?",
    ]
}
---

        The structure of a Character is:

          - name:           string    the name of the character that will used to detect when to quote the character.
            (REQUIRED)

          - sentences:      string[]  the sentences related to the character that will be used when quoting it.
            (REQUIRED)=

      - templates:        string  path to the folder containing Kadok templates. See the current available templates.
        (REQUIRED)                Then must be all present for Kadok to work properly.`,
	Run: func(cmd *cobra.Command, args []string) {
		bot.Run(token, properties)
	},
}

// init flags for runCmd
func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&token, "token", "t", "", "Bot's token to connect to discord")
	runCmd.MarkFlagRequired("token")
	runCmd.Flags().StringVarP(&properties, "properties", "p", "config/properties.yaml", "Path to the properties file")
}
