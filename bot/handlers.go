// This folder contains handlers for managing discord events

package bot

import (
	"fmt"
	"strings"

	"github.com/Terag/kadok/characters"
	"github.com/Terag/kadok/security"
	"github.com/bwmarrin/discordgo"
)

// registerHandlers add the event handlers used by Kadok to interact with Discord
func registerHandlers(dg *discordgo.Session) error {
	dg.AddHandler(messageCreate)
	return nil
}

// messageCreate is the handler that manages message created events
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Ignore all messages that come from another guild
	if m.GuildID != Configuration.Guild.ID {
		return
	}

	// Prevent the bot from crashing in the event of the goroutine panicking
	// It also returns an "Oups problème !" to inform the user that an error happend
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Kadok paaaaaanic !")
			fmt.Println(r)
			_, err := s.ChannelMessageSend(m.ChannelID, "Oups problème !")
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	roles, err := GetUserRoles(s, m)
	if err != nil {
		fmt.Println("Error retrieving user roles")
		return
	}
	isGranted := security.MakeIsGranted(Configuration.Security.RolesHierarchy, roles)

	// Check if the first letters of the message match with the bot command prefix
	if len(m.Content) >= len(Configuration.Prefix) && strings.ToUpper(m.Content[:len(Configuration.Prefix)]) == strings.ToUpper(Configuration.Prefix) {
		call := strings.Fields(m.Content)
		action, executeAction := ResolveAction(&RootAction, call[1:])

		if isGranted(action) {
			response, err := executeAction(s, m)
			if err != nil {
				fmt.Println(err)
				_, err = s.ChannelMessageSend(m.ChannelID, "Oups problème !")
				if err != nil {
					fmt.Println(err)
				}
				return
			}
			if response != "" {
				_, err = s.ChannelMessageSend(m.ChannelID, response)
				if err != nil {
					fmt.Println(err)
				}
			}
			return
		}
		return
	}

	// Check if a character can be quoted then quote it
	if isGranted(security.CallCharacter) {
		quote, err := characters.GetQuoteFromMessage(Configuration.Characters.List, m.Content)
		if err != nil {
			return
		}
		if quote != "" {
			_, err = s.ChannelMessageSend(m.ChannelID, "> "+quote)
			if err != nil {
				fmt.Println(err)
			}
		}
		return
	}
}
