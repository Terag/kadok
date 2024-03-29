package bot

import (
	"bytes"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/bwmarrin/discordgo"
	"github.com/terag/kadok/internal/http"
	"github.com/terag/kadok/internal/info"
	"github.com/terag/kadok/internal/utils"
	"github.com/terag/kadok/pkg/security"
)

// TplParams as the standard information passed to an action template
type TplParams struct {
	Message       discordgo.Message
	Info          info.Info
	Configuration Properties
	Parameters    []string
	Data          map[string]interface{}
}

// ExecuteAction is the function type to implement with an Action
type ExecuteAction func(bc *BotContext, s *discordgo.Session, m *discordgo.MessageCreate) (string, error)

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
	if len(parameters) == 1 && strings.ToUpper(parameters[0]) == "AIDE" {
		// If yes, returns an ExecuteAction that returns Action's information
		return func(bc *BotContext, s *discordgo.Session, m *discordgo.MessageCreate) (string, error) {
			return action.Information, nil
		}
	}
	// If no, returns an execute action with the parameters initialized
	return func(bc *BotContext, s *discordgo.Session, m *discordgo.MessageCreate) (string, error) {
		return action.Execute(bc, s, m, parameters)
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
	GetData func(bc *BotContext, s *discordgo.Session, m *discordgo.MessageCreate, parameters []string) map[string]interface{}
	// Path to the action template
	Template string
}

// return a FuncMap available within all our templates as helpers
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
	// In order to pass the function map to the file you need to already have an instance of template.
	// In order for the template to be valid it needs to have the file name as template name which explains
	// the New with file base then the functions map then the file parsing.
	return template.Must(template.New(filepath.Base(fullPath)).Funcs(getFuncMap(parameters)).ParseFiles(fullPath))
}

// Execute is the function to execute when the action is called
func (a *Action) Execute(bc *BotContext, s *discordgo.Session, m *discordgo.MessageCreate, parameters []string) (string, error) {
	var tpl bytes.Buffer
	var err = a.GetTemplate(parameters).Execute(&tpl, &TplParams{
		*m.Message,
		info.GetInfo(),
		Configuration,
		parameters,
		a.GetData(bc, s, m, parameters),
	})
	return tpl.String(), err
}

// GetPermission is used with the security module to check the permission of the action
func (a *Action) GetPermission() security.Permission {
	return a.Permission
}

// EmptyData is an helper to instanciate an action that doesn't rely on data
func EmptyData(bc *BotContext, s *discordgo.Session, m *discordgo.MessageCreate, parameters []string) map[string]interface{} {
	return make(map[string]interface{})
}

