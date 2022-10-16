package api

import (
	"net/http"
)

// createProblem500Default return the default problem for errors 500
func createProblem500Default() Problem {
	return Problem{
		Type:   nullable("about:blank"),
		Title:  nullable("Internal Error"),
		Detail: nullable("Oups, something went wrong! Please try later"),
		Status: nullable(int32(http.StatusInternalServerError)),
	}
}
