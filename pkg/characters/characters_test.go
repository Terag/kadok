package characters

import (
	"testing"

	"gopkg.in/yaml.v3"
)

const charactersFolderPath = "../../assets/characters"

func TestCharactersPropertiesUnmarshalYAML(t *testing.T) {
	propertiesYAML := []byte("folder: \"" + charactersFolderPath + "\"")
	var properties Properties
	err := yaml.Unmarshal(propertiesYAML, &properties)
	if err != nil {
		t.Errorf("Error loading characters module properties: %v", err)
	}
	if len(properties.List) == 0 {
		t.Errorf("Error, no characters were loaded, got: %v want more", len(properties.List))
	}
}

func TestMakeCharactersSliceFromFolder(t *testing.T) {
	characters, err := MakeCharactersSliceFromFolder(charactersFolderPath)
	if err != nil {
		t.Errorf("Error loading characters from folder, error: %v", err)
	}
	if len(characters) == 0 {
		t.Errorf("Error, no characters were loaded, got: %v want more", len(characters))
	}
}

func TestGetQuoteFromMessage(t *testing.T) {
	characters, err := MakeCharactersSliceFromFolder(charactersFolderPath)
	_, err = GetQuoteFromMessage(characters, "Bonjour monsieur arthur")
	if err != nil {
		t.Errorf("Error finding a quote from message, error: %v", err)
	}
}
