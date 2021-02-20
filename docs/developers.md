# Developers Documentation

The purpose of this page is to help you understand how Kadok is organized.

## Code Organization

The main package holds the global structure of the bot and is also the only package importing and using the SDK DiscordGo.
Other packages are responsible to deliver domain specific features such as security features or characters features.

When you find yourself developing features for a new category, please create the associated package containing
the functions associated to your category.

If the category needs a configuration or load external files, you should take example on other packages and create
a custom Properties structure with its own UnmarshalYAML implementation. You then can add it to the Properties structure
in the main package.

## How to add a Feature

In this example, we assume that the functions required by your feature already exist in the right package.
The objective here is to add it to the current action tree and make it callable by the users.

The action tree is the hierarchy of possible actions. Most of the actions have the option to implement a help feature that can
called like this `kadok <you-command> help` (example: `kadok minecraft help`, it can also be called on the root command `kadok help` for general information about the bot).
Example of commands hierarchy:

```
 kadok
 ├─ aqui
 ├─ ping
 ├─ pong
 ├─ minecraft
 |    ├─ status
 |    ├─ whitelist <username>
 |    └─ blacklist <username>
 └─ dice <dices-pattern>
```

In Kadok, you can add an action to the current action tree by creating an Action structure.

```go
package main

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
```

For that you simply have to add a new Action to the list of available Actions as a global variable.

Don't forget to add a permission in security/rbac.go if you need a new permission. It will directly by available for use in the configuration files.

You then need to reference it in the parent action for it to be added to the hierarchy.
The key in the map should be UPPERCASE as the bot is case-insensitive.

The RootAction is the root node of the tree. It is called first by Kadok to resolve the action.

Examples:
```go
package main

var(
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
```
