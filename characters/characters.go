// The characters package contains all the features related to characters and quoting them.
//
// As of today the package allows the loading of Characters structures that are a name and a list of associated quote.
// It also implement basic functionalities on top of the Character structure.
package characters

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

// Properties structure for the package characters
// Yaml structure to add to the properties' hierarchy :
//
//  folder: <characters_folder_path> (folder path where the characters information files are stored)
type Properties struct {
	Folder string
	List   []Character
}

// UnmarshalYAML implementation for the package Properties
func (properties *Properties) UnmarshalYAML(unmarshal func(interface{}) error) error {

	type PropertiesYAML struct {
		Folder string `yaml:"folder"`
	}

	var propertiesYAML PropertiesYAML
	err := unmarshal(&propertiesYAML)
	if err != nil {
		return err
	}

	characters, err := MakeCharactersSliceFromFolder(propertiesYAML.Folder)
	if err != nil {
		return err
	}
	properties.Folder = propertiesYAML.Folder
	properties.List = characters

	return nil
}

// Character to retrieve quotes from. A Character is a name with a list of sentences.
type Character struct {
	Name      string   `json:"name"`
	Sentences []string `json:"sentences"`
}

// MakeCharactersSliceFromFolder return a slice of characters generated from a folder containing the characters files.
func MakeCharactersSliceFromFolder(folder string) ([]Character, error) {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		fmt.Println("Error loading characters folder")
		return make([]Character, 0), err
	}

	characters := make([]Character, 0, len(files))
	for _, file := range files {
		filename := filepath.Join(folder, file.Name())

		jsonFile, err := os.Open(filename)
		if err != nil {
			fmt.Println("Error getting " + filename + "'s sentences and name")
		} else {
			byteValue, _ := ioutil.ReadAll(jsonFile)
			var character Character
			err = json.Unmarshal(byteValue, &character)
			if err != nil {
				fmt.Println("Error loading character file " + filename)
				log.Fatal(err)
			}
			characters = append(characters, character)
			fmt.Println(character.Name + " loaded succesfully!")

			err = jsonFile.Close()
			if err != nil {
				fmt.Println(character.Name + " error closing " + filename)
			}
		}
	}
	return characters, nil
}

// GetQuoteFromMessage returns a random quote from a character if its name is found in the message
func GetQuoteFromMessage(characters []Character, message string) (string, error) {
	for _, character := range characters {
		upperCaseMessage := strings.Split(strings.ToUpper(message), " ")
		upperCaseCharacter := strings.ToUpper(character.Name)

		for _, word := range upperCaseMessage {
			if upperCaseCharacter == word {
				if len(character.Sentences) < 1 {
					return "", errors.New("Error, empty character quotes for " + character.Name)
				}

				index := rand.Intn(len(character.Sentences))

				return character.Sentences[index], nil
			}
		}
	}
	return "", errors.New("no quote found")
}
