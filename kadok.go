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
	"sync"
	"syscall"
	"time"

	"github.com/Terag/kadok/security"
	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v2"
)

// Properties structure. Default from properties.yaml
type Properties struct {
	Characters struct {
		Folder string `yaml:"folder"`
		List   []Character
	}
	Audio struct {
		Folder string `yaml:"folder"`
	}
	Security struct {
		RolesConfiguration string `yaml:"roles"`
	}
}

// Variables used globally
var (
	Token         string
	Configuration Properties
	Buffer        [][]byte
	RolesTree     security.RolesTree
)

const sampleRate = 48000
const channels = 2 // 1 mono; 2 stereo

var mutex = &sync.Mutex{}

func init() {
	rand.Seed(time.Now().UnixNano())

	var configPath string
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&configPath, "p", "properties.yaml", "Properties file")
	flag.Parse()

	loadConfiguration(configPath)
}

func loadConfiguration(path string) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Config file didn't load", r)
		}
	}()

	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(configFile, &Configuration)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	RolesTree, err = security.MakeRolesTreeFromFile(Configuration.Security.RolesConfiguration)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	loadCharacters()
}

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
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	roles, err := GetUserRoles(s, m)
	if err != nil {
		fmt.Println("Error retrieving user roles")
		return
	}
	isGranted := security.MakeIsGranted(RolesTree, roles)

	// If the message is "ping" reply with "Pong!"
	if strings.ToUpper(m.Content) == "PING" {
		s.ChannelMessageSend(m.ChannelID, "À Kadoc ! À Kadoc ! Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if strings.ToUpper(m.Content) == "PONG" {
		s.ChannelMessageSend(m.ChannelID, "À Kadoc ! À Kadoc ! Ping!")
	}

	if strings.ToUpper(m.Content) == "KADOK HELP" {
		if isGranted(security.GetHelp) {
			message := ""
			message += "\nTatan elle fait du flan, elle m'a aussi dit de dire des choses intelligentes si on m'appel: 'AKadok'"
			message += "\n'Kadok aqui' ? Je dis tous mes amis !"
			s.ChannelMessageSend(m.ChannelID, message)
		}
		return
	}

	if strings.ToUpper(m.Content) == "KADOK AQUI" {
		if isGranted(security.GetCharacterList) {
			displayAvailableCharacters(s, m)
		}
		return
	}

	if strings.ToUpper(m.Content) == "AKADOK" {
		s.ChannelMessageSend(m.ChannelID, "Mordu, mordu mordu moooooooooooooooordu mordu mordu mordu mordu mordu mordu mordu mordu mordu mordu mordu morduuuuuuuuuuuuuuuuuuuuuuuuuuuuu!!!!")
		return
	}

	if isGranted(security.CallCharacter) {
		handleCalledCharacter(s, m)
	}
}
