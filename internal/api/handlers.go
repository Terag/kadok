package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// httpErrorHandler used to return errors in the format of the [RFC7807: Problem Details](https://www.rfc-editor.org/rfc/rfc7807)
func httpErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	c.Logger().Error(err)
	var problem Problem
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		problem = Problem{
			Type:   nullable("about:blank"),
			Title:  nullable(httpStatusCodeToMessage(code)),
			Detail: nullable(fmt.Sprintf("%v", he.Message)),
			Status: nullable(int32(code)),
		}
	} else {
		problem = createProblem500Default()
	}
	c.JSON(code, problem)
}
