package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/hajimehoshi/go-mp3"
	"gopkg.in/hraban/opus.v2"
	"gopkg.in/yaml.v2"
)

// Properties structure. Default from properties.yaml
type Properties struct {
	Characters struct {
		Folder string `yaml:"folder"`
		List   []struct {
			Name string `yaml:"name"`
			File string `yaml:"file"`
			Data Character
		}
	}
	Audio struct {
		Folder string `yaml:"folder"`
	}
}

// Variables used globally
var (
	Token         string
	Configuration Properties
	Buffer        [][]byte
)

const sampleRate = 48000
const channels = 2 // 1 mono; 2 stereo

var mutex = &sync.Mutex{}

// Character of Kaamelott to retrieve sentences from
type Character struct {
	Sentences []string `json:"sentences"`
}

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
		panic(err)
	}

	for index := range Configuration.Characters.List {
		jsonFile, err := os.Open(Configuration.Characters.Folder + "/" + Configuration.Characters.List[index].File)

		if err != nil {
			fmt.Println("Error getting " + Configuration.Characters.List[index].Name + "'s sentences")
		} else {
			byteValue, _ := ioutil.ReadAll(jsonFile)
			json.Unmarshal(byteValue, &Configuration.Characters.List[index].Data)
			fmt.Println(Configuration.Characters.List[index].Name + " loaded succesfully!")
		}

		jsonFile.Close()
	}
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

	// If the message is "ping" reply with "Pong!"
	if strings.ToUpper(m.Content) == "PING" {
		s.ChannelMessageSend(m.ChannelID, "À Kadoc ! À Kadoc ! Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if strings.ToUpper(m.Content) == "PONG" {
		s.ChannelMessageSend(m.ChannelID, "À Kadoc ! À Kadoc ! Ping!")
	}

	if strings.ToUpper(m.Content) == "KADOK HELP" {
		message := ""
		message += "\nTatan elle fait du flan, elle m'a aussi dit de dire des choses intelligentes si on m'appel: 'AKadok'"
		message += "\nLe caca des pigeons, c'est caca. J'ai pleins d'amis à Kaamelott:"
		for _, character := range Configuration.Characters.List {
			message += "\n- " + character.Name
		}
		s.ChannelMessageSend(m.ChannelID, message)
		return
	}

	if strings.ToUpper(m.Content) == "AKADOK" {
		mutex.Lock()
		defer mutex.Unlock()
		audioFiles, err := ioutil.ReadDir(Configuration.Audio.Folder)
		index := rand.Intn(len(audioFiles))
		Buffer := make([][]byte, 0)
		go loadSound(Configuration.Audio.Folder + "/" + audioFiles[index].Name())

		// Find the channel that the message came from.
		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			// Could not find channel.
			return
		}

		// Find the guild for that channel.
		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			// Could not find guild.
			return
		}

		// Look for the message sender in that guild's current voice states.
		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				err = playSound(s, g.ID, vs.ChannelID)
				if err != nil {
					fmt.Println("Error playing sound:", err)
				}

				return
			}
		}
		return
	}

	for _, character := range Configuration.Characters.List {
		if strings.Index(strings.ToUpper(m.Content), strings.ToUpper(character.Name)) > -1 {
			if len(character.Data.Sentences) < 1 {
				s.ChannelMessageSend(m.ChannelID, "Mordu, mordu mordu moooooooooooooooordu mordu mordu mordu mordu mordu mordu mordu mordu mordu mordu mordu morduuuuuuuuuuuuuuuuuuuuuuuuuuuuu!!!!")
				return
			}

			index := rand.Intn(len(character.Data.Sentences))

			s.ChannelMessageSend(m.ChannelID, character.Data.Sentences[index])
			return
		}
	}
}

// loadSound attempts to load an encoded sound file from disk.
func loadSound(path string) {

	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	decoder, err := mp3.NewDecoder(file)
	if err != nil {
		return
	}

	opusEnc, err := opus.NewEncoder(sampleRate, channels, opus.AppVoIP)
	if err != nil {
		fmt.Println("Error creating opus encoder :", err)
		return
	}

	var writer io.WriteCloser

	var inBuff = make([]byte, 1024)

	for {
		// Read opus frame length from mp3 file.
		_, err := decoder.Read(inBuff)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return
		}

		if err != nil {
			fmt.Println("Error reading file :", err)
			return
		}

		n, err := enc.Encode(pcm, inBuff)
		if err != nil {
			fmt.Println("Error encoding Opus : ", err)
			return
		}
		inBuff = inBuff[:n]

		// Append decoded mp3 data to the buffer channel.
		Buffer = append(Buffer, inBuff)
	}
}

// playSound plays the current buffer to the provided channel.
func playSound(s *discordgo.Session, guildID, channelID string) (err error) {

	//Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(250 * time.Millisecond)

	// Start speaking.
	vc.Speaking(true)

	// Send the buffer data.
	for row := range Buffer {
		vc.OpusSend <- row
	}

	// Stop speaking
	vc.Speaking(false)

	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)

	// Disconnect from the provided voice channel.
	vc.Disconnect()

	return nil
}
