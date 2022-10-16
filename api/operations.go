package api

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Terag/kadok/info"
	"github.com/labstack/echo/v4"
)

// Get Well-Known
// (GET /.well-known)
func (kas KadokApiServer) GetWellKnown(ctx echo.Context) error {
	infoUrl := url.URL{
		Host:   kas.Hostname,
		Path:   kas.BasePath,
		Scheme: ctx.Scheme(),
	}
	infoUrl = *infoUrl.JoinPath("/.well-known/info")
	infoUrlString := infoUrl.String()
	return ctx.JSON(http.StatusOK, GetWellKnownResponse200{
		Info: nullable(infoUrlString),
	})
}

// Get Well-Known Info
// (GET /.well-known/info)
func (kas KadokApiServer) GetWellKnownInfo(ctx echo.Context) error {
	i := info.GetInfo()
	infoResponse := GetWellKnownInfoResponse200{
		GitCommit:   nullable(i.GitCommit),
		GoVersion:   nullable(i.GoVersion),
		LicenseName: nullable(i.LicenseName),
		LicenseURL:  nullable(i.LicenseURL),
		Url:         nullable(i.URL),
		Version:     nullable(i.Version),
	}
	//Build Date
	bDate, err := time.Parse(time.RFC3339, i.BuildDate)
	if err == nil {
		infoResponse.BuildDate = nullable(bDate)
	}
	// Contributors
	cs := strings.Split(i.Contributors, ",")
	cis := make([]InfoContributor, len(cs))
	for i, c := range cs {
		c = strings.TrimSpace(c)
		cis[i] = InfoContributor{
			Username: nullable(c),
		}
	}
	infoResponse.Contributors = &cis
	return ctx.JSON(http.StatusOK, infoResponse)
}
