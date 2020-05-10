package main

import (
	"github.com/Terag/kadok/security"
	"github.com/bwmarrin/discordgo"
	"strings"
)

// ExecuteAction is the function type to implement with an Action
type ExecuteAction func(s *discordgo.Session, m *discordgo.MessageCreate) (string, error)

// ResolveAction is used to resolve the action that should be executed from a command.
// It returns the action to execute and the function to execute it as well.
func ResolveAction(rootAction *Action, call []string) (*Action, ExecuteAction) {
	// It first check a sub action with the name of the 1st parameter exist
	if len(call) > 0 {
		if subAction, ok := (*rootAction).SubActions[strings.ToUpper(call[0])]; ok {
			// If yes, move to this action for the resolve and reduce de number of parameters
			return ResolveAction(subAction, call[1:])
		}
	}
	// If no, returns the action with its associated execute action.
	return rootAction, MakeExecuteAction(rootAction, call)
}

// MakeExecuteAction from an action
func MakeExecuteAction(action *Action, parameters []string) ExecuteAction {
	// Check if there is only 1 parameter with the value "HELP"
	if len(parameters) == 1 && strings.ToUpper(parameters[0]) == "HELP" {
		// If yes, returns an ExecuteAction that returns Action's information
		return func(s *discordgo.Session, m *discordgo.MessageCreate) (string, error) {
			return action.Information, nil
		}
	}
	// If no, returns an execute action with the parameters initialized
	return func(s *discordgo.Session, m *discordgo.MessageCreate) (string, error) {
		return action.Execute(s, m, parameters)
	}
}

// Action structure containing information and functions regarding an action
type Action struct {
	// The permission required to use this action. Use the EmptyPermission if no permission is required
	Permission security.Permission
	// Information regarding the action. It is the message that should be displayed when requesting for help
	Information string
	// The sub actions in the hierarchy. They do not inherit the permission requirement
	SubActions map[string]*Action
	// The function to execute when the action is called
	Execute func(s *discordgo.Session, m *discordgo.MessageCreate, parameters []string) (string, error)
}

// GetPermission is used with the security module to check the permission of the action
func (action *Action) GetPermission() security.Permission {
	return action.Permission
}

// NotImplementedExecute to use if Execute is not implemented
var NotImplementedExecute = func(s *discordgo.Session, m *discordgo.MessageCreate, parameters []string) (string, error) {
	return "Oups ! Je sais pas faire !", nil
}

var (
	// PingAction for Kadok to respond pong
	PingAction = Action{
		security.GetCharacterList,
		"Kadok il dit Pong!",
		map[string]*Action{},
		func(s *discordgo.Session, m *discordgo.MessageCreate, parameters []string) (string, error) {
			return "À Kadoc ! À Kadoc ! Pong!", nil
		},
	}

	// PongAction for Kadok to respond ping
	PongAction = Action{
		security.GetCharacterList,
		"Kadok il dit Pong!",
		map[string]*Action{},
		func(s *discordgo.Session, m *discordgo.MessageCreate, parameters []string) (string, error) {
			return "À Kadoc ! À Kadoc ! Ping!", nil
		},
	}

	// GetCharactersAction to retrieve the list of available characters
	GetCharactersAction = Action{
		security.GetCharacterList,
		"C'est Kadok qui a pleins d'amis",
		map[string]*Action{},
		func(s *discordgo.Session, m *discordgo.MessageCreate, parameters []string) (string, error) {
			message := ""
			message += "\nLe caca des pigeons, c'est caca. Si tu parle d'un de mes amis, je te dirai ce qu'il a dit."
			message += "\nJ'ai pleins d'amis à Kaamelott:"
			for _, character := range Configuration.Characters.List {
				message += "\n- " + character.Name
			}
			return message, nil
		},
	}

	// RootAction is the first action call by Kadok for resolve
	RootAction = Action{
		security.GetHelp,
		"Tatan elle fait du flan, elle m'a aussi dit de dire des choses intelligentes si on m'appel: 'AKadok' \n'Kadok aqui' ? Je dis tous mes amis !",
		map[string]*Action{
			"AQUI": &GetCharactersAction,
			"PING": &PingAction,
			"PONG": &PongAction,
		},
		NotImplementedExecute,
	}
)
