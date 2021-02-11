package main

import (
	"bytes"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Terag/kadok/security"
	"github.com/bwmarrin/discordgo"
)

// TplParams as the standard information passed to an action template
type TplParams struct {
	Message    discordgo.Message
	Infos      Infos
	Parameters []string
	Data       map[string]interface{}
}

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
	// Executed on call to format data passed to template
	GetData func(s *discordgo.Session, m *discordgo.MessageCreate, parameters []string) map[string]interface{}
	// Path to the action template
	Template string
}

func getFuncMap(parameters []string) template.FuncMap {
	return template.FuncMap(template.FuncMap{
		"withParams": func(param string) bool {
			for _, p := range parameters {
				if p == param {
					return true
				}
			}
			return false
		},
	})
}

// GetTemplate resolves template file
func (a *Action) GetTemplate(parameters []string) *template.Template {
	var fullPath = path.Join(Configuration.Templates, a.Template)
	return template.Must(template.New(filepath.Base(fullPath)).Funcs(getFuncMap(parameters)).ParseFiles(fullPath))
}

// Execute is the function to execute when the action is called
func (a *Action) Execute(s *discordgo.Session, m *discordgo.MessageCreate, parameters []string) (string, error) {
	var tpl bytes.Buffer
	var err = a.GetTemplate(parameters).Execute(&tpl, &TplParams{
		*m.Message,
		GetInfos(),
		parameters,
		a.GetData(s, m, parameters),
	})
	return tpl.String(), err
}

// GetPermission is used with the security module to check the permission of the action
func (a *Action) GetPermission() security.Permission {
	return a.Permission
}

// EmptyData is an helper to instanciate an action that doesn't rely on data
func EmptyData(s *discordgo.Session, m *discordgo.MessageCreate, parameters []string) map[string]interface{} {
	return make(map[string]interface{})
}

var (
	// StatusAction for kadok to respond with general information about the itself
	StatusAction = Action{
		security.EmptyPermission,
		"\nJe te dis ce que j'ai dans mon ventre !",
		map[string]*Action{},
		EmptyData,
		"status.tmpl",
	}

	// GetCharactersAction to retrieve the list of available characters
	GetCharactersAction = Action{
		security.GetCharacterList,
		"C'est Kadok qui a pleins d'amis",
		map[string]*Action{},
		func(s *discordgo.Session, m *discordgo.MessageCreate, parameters []string) map[string]interface{} {
			return map[string]interface{}{
				"Characters": Configuration.Characters.List,
			}
		},
		"aqui.tmpl",
	}

	// RootAction is the first action call by Kadok for resolve
	RootAction = Action{
		security.GetHelp,
		"Tatan elle fait du flan, elle m'a aussi dit de dire des choses intelligentes si on m'appel: 'AKadok' \n'Kadok aqui' ? Je dis tous mes amis !",
		map[string]*Action{
			"AQUI":  &GetCharactersAction,
			"TATAN": &StatusAction,
		},
		EmptyData,
		"404.tmpl",
	}
)
