package echoext

import (
	"github.com/Direct-Debit/go-commons/errlib"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"strconv"
)

const PageDataKey = "page_data"

type PageData struct {
	Limit  int
	Cursor string
}

func PageHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		size, err := strconv.Atoi(ctx.QueryParam("limit"))
		if errlib.DebugError(err, "Could not parse limit") {
			return echo.NewHTTPError(400, "Invalid Page Size")
		}
		cursor := ctx.QueryParam("cursor")

		ctx.Set(PageDataKey, PageData{Limit: size, Cursor: cursor})
		return next(ctx)
	}
}

func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("server", "Damage Per Second")
		return next(c)
	}
}

func JSONRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		m, err := DecodeJSONBody(c)
		if err != nil {
			log.Warnf("Failed to parse JSON for request on %s: %v", c.Path(), err)
		}
		for k, v := range m {
			c.Set(k, v)
		}

		return next(c)
	}
}

// ParseQueryArgs parses query args and puts them directly into the request context
func ParseQueryArgs(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		vls := c.QueryParams()
		for k, vs := range vls {
			if len(vs) == 1 {
				c.Set(k, vs[0])
			} else {
				c.Set(k, vs)
			}
		}

		return next(c)
	}
}
