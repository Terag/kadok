package radio

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/terag/kadok/internal/graphql"
	"github.com/terag/kadok/internal/http"
)

type RadioFrance struct {
	Client http.Client
	Url    url.URL
	ApiKey string
}

func (rf *RadioFrance) GetStations() ([]Station, error) {
	dto, err := rf.GetRadioFranceBrands()
	if err != nil {
		return []Station{}, err
	}
	return dto.GetStations(), nil
}

func (rf RadioFrance) GetName() string {
	return "Radio France"
}

func (rf *RadioFrance) GetRadioFranceBrands() (RadioFranceDto, error) {
	var rfd RadioFranceDto

	body, _ := json.Marshal(&graphql.Query{
		Query: `{
            brands {
                id
                title
                baseline
                description
                liveStream
                localRadios {
                    id
                    title
                    description
                    liveStream
                }
                webRadios {
                    id
                    title
                    description
                    liveStream
                }
            }
        }`,
	})

	response, err := rf.Client.Execute(http.Request{
		Method:   "POST",
		Url:      rf.Url,
		CacheKey: "radio_france_brands",
		Headers: []http.Header{
			{
				Key: "Content-Type",
				Values: []string{
					"application/json",
				},
			},
			{
				Key: "X-Token",
				Values: []string{
					rf.ApiKey,
				},
			},
		},
		Body: body,
	})
	if err != nil {
		return rfd, err
	}

	if response.StatusCode != 200 {
		return rfd, errors.New("Error retrieving Radio France Brands: " + string(response.Body))
	}

	err = json.Unmarshal(response.Body, &rfd)
	if err != nil {
		return rfd, err
	}

	return rfd, nil
}

type RadioFranceDto struct {
	Data struct {
		Brands []struct {
			Id          string `json:"id"`
			Title       string `json:"title"`
			Baseline    string `json:"baseline"`
			Description string `json:"description"`
			LiveStream  string `json:"liveStream"`
			LocalRadios []struct {
				Id          string `json:"id"`
				Title       string `json:"title"`
				Description string `json:"description"`
				LiveStream  string `json:"liveStream"`
			} `json:"localRadios"`
			WebRadios []struct {
				Id          string `json:"id"`
				Title       string `json:"title"`
				Description string `json:"description"`
				LiveStream  string `json:"liveStream"`
			} `json:"webRadios"`
		} `json:"brands"`
	} `json:"data"`
}

func (rfd RadioFranceDto) GetStations() []Station {
	var stations []Station
	for _, s := range rfd.Data.Brands {
		streamUrl, err := url.Parse(s.LiveStream)
		if err != nil {
			fmt.Printf("Error parse Stream Url for Radio France %v: %v\n", s.Id, err)
		} else {
			stations = append(stations, Station{
				Id:          s.Id,
				Name:        s.Title,
				Type:        LiveStation,
				Summary:     s.Baseline,
				Description: s.Description,
				StreamUrl:   *streamUrl,
			})
		}
		for _, ws := range s.WebRadios {
			if err != nil {
				fmt.Printf("Error parse Stream Url for Radio France %v: %v\n", ws.Id, err)
			} else {
				stations = append(stations, Station{
					Id:          ws.Id,
					Name:        ws.Title,
					Type:        LiveStation,
					Summary:     ws.Description,
					Description: ws.Description,
				})
			}
		}
		for _, ls := range s.LocalRadios {
			if err != nil {
				fmt.Printf("Error parse Stream Url for Radio France %v: %v\n", ls.Id, err)
			} else {
				stations = append(stations, Station{
					Id:          ls.Id,
					Name:        ls.Title,
					Type:        LiveStation,
					Summary:     ls.Description,
					Description: ls.Description,
				})
			}
		}
	}
	return stations
}
