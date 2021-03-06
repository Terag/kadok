package main

import (
	"github.com/bwmarrin/discordgo"
)

// GetUserRoles retrieves the list of roles names of the User who sent the message.
func GetUserRoles(s *discordgo.Session, m *discordgo.MessageCreate) ([]string, error) {
	cGR, cMR := make(chan map[string]string), make(chan []string)
	cGRerr, cMRerr := make(chan error), make(chan error)

	// goroutine to retrieve Guild roles from the sender's Guild.
	go func(s *discordgo.Session, guildID string, c chan map[string]string, cerr chan error) {
		roles, err := s.GuildRoles(guildID)
		if err != nil {
			cerr <- err
		} else {
			rolesMap := make(map[string]string)
			for _, role := range roles {
				rolesMap[role.ID] = role.Name
			}
			c <- rolesMap
		}
		close(c)
		close(cerr)
	}(s, m.GuildID, cGR, cGRerr)

	// goroutine to retrieve Guild roles' Ids of the sender.
	go func(s *discordgo.Session, guildID string, userID string, c chan []string, cerr chan error) {
		member, err := s.GuildMember(guildID, userID)
		if err != nil {
			cerr <- err
		} else {
			c <- member.Roles
		}
		close(c)
		close(cerr)
	}(s, m.GuildID, m.Author.ID, cMR, cMRerr)

	// Merge Guild roles and User roles into a list of user's roles names.
	roles := make([]string, 0)
	select {
	case guildRoles := <-cGR:
		select {
		case memberRoles := <-cMR:
			for _, roleID := range memberRoles {
				roles = append(roles, guildRoles[roleID])
			}
		case err := <-cMRerr:
			return nil, err
		}
	case err := <-cGRerr:
		return nil, err
	}
	return roles, nil
}
