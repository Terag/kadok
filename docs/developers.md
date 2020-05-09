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

In Kadok, you have two ways of adding an action to the current action tree:

- The most customizable is by implementing the interface **ActionResolver** from the file actions.go.
- The out of the box one is by creating a new **Action** global variable

```go
package main

type ExecuteAction func(s *discordgo.Session, m *discordgo.MessageCreate) error

type ActionResolver interface {
	// Required by the security package.
    GetPermission() security.Permission
    // Resolve the call and return the ActionResolver with a function to execute it.
    ResolveAction(call []string) (ActionResolver, ExecuteAction)
}

type Action struct {
    // The permission required to use this action. Use the EmptyPermission if no permission is required
    Permission  security.Permission
    // The sub actions in the hierarchy. They do not inherit the permission requirement
    SubActions  map[string]ActionResolver
    // The function to execute when the action is called
    Execute     func(s *discordgo.Session, m *discordgo.MessageCreate, parameters []string) error
    // The function to display information about the action
    DisplayHelp func(s *discordgo.Session, m *discordgo.MessageCreate) error
}
```

In this example we will use the out of the box option.
For that you simply have to add a new Action to the list of available Actions as a global variable.

Don't forget to add a permission in security/rbac.go if you need a new permission. It will directly by available for use in the configuration files.

You then need to reference it in the parent action for it to be added to the hierarchy (it is true for both methods).
The key in the map should be UPPERCASE as the bot is case-insensitive.

The RootAction is the root node of the tree. It is called first by Kadok to resolve the action.

Create the Action:

```go
package main

var(
	PingAction = Action{
		security.GetCharacterList,
		map[string]ActionResolver{},
		func(s *discordgo.Session, m *discordgo.MessageCreate, parameters []string) error {
			_, err := s.ChannelMessageSend(m.ChannelID, "À Kadoc ! À Kadoc ! Pong!")
			return err
		},
		func(s *discordgo.Session, m *discordgo.MessageCreate) error {
			_, err := s.ChannelMessageSend(m.ChannelID, "Kadok répond pong!")
			return err
		},
	}

	RootAction = Action{
		security.GetHelp,
		map[string]ActionResolver{
			"PONG": &PongAction,
		},
		func(s *discordgo.Session, m *discordgo.MessageCreate, parameters []string) error {
			return errors.New("Not a valid action")
		},
		func(s *discordgo.Session, m *discordgo.MessageCreate) error {
			message := ""
			message += "\nTatan elle fait du flan, elle m'a aussi dit de dire des choses intelligentes si on m'appel: 'AKadok'"
			message += "\n'Kadok aqui' ? Je dis tous mes amis !"
			_, err := s.ChannelMessageSend(m.ChannelID, message)
			return err
		},
	}
)
```