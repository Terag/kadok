package radio

import (
	"net/url"
)

type StationType int

// Type of radio stations
const (
	//LiveStation - classic radio station
	LiveStation StationType = iota
	//CallCharacter - local or regional stations
	LocalStation
	//GetCharacterList - web station
	WebStation
)

type Station struct {
	Id          string
	Name        string
	Type        StationType
	Summary     string
	Description string
	StreamUrl   *url.URL
}
