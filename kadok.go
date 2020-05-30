// Kadok is a bot developed for the Discord Guild "Les petits pedestres". It aims to provide fun and useful functionnalities for the Guilde Members.
//
// The code of Kadok is open source, you can use it and modify it as long as you follow the GNU_v3 licence terms.
//
// The current available features are:
//
// 1. Play ping pong
//
// 2. Quote characters from Kaamelott universe (and even more!)
//
// 3. Manage features' permission based on Discord roles
//
// This documentation is mainly to help developers to understand how the inside of Kadok works
// To find more on how to start and configure the bot, please visit the wiki page: https://github.com/terag/kadok/wiki
//
// Main package holds the global structure of bot and is also the only package importing the SDK DiscordGo.
// Other packages are responsible to deliver domain specific features such as security features or characters features.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Terag/kadok/characters"
	"github.com/Terag/kadok/security"
	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v2"
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
	Characters characters.Properties `yaml:"characters"`
	Security   security.Properties   `yaml:"security"`
}

// Variables used globally
var (
	// Bot token to access Discord API
	Token string
	// Instance of properties structure
	Configuration Properties
)

// Variables set at compilation time. Used to provide general information about the bot
var (
	// Version of the bot. It correspond to the associated git tag
	Version = "0.1.0"
	// UTC date of the build run
	BuildDate string
	// Git commit reference of the build
	GitCommit string
	// Version of go used to build Kadok
	GoVersion string
	// License name is static
	LicenseName = "GNU General Public License v3.0"
	// License url is static
	LicenseUrl = "https://www.gnu.org/licenses/gpl-3.0-standalone.html"
	// Short description
	About = "Kadok is a Discord bot firstly developed for the Guild \"Les petits pedestres\". It aims to provide fun and useful functionalities for the Guild Members."
)

// Called before main to initialize global variables and configuration properties
// It gets the flags and unmarshal the properties.yaml file, populating the global configuration.
func init() {
	rand.Seed(time.Now().UnixNano())

	var configPath string
	helpFlag := flag.Bool("h", false, "Print this message and exit.")
	versionFlag := flag.Bool("v", false, "Print Kadok version and build information")
	infoFlag := flag.Bool("i", false, "Print Kadok information and credits")
	flag.StringVar(&Token, "t", "", "Bot's token to connect to discord")
	flag.StringVar(&configPath, "p", "properties.yaml", "Properties file, default 'properties.yaml'")
	flag.Parse()

	switch {
	case *helpFlag:
		flag.Usage()
		os.Exit(0)
	case *versionFlag:
		fmt.Println("Version: " + Version)
		if GitCommit != "" {
			fmt.Println("Build commit: " + GitCommit)
		}
		if BuildDate != "" {
			fmt.Println("Build date: " + BuildDate)
		}
		if GoVersion != "" {
			fmt.Print("Go: " + GoVersion)
		}
		os.Exit(0)
	case *infoFlag:
		fmt.Println(About)
		fmt.Println("Licensed under " + LicenseName)
		fmt.Println("Full license: " + LicenseUrl)
		os.Exit(0)
	}

	if Token == "" {
		fmt.Println("Use -t to specify the discord's bot token")
		return
	}

	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(configFile, &Configuration)
	if err != nil {
		log.Fatal(err)
	}
}

// main function starting the bot
func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	err = dg.Close()
	if err != nil {
		fmt.Println(err)
	}
}

// This function will be called (due to AddHandler from main()) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
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
			_, err = s.ChannelMessageSend(m.ChannelID, quote)
			if err != nil {
				fmt.Println(err)
			}
		}
		return
	}
}
