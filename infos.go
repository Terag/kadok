package main

import "runtime"

// Variables set at compilation time. Used to provide general information about the bot
var (
	// Version of the bot. It correspond to the associated git tag
	Version = "0.1.0"
	// UTC date of the build run
	BuildDate string
	// Git commit reference of the build
	GitCommit string
	// License name is static
	LicenseName = "GNU General Public License v3.0"
	// License url is static
	LicenseURL = "https://www.gnu.org/licenses/gpl-3.0-standalone.html"
	// Short description
	About = "Kadok is a Discord bot firstly developed for the Guild \"Les petits pedestres\". It aims to provide fun and useful functionalities for the Guild Members."
	// Url of the documentation
	URL = "https://kadok.pedestres.fr"
	// git contributors
	Contributors = "kadok_team"
)

type Infos struct {
	About        string
	LicenseName  string
	LicenseURL   string
	Version      string
	GitCommit    string
	BuildDate    string
	GoVersion    string
	URL          string
	Contributors string
}

func GetInfos() Infos {
	return Infos{
		About,
		LicenseName,
		LicenseURL,
		Version,
		GitCommit,
		BuildDate,
		runtime.Version(),
		URL,
		Contributors,
	}
}
