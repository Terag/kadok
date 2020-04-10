package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Character of Kaamelott to retrieve sentences from
type Character struct {
	Name      string   `json:"name"`
	Sentences []string `json:"sentences"`
}

func loadCharacters() {
	files, err := ioutil.ReadDir(Configuration.Characters.Folder)
	if err != nil {
		fmt.Println("Error loading characters folder")
		log.Fatal(err)
	}

	for _, file := range files {
		filename := filepath.Join(Configuration.Characters.Folder, file.Name())

		jsonFile, err := os.Open(filename)
		if err != nil {
			fmt.Println("Error getting " + filename + "'s sentences qnd nqme")
		} else {
			byteValue, _ := ioutil.ReadAll(jsonFile)
			var character Character
			err = json.Unmarshal(byteValue, &character)
			if err != nil {
				fmt.Println("Error loading character file " + filename)
				log.Fatal(err)
			}
			Configuration.Characters.List = append(Configuration.Characters.List, character)
			fmt.Println(character.Name + " loaded succesfully!")
		}

		jsonFile.Close()
	}
}

func displayAvailableCharacters(s *discordgo.Session, m *discordgo.MessageCreate) {
	message := ""
	message += "\nLe caca des pigeons, c'est caca. Si tu parle d'un de mes amis, je te dirai ce qu'il a dit."
	message += "\nJ'ai pleins d'amis à Kaamelott:"
	for _, character := range Configuration.Characters.List {
		message += "\n- " + character.Name
	}
	s.ChannelMessageSend(m.ChannelID, message)
}

func handleCalledCharacter(s *discordgo.Session, m *discordgo.MessageCreate) {
	for _, character := range Configuration.Characters.List {
		if strings.Index(strings.ToUpper(m.Content), strings.ToUpper(character.Name)) > -1 {
			if len(character.Sentences) < 1 {
				s.ChannelMessageSend(m.ChannelID, "Mordu, mordu mordu moooooooooooooooordu mordu mordu mordu mordu mordu mordu mordu mordu mordu mordu mordu morduuuuuuuuuuuuuuuuuuuuuuuuuuuuu!!!!")
				return
			}

			index := rand.Intn(len(character.Sentences))

			s.ChannelMessageSend(m.ChannelID, character.Sentences[index])
			return
		}
	}
}
