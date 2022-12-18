// Bot package holds the global functionalities of the bot and is also the only package importing the SDK DiscordGo.
// Other packages are responsible to deliver domain specific features such as security features or characters features.

package bot

import (
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/terag/kadok/internal/cache"
	"github.com/terag/kadok/internal/http"
	"github.com/terag/kadok/pkg/characters"
	"github.com/terag/kadok/pkg/radio"
	"github.com/terag/kadok/pkg/security"
	"gopkg.in/yaml.v3"
)

// Properties is used to easily load packages' properties from a single yaml properties file
//
// In order to simplify the initialization of Kadok and left the responsibility of the configuration format to the relevant packages.
// Each package CAN implement a Properties structure and add it to the main package structure.
// The configuration structure will automatically be taken into account into the bot properties file.
//
// Properties is a struct in which each field is a Properties structure of a sub-package.
// The Properties structure of a package MUST implement the Unmarshaler interface https://godoc.org/gopkg.in/yaml.v2#Unmarshaler
type Properties struct {
	//Prefix value used to call the bot
	Prefix     string                `yaml:"prefix"`
	Guild      Guild                 `yaml:"guild"`
	Characters characters.Properties `yaml:"characters"`
	Security   security.Properties   `yaml:"security"`
	Radio      radio.Properties      `yaml:"radio"`
	Templates  string                `yaml:"templates"`
}

type Guild struct {
	Name string `yaml:"name"`
	ID   string
}

var (
	Configuration Properties
	Context       BotContext
)

// Run starts the bot. and call for registering the handlers.
// In case of error when launching the server, `onError` is called
// Once the server is launched, call onReady and wait on its return to close the server.
func Run(onReady func(), onError func(error), token string, configPath string) {

	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		onError(errors.New("Kadok bot error: " + err.Error()))
		return
	}

	err = yaml.Unmarshal(configFile, &Configuration)
	if err != nil {
		onError(errors.New("Kadok bot error: " + err.Error()))
		return
	}

	numberRoles := len(Configuration.Security.RolesHierarchy.Buffer)
	numberClans := len(Configuration.Security.RolesHierarchy.GetClans())
	numberGroups := len(Configuration.Security.RolesHierarchy.GetGroups())
	fmt.Println("Roles successfully loaded: " + fmt.Sprint(numberRoles) + " (Clans: " + fmt.Sprint(numberClans) + " - Groups: " + fmt.Sprint(numberGroups) + ")")

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		onError(errors.New("Kadok bot error: " + err.Error()))
		return
	}

	memCache := cache.NewMemoryCache(time.Duration(60 * time.Second))
	bc := NewBotContext(dg, memCache, http.NewHttpClient(memCache, time.Duration(30*time.Second)))

	// Register the messageCreate func as a callback for MessageCreate events.
	err = registerHandlers(dg, &bc)
	if err != nil {
		fmt.Println("error registering handler: ", err)
	}

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		onError(errors.New("Kadok bot error: " + err.Error()))
		return
	}

	// Cleanly close down the Discord session.
	defer func() {
		err = dg.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	// Check that bot is in the right guild
	guilds := dg.State.Guilds
	for _, guild := range guilds {
		guild, err = dg.Guild(guild.ID)
		if err != nil {
			onError(errors.New("Kadok bot error: " + err.Error()))
			return
		}
		if guild.Name == Configuration.Guild.Name {
			Configuration.Guild.ID = guild.ID
			break
		}
	}
	if Configuration.Guild.ID == "" {
		application, err := dg.Application("@me")
		if err != nil {
			onError(errors.New("Kadok bot error: " + err.Error()))
			return
		}
		fmt.Println("Configured guild \"", Configuration.Guild.Name, "\" not found. The bot must be invited in the right guild using the following link: https://discord.com/api/oauth2/authorize?client_id="+application.ID+"&scope=bot&permissions=8")
		onError(errors.New("Kadok bot error: Configured guild not found. The bot must be invited in the right guild using the following link: https://discord.com/api/oauth2/authorize?client_id=" + application.ID + "&scope=bot&permissions=8"))
		return
	}

	fmt.Println("Bot is now running on the Guild [" + Configuration.Guild.Name + "] and can be called with the Prefix [" + Configuration.Prefix + "].")

	onReady()
}
