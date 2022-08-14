// Bot package holds the global functionalities of the bot and is also the only package importing the SDK DiscordGo.
// Other packages are responsible to deliver domain specific features such as security features or characters features.

package bot

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Terag/kadok/characters"
	"github.com/Terag/kadok/security"
	"github.com/bwmarrin/discordgo"
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
	Templates  string                `yaml:"templates"`
}

type Guild struct {
	Name string `yaml:"name"`
	ID   string
}

var (
	Configuration Properties
)

// Run starts the bot. and call for registering the handlers.
func Run(token string, configPath string) {

	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal("yamlFile.Get err: ", err)
	}

	err = yaml.Unmarshal(configFile, &Configuration)
	if err != nil {
		log.Fatal(err)
	}
	numberRoles := len(Configuration.Security.RolesHierarchy.Buffer)
	numberClans := len(Configuration.Security.RolesHierarchy.GetClans())
	numberGroups := len(Configuration.Security.RolesHierarchy.GetGroups())
	fmt.Println("Roles successfully loaded: " + fmt.Sprint(numberRoles) + " (Clans: " + fmt.Sprint(numberClans) + " - Groups: " + fmt.Sprint(numberGroups) + ")")

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session: ", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	err = registerHandlers(dg)
	if err != nil {
		fmt.Println("error registering handler: ", err)
	}

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection: ", err)
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
			log.Fatal(err)
		}
		if guild.Name == Configuration.Guild.Name {
			Configuration.Guild.ID = guild.ID
			break
		}
	}
	if Configuration.Guild.ID == "" {
		application, err := dg.Application("@me")
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal("Configured guild not found. The bot must be invited in the right guild using the following link: https://discord.com/api/oauth2/authorize?client_id=" + application.ID + "&scope=bot&permissions=8")
	}

	fmt.Println("Bot is now running on the Guild [" + Configuration.Guild.Name + "] and can be called with the Prefix [" + Configuration.Prefix + "].")
	fmt.Println("Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	err = dg.Close()
	if err != nil {
		fmt.Println(err)
	}
}
