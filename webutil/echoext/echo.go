package echoext

import (
	"encoding/json"
	"github.com/Direct-Debit/go-commons/errlib"
	"github.com/labstack/echo/v4"
	"net/http"
)

func ErrorHandlerJSON(err error, c echo.Context) {
	var msg interface{}

	code := http.StatusInternalServerError
	msg = "An internal server error has occurred. Our engineers have been notified. Please try again later."

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
	} else {
		errlib.ErrorError(err, "An internal server error occurred")
	}

	jsonErr := c.JSON(code, map[string]interface{}{
		"error": msg,
		"code":  code,
	})
	errlib.ErrorError(jsonErr, "Echo failed to JSON-ify error response")
}

func DecodeJSONBody(c echo.Context) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&m)
	return m, err
}
