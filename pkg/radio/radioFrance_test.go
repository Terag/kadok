package radio

import (
	"net/url"
	"testing"

	"github.com/terag/kadok/internal/http"
)

var QueryBrandsResponse string = `{
    "data": {
        "brands": [
            {
                "id": "FRANCEINTER",
                "title": "France Inter",
                "baseline": "Le direct de France Inter",
                "description": "Joyeuse, savante et populaire, France Inter est la radio généraliste de service public ",
                "liveStream": "https://icecast.radiofrance.fr/franceinter-midfi.mp3?id=openapi",
                "localRadios": null,
                "webRadios": null
            },
            {
                "id": "FRANCEINFO",
                "title": "franceinfo",
                "baseline": "Et tout est plus clair",
                "description": "L'actualité en direct et en continu avec le média global du service public",
                "liveStream": "https://icecast.radiofrance.fr/franceinfo-midfi.mp3?id=openapi",
                "localRadios": null,
                "webRadios": null
            },
            {
                "id": "FRANCEMUSIQUE",
                "title": "France Musique",
                "baseline": "Ce monde a besoin de musique",
                "description": "Classique mais pas que… Chaque jour, la musique se vit intensément sur France Musique.",
                "liveStream": "https://icecast.radiofrance.fr/francemusique-midfi.mp3?id=openapi",
                "localRadios": null,
                "webRadios": [
                    {
                        "id": "FRANCEMUSIQUE_CLASSIQUE_EASY",
                        "title": "Classique Easy",
                        "description": "Pour vous concentrer, travailler ou vous détendre, écoutez un flux musical continu diffusant les titres incontournables de la musique classique. ",
                        "liveStream": "https://icecast.radiofrance.fr/francemusiqueeasyclassique-midfi.mp3?id=openapi"
                    },
                    {
                        "id": "FRANCEMUSIQUE_CLASSIQUE_PLUS",
                        "title": "Classique Plus",
                        "description": "France Musique vous emmène plus loin dans le monde du classique avec une sélection d’œuvres moins ou peu connues à écouter en intégralité. ",
                        "liveStream": "https://icecast.radiofrance.fr/francemusiqueclassiqueplus-midfi.mp3?id=openapi"
                    },
                    {
                        "id": "FRANCEMUSIQUE_CONCERT_RF",
                        "title": "Concerts Radio France",
                        "description": "Écoutez une grande sélection d'enregistrements de concerts des orchestres et des chœurs de Radio France, avec des solistes et chefs d'orchestre de renom. ",
                        "liveStream": "https://icecast.radiofrance.fr/francemusiqueconcertsradiofrance-midfi.mp3?id=openapi"
                    },
                    {
                        "id": "FRANCEMUSIQUE_OCORA_MONDE",
                        "title": "Ocora Musiques du Monde",
                        "description": "Voyagez à travers la musique traditionnelle du monde entier avec les innombrables enregistrements du label Ocora disponibles à l’écoute en ligne. ",
                        "liveStream": "https://icecast.radiofrance.fr/francemusiqueocoramonde-midfi.mp3?id=openapi"
                    },
                    {
                        "id": "FRANCEMUSIQUE_LA_JAZZ",
                        "title": "La Jazz",
                        "description": "Écoutez des enregistrements acoustiques rares et des reprises électriques des standards du jazz depuis ses débuts jusqu’au jazz contemporain d’aujourd’hui. ",
                        "liveStream": "https://icecast.radiofrance.fr/francemusiquelajazz-midfi.mp3?id=openapi"
                    },
                    {
                        "id": "FRANCEMUSIQUE_LA_CONTEMPORAINE",
                        "title": "La Contemporaine",
                        "description": "Écoutez une musique pleine de liberté et d'innovation, des œuvres expérimentales et classiques de compositeurs vivants de renommée internationale. ",
                        "liveStream": "https://icecast.radiofrance.fr/francemusiquelacontemporaine-midfi.mp3?id=openapi"
                    },
                    {
                        "id": "FRANCEMUSIQUE_EVENEMENTIELLE",
                        "title": "La B.O. Musiques de Films",
                        "description": "Écoutez les plus beaux extraits de bandes originales, une sélection de musiques de films célèbres, mais aussi d’œuvres rares issues de la plus grande discothèque d’Europe. ",
                        "liveStream": "https://icecast.radiofrance.fr/francemusiquelabo-midfi.mp3?id=openapi"
                    }
                ]
            },
            {
                "id": "FRANCECULTURE",
                "title": "France Culture",
                "baseline": "France Culture, l'esprit d'ouverture",
                "description": "",
                "liveStream": "https://icecast.radiofrance.fr/franceculture-lofi.mp3?id=openapi",
                "localRadios": null,
                "webRadios": null
            }
        ]
    }
}`

type FakeHttpClient struct {
	QueryResponse []byte
}

func (fhc FakeHttpClient) Execute(request http.Request) (http.Response, error) {
	return http.Response{
		StatusCode: 200,
		CacheHit:   false,
		Body:       fhc.QueryResponse,
	}, nil
}

func (fhc FakeHttpClient) OpenStream(request http.Request) (http.ResponseStream, error) {
	return http.ResponseStream{
		StatusCode: 200,
		CacheHit:   false,
		Stream:     nil,
	}, nil
}

func TestRadioFranceGetStations(t *testing.T) {
	radio := RadioFrance{
		Url: func() url.URL {
			radioUrl, _ := url.ParseRequestURI("https://openapi.radiofrance.fr/v1/graphql")
			return *radioUrl
		}(),
	}
	stations, err := radio.GetStations(FakeHttpClient{
		QueryResponse: []byte(QueryBrandsResponse),
	})
	if err != nil {
		t.Errorf("Failed get stations: %v", err.Error())
		return
	}
	if len(stations) < 1 {
		t.Errorf("Error no stations in stations list")
	}
}