var (
	// StatusAction for kadok to respond with general information about the itself
	StatusAction = Action{
		security.EmptyPermission,
		"\nJe te dis ce que j'ai dans mon ventre ! `plus` si t'en veux plus !",
		map[string]*Action{},
		EmptyData,
		"status.tmpl",
	}

	// GetCharactersAction to retrieve the list of available characters
	GetCharactersAction = Action{
		security.GetCharacterList,
		"C'est Kadok qui a pleins d'amis",
		map[string]*Action{},
		func(bc *BotContext, s *discordgo.Session, m *discordgo.MessageCreate, parameters []string) map[string]interface{} {
			return map[string]interface{}{
				"Characters": Configuration.Characters.List,
			}
		},
		"aqui.tmpl",
	}

	GroupRootAction = Action{
		security.EmptyPermission,
		"Toutes les actions qui concernent les groupes !" +
			"> `kadok groupe liste`\n" +
			"> `kadok groupe rejoindre <groupId|groupName>`\n" +
			"> `kadok groupe quitter <groupId|groupName>`",
		map[string]*Action{
			"LISTE":     &GroupListAction,
			"REJOINDRE": &GroupJoinAction,
			"QUITTER":   &GroupLeaveAction,
		},
		EmptyData,
		"404.tmpl",
	}

	GroupListAction = Action{
		security.EmptyPermission,
		"La liste des groupes disponibles !",
		map[string]*Action{},
		func(bc *BotContext, s *discordgo.Session, m *discordgo.MessageCreate, parameters []string) map[string]interface{} {
			groups := Configuration.Security.RolesHierarchy.GetGroups()
			return map[string]interface{}{
				"Groups": groups,
			}
		},
		"groupList.tmpl",
	}

	GroupJoinAction = Action{
		security.EmptyPermission,
		"C'est ici qu'on rejoint un groupe !\n> Commande : `kadok groupe rejoindre <groupId|groupName>`",
		map[string]*Action{},
		func(bc *BotContext, s *discordgo.Session, m *discordgo.MessageCreate, parameters []string) map[string]interface{} {
			roleReference := strings.Join(parameters, " ")
			addRole := MakeAddRole(s, m)
			group, err := Configuration.Security.RolesHierarchy.JoinGroup(addRole, roleReference)
			return map[string]interface{}{
				"Username": m.Author.Username,
				"Group":    group,
				"Error":    err,
			}
		},
		"groupJoin.tmpl",
	}

	GroupLeaveAction = Action{
		security.EmptyPermission,
		"C'est ici qu'on part d'un groupe ! Snif..\n> Commande : `kadok groupe quitter <groupId|groupName>`",
		map[string]*Action{},
		func(bc *BotContext, s *discordgo.Session, m *discordgo.MessageCreate, parameters []string) map[string]interface{} {
			roleReference := strings.Join(parameters, " ")
			removeRole := MakeRemoveRole(s, m)
			group, err := Configuration.Security.RolesHierarchy.LeaveGroup(removeRole, roleReference)
			return map[string]interface{}{
				"Username": m.Author.Username,
				"Group":    group,
				"Error":    err,
			}
		},
		"groupLeave.tmpl",
	}

	ClanRootAction = Action{
		security.EmptyPermission,
		"Toutes les actions qui concernent les clans !\n" +
			"> `kadok clan liste`\n" +
			"> `kadok clan rejoindre <clanId|clanName>`\n" +
			"> `kadok clan quitter <clanId|clanName>`",
		map[string]*Action{
			"LISTE":     &ClanListAction,
			"REJOINDRE": &ClanJoinAction,
			"QUITTER":   &ClanLeaveAction,
		},
		EmptyData,
		"404.tmpl",
	}

	ClanListAction = Action{
		security.EmptyPermission,
		"La liste des clans disponibles !",
		map[string]*Action{},
		func(bc *BotContext, s *discordgo.Session, m *discordgo.MessageCreate, parameters []string) map[string]interface{} {
			clans := Configuration.Security.RolesHierarchy.GetClans()
			return map[string]interface{}{
				"Clans": clans,
			}
		},
		"clanList.tmpl",
	}

	ClanJoinAction = Action{
		security.EmptyPermission,
		"C'est ici qu'on rejoint un clan !\n> Commande : `kadok clan rejoindre <clanId|clanName>`",
		map[string]*Action{},
		func(bc *BotContext, s *discordgo.Session, m *discordgo.MessageCreate, parameters []string) map[string]interface{} {
			roleReference := strings.Join(parameters, " ")
			addRole := MakeAddRole(s, m)
			removeRole := MakeRemoveRole(s, m)
			clan, err := Configuration.Security.RolesHierarchy.JoinClan(addRole, removeRole, roleReference)
			return map[string]interface{}{
				"Username": m.Author.Username,
				"Clan":     clan,
				"Error":    err,
			}
		},
		"clanJoin.tmpl",
	}

	ClanLeaveAction = Action{
		security.EmptyPermission,
		"C'est ici qu'on part d'un clan ! Snif..\n> Commande : `kadok clan quitter <clanId|clanName>`",
		map[string]*Action{},
		func(bc *BotContext, s *discordgo.Session, m *discordgo.MessageCreate, parameters []string) map[string]interface{} {
			roleReference := strings.Join(parameters, " ")
			removeRole := MakeRemoveRole(s, m)
			clan, err := Configuration.Security.RolesHierarchy.LeaveClan(removeRole, roleReference)
			return map[string]interface{}{
				"Username": m.Author.Username,
				"Clan":     clan,
				"Error":    err,
			}
		},
		"clanLeave.tmpl",
	}

	AudioRootAction = Action{
		security.EmptyPermission,
		"Tu peux gerer tous les trucs audio que je fais !\n" +
			"> `kadok audio stop`\n",
		map[string]*Action{
			"STOP": &AudioStopAction,
		},
		EmptyData,
		"404.tmpl",
	}

	AudioStopAction = Action{
		security.EmptyPermission,
		"Je m'arrete de parler 'snif snif' !\n",
		map[string]*Action{},
		func(bc *BotContext, s *discordgo.Session, m *discordgo.MessageCreate, parameters []string) map[string]interface{} {
			err := bc.Voice.Stop()
			return map[string]interface{}{
				"Username": m.Author.Username,
				"Error":    err,
			}
		},
		"audioStop.tmpl",
	}

	RadioRootAction = Action{
		security.EmptyPermission,
		"Toutes les radios que tu peux écouter !\n" +
			"> `kadok radio liste <pageIndex>`\n" +
			"> `kadok radio ecoute <pageId>`\n" +
			"> `kadok radio info <radioId>`\n",
		map[string]*Action{
			"LISTE":  &RadioListAction,
			"INFO":   &RadioInfoAction,
			"ECOUTE": &RadioListenAction,
		},
		EmptyData,
		"404.tmpl",
	}

	RadioListAction = Action{
		security.EmptyPermission,
		"La liste des radios disponibles !",
		map[string]*Action{},
		func(bc *BotContext, s *discordgo.Session, m *discordgo.MessageCreate, parameters []string) map[string]interface{} {
			stations, err := Configuration.Radio.France.GetStations(bc.HttpClient)
			if err != nil {
				return map[string]interface{}{
					"Error": err,
				}
			}
			requestedPage := strings.Join(parameters, " ")
			nbPage, currentPage, indexStart, indexEnd := utils.CalculatePagination(len(stations), 10, requestedPage)
			return map[string]interface{}{
				"Username":    m.Author.Username,
				"NbPage":      nbPage,
				"CurrentPage": currentPage,
				"Stations":    stations[indexStart:indexEnd],
				"Error":       err,
			}
		},
		"radioList.tmpl",
	}

	RadioInfoAction = Action{
		security.EmptyPermission,
		"Info d'une radio en particulier !",
		map[string]*Action{},
		func(bc *BotContext, s *discordgo.Session, m *discordgo.MessageCreate, parameters []string) map[string]interface{} {
			stationId := strings.Join(parameters, " ")
			station, err := Configuration.Radio.France.GetStation(bc.HttpClient, stationId)
			return map[string]interface{}{
				"Username": m.Author.Username,
				"Station":  station,
				"Error":    err,
			}
		},
		"radioInfo.tmpl",
	}

	RadioListenAction = Action{
		security.EmptyPermission,
		"Ecoute la radio !",
		map[string]*Action{},
		func(bc *BotContext, s *discordgo.Session, m *discordgo.MessageCreate, parameters []string) map[string]interface{} {
			returnError := func(err error) map[string]interface{} {
				return map[string]interface{}{
					"Username": m.Author.Username,
					"Error":    err,
				}
			}
			stationId := strings.Join(parameters, " ")
			station, err := Configuration.Radio.France.GetStation(bc.HttpClient, stationId)
			voiceChannelID, err := GetUserCurrentVoiceChannel(s, m.GuildID, m.Author.ID)
			if err != nil {
				return returnError(err)
			}

			response, err := bc.HttpClient.OpenStream(http.Request{
				Method: "GET",
				Url:    *station.StreamUrl,
			})
			if err != nil {
				return returnError(err)
			}

			err = bc.Voice.Play(response.Stream, m.GuildID, voiceChannelID)
			if err != nil {
				return returnError(err)
			}

			return map[string]interface{}{
				"Username": m.Author.Username,
				"Station":  station,
			}
		},
		"radioListen.tmpl",
	}

	// RootAction is the first action call by Kadok for resolve
	RootAction = Action{
		security.GetHelp,
		"Tatan elle fait du flan, elle m'a aussi dit de dire des choses intelligentes si on m'appel: \n" +
			"> - `kadok aide` je te dis ce que je fais ! Et `kadok <commande> aide` je te donne plus de details !\n" +
			"> - `kadok aqui` ? Je dis tous mes amis !\n" +
			"> - `kadok tatan` je te parle de moi !\n" +
			"> - `kadok groupe <liste|rejoindre|quitter>` Pour voir et rejoindre un groupe ! Tu peux etre dans autant de groupes que tu veux !\n" +
			"> - `kadok clan <liste|rejoindre|quitter>` Pour voir et rejoindre un clan ! Tu peux avoir seulement un clan !\n" +
			"> - `kadok radio <liste|info|ecoute>` Pour voir toutes les radios à écouter !\n" +
			"> - `kadok audio <stop>` Pour que j'arrete de parler !\n",
		map[string]*Action{
			"AQUI":   &GetCharactersAction,
			"TATAN":  &StatusAction,
			"GROUPE": &GroupRootAction,
			"CLAN":   &ClanRootAction,
			"RADIO":  &RadioRootAction,
			"AUDIO":  &AudioRootAction,
		},
		EmptyData,
		"404.tmpl",
	}
)
