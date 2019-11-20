package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables used globally
var (
	Token    string
	Kadok    Character
	Perceval Character
)

var buffer = make([][]byte, 0)

// Character of Kaamelott to retrieve sentences from
type Character struct {
	Sentences []string `json:"sentences"`
}

func init() {

	jsonFile, err := os.Open("characters/kadok.json")

	if err != nil {
		fmt.Println("Error getting Kadok's sentences")
	} else {
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &Kadok)
		fmt.Println("Kadok loaded succesfully!")
	}

	jsonFile.Close()
	jsonFile, err = os.Open("characters/perceval.json")

	if err != nil {
		fmt.Println("Error getting Perceval's sentences")
	} else {
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &Perceval)
		fmt.Println("Perceval loaded succesfully!")
	}

	jsonFile.Close()

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
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

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
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
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// If the message is "ping" reply with "Pong!"
	if strings.ToUpper(m.Content) == "PING" {
		s.ChannelMessageSend(m.ChannelID, "À Kadoc ! À Kadoc ! Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if strings.ToUpper(m.Content) == "PONG" {
		s.ChannelMessageSend(m.ChannelID, "À Kadoc ! À Kadoc ! Ping!")
	}

	if strings.Index(strings.ToUpper(m.Content), "KADOK") > -1 {
		if len(Kadok.Sentences) < 1 {
			s.ChannelMessageSend(m.ChannelID, "Mordu, mordu mordu moooooooooooooooordu mordu mordu mordu mordu mordu mordu mordu mordu mordu mordu mordu morduuuuuuuuuuuuuuuuuuuuuuuuuuuuu!!!!")
			return
		}

		index := rand.Intn(len(Kadok.Sentences))

		s.ChannelMessageSend(m.ChannelID, Kadok.Sentences[index])
	}

	if strings.Index(strings.ToUpper(m.Content), "PERCEVAL") > -1 {
		if len(Perceval.Sentences) < 1 {
			s.ChannelMessageSend(m.ChannelID, "Mordu, mordu mordu moooooooooooooooordu mordu mordu mordu mordu mordu mordu mordu mordu mordu mordu mordu morduuuuuuuuuuuuuuuuuuuuuuuuuuuuu!!!!")
			return
		}

		index := rand.Intn(len(Perceval.Sentences))

		s.ChannelMessageSend(m.ChannelID, Perceval.Sentences[index])
	}
}

func sayRandom(s *discordgo.Session) {

}
