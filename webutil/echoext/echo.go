package echoext

import (
	"encoding/json"
	"github.com/Direct-Debit/go-commons/errlib"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"html/template"
	"io"
	"net/http"
)

func ErrorHandlerJSON(err error, c echo.Context) {
	if c.Response().Committed {
		// workaround to multiple error handler calls:
		// https://github.com/labstack/echo/issues/1948#issuecomment-1068045153
		return
	}

	var msg interface{}

	code := http.StatusInternalServerError
	msg = "An internal server error has occurred. Our engineers have been notified. Please try again later."

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
		errlib.DebugError(err, "A client error has occurred")
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

type TemplateRenderer struct {
	dev       bool
	glob      string
	templates *template.Template
}

// NewTemplateRenderer creates an echo.Renderer that uses all the template files matching the given glob.
// If dev is set to true, the renderer will reread the templates on every Render.
func NewTemplateRenderer(glob string, dev bool) echo.Renderer {
	return TemplateRenderer{
		dev:       dev,
		glob:      glob,
		templates: template.Must(template.ParseGlob(glob)),
	}
}

func (t TemplateRenderer) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	if t.dev {
		log.Debugf("Rereading templates for dev environment")
		t.templates = template.Must(template.ParseGlob(t.glob))
	}
	return t.templates.ExecuteTemplate(w, name, data)
}
