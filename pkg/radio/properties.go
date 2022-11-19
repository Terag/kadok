package radio

import (
	"net/url"
)

type Properties struct {
	France *RadioFrance
}

func (p *Properties) GetStations() ([]Station, error) {
	if p.France != nil {
		return p.France.GetStations()
	} else {
		return []Station{}, nil
	}
}

func (p *Properties) GetStation(id string) (Station, error) {
	if p.France != nil {
		return p.France.GetStation(id)
	} else {
		return Station{}, nil
	}
}

// UnmarshalYAML implementation for the package Properties
func (properties *Properties) UnmarshalYAML(unmarshal func(interface{}) error) error {

	type PropertiesYAML struct {
		France struct {
			Enable bool `yaml:"enable"`
			Api    struct {
				Url   string `yaml:"url"`
				Token string `yaml:"token"`
			} `yaml:"api"`
		} `yaml:"france"`
	}

	var propertiesYAML PropertiesYAML
	err := unmarshal(&propertiesYAML)
	if err != nil {
		return err
	}

	if propertiesYAML.France.Enable {
		url, err := url.Parse(propertiesYAML.France.Api.Url)
		if err != nil {
			return err
		}

		properties.France = &RadioFrance{
			Url:    *url,
			ApiKey: propertiesYAML.France.Api.Token,
		}
	}

	return nil
}
