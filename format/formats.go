package format

import "fmt"

const (
	DateShort6          = "060102"
	DateShort6Slashes   = "06/01/02"
	DateShort8          = "20060102"
	DateShortSlashes    = "2006/01/02"
	DateShortDashes     = "2006-01-02"
	DateTimeShort       = "02/01/2006 15:04"
	DateTimeShortDashes = "2006-01-02 15:04:05"
)

func CentToCommaRand(cent int) string {
	r := cent / 100
	c := cent % 100
	return fmt.Sprintf("%d,%02d", r, c)
}
