package main

import (
	"github.com/bwmarrin/discordgo"
)

// GetUserRoles retrieves the list of roles names of the User who sent the message
func GetUserRoles(s *discordgo.Session, m *discordgo.MessageCreate) ([]string, error) {
	cGR, cMR := make(chan map[string]string), make(chan []string)
	cGRerr, cMRerr := make(chan error), make(chan error)

	// async calls to discord to retrieve Guild roles and GuildMember roles
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

	// Merge match Guild roles and User roles if no error is returned
	roles := []string{}
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
