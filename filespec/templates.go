package filespec

import (
	"fmt"
	"github.com/Direct-Debit/go-commons/errlib"
	"github.com/Direct-Debit/go-commons/fileio"
	"strings"
	"text/template"
)

func ProcessTemplate(location string, data interface{}) string {
	storage := fileio.CurrStorage()
	tmpContent, err := storage.Load(location)
	errlib.PanicError(err, fmt.Sprintf("Couldn't load %v template", location))

	tmp, err := template.New(location).Parse(tmpContent)
	errlib.PanicError(err, fmt.Sprintf("Couldn't parse %v template", location))
	var builder strings.Builder
	errlib.PanicError(tmp.Execute(&builder, data), fmt.Sprintf("Couldn't write %v template", location))
	return builder.String()
}
