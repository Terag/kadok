package characters

import "testing"

func TestMakeCharactersSliceFromFolder(t *testing.T) {
	characters, err := MakeCharactersSliceFromFolder("../assets/characters")
	if err != nil {
		t.Errorf("Error loading characters from folder, error: %v", err)
	}
	if len(characters) == 0 {
		t.Errorf("Error, no characters were loaded, got: %v want more", len(characters))
	}
}

func TestGetQuoteFromMessage(t *testing.T) {
	characters, err := MakeCharactersSliceFromFolder("../assets/characters")
	_, err = GetQuoteFromMessage(characters, "Bonjour monsieur arthur")
	if err != nil {
		t.Errorf("Error finding a quote from message, error: %v", err)
	}
}
